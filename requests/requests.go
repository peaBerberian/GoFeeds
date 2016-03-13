package requests

import "net/http"
import "io/ioutil"
import "log"
import "github.com/peaberberian/OscarGoGo/config"
import "github.com/peaberberian/OscarGoGo/format"

// Launch requests on every given website "feedLink", parse them, and
// return them in the format.FeedFormat.
//
// This function will check if the cache can be used on each request
// intelligently (present for the website and does not exceed the
// given cacheTimeout).
//
// Because this function launch various error-prones routines, and does
// not exit on any of them, logs have been added in strategic points.
func GetFeeds(websites []config.Website, cache *feedCache) []format.FeedFormat {
	var res []format.FeedFormat
	var websitesToFetch []config.Website

	// Checks cache
	for _, web := range websites {
		webCache, errCache := cache.GetCacheForId(web.Id)
		if errCache == nil {
			res = append(res, webCache)
		} else {
			websitesToFetch = append(websitesToFetch, web)
		}
	}

	fetchedFeeds := fetchFeeds(websitesToFetch, cache)
	res = append(res, fetchedFeeds...)
	return res
}

func fetchFeeds(webs []config.Website, cache *feedCache) []format.FeedFormat {
	var chanRequest []chan httpResponse
	var res []format.FeedFormat

	// performs the requests for every given websites
	for i, web := range webs {
		chanRequest = append(chanRequest, make(chan httpResponse))

		log.Printf("launching for %s\n", web.FeedLink)
		go fetchUrl(web.FeedLink, chanRequest[i])
	}

	// handle each response
	for i, web := range webs {
		// blocking until response i arrive
		// TODO take it in first arrived order
		response := <-chanRequest[i]

		if response.err != nil {
			log.Printf("HTTP Error for %s: %s", web.SiteName, response.err)
		} else {
			log.Printf("Response received for %s\n", web.SiteName)
			feedRes, errParse := format.ParseFeed(response.body, web)
			if errParse != nil {
				log.Printf("XML Parsing error for %s: %s", web.SiteName, errParse)
			} else {
				cache.SetCacheForId(web.Id, feedRes)
				res = append(res, feedRes)
			}
		}
	}
	return res
}

// Struct of a http request response returned by fetchUrl
type httpResponse struct {
	err  error  // error: http error, response parsing error
	body []byte // complete response body
}

// Fetch url and return to the given channel the http response once the
// request is finished.
// Close the channel when everything is done
func fetchUrl(url string, c chan<- httpResponse) {
	defer close(c)
	resp, err := http.Get(url)
	if err != nil {
		c <- httpResponse{err: err}
		return
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- httpResponse{err: errRead, body: body}
		return
	}

	c <- httpResponse{err: nil, body: body}
}

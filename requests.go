package main

import "net/http"
import "io/ioutil"
import "time"
import "log"

// Launch requests on every given website "feedLink", parse them, and
// return them in the feedFormat.
//
// This function will check if the cache can be used on each request
// intelligently (present for the website and does not exceed the
// given cacheTimeout).
//
// Because this function launch various error-prones routines, and does
// not exit on any of them, logs have been added in strategic points.
func fetchWebsitesRss(websites []website, cache *feedCache, cacheTimeout int) []feedFormat {
	var res []feedFormat
	var _fetch = func(webs []website) {
		var c []chan httpResponse

		// performs the requests for every given websites
		for i, web := range webs {
			c = append(c, make(chan httpResponse))

			log.Printf("launching for %s\n", web.feedLink)
			go fetchUrl(web.feedLink, c[i])
		}

		// handle each response
		for i, web := range webs {
			// blocking
			response := <-c[i]

			if response.err != nil {
				log.Printf("HTTP Error for %s: %s", web.siteName, response.err)
			} else {
				log.Printf("Response received for %s\n", web.siteName)
				feedRes, errParse := parseFeed(response.body, web)
				if errParse != nil {
					log.Printf("XML Parsing error for %s: %s", web.siteName, errParse)
				} else {
					cache.set(web.id, feedRes)
					res = append(res, feedRes)
				}
			}
		}
	}

	// calculate here if we can use the cache for the wanted requests
	if cacheTimeout > 0 {
		var websitesToFetch []website

		// Checks cache
		for _, web := range websites {
			var shouldFetch = true
			webCache, errCache := cache.get(web.id)
			if errCache == nil {
				var cacheDate = webCache.date
				var deltaNano = time.Now().Nanosecond() - cacheDate.Nanosecond()
				if (deltaNano / 1000) < cacheTimeout {
					shouldFetch = false
					res = append(res, webCache.cache)
				}
			}
			if shouldFetch {
				websitesToFetch = append(websitesToFetch, web)
			}
		}
		_fetch(websitesToFetch)
	} else {
		// fetch all websites directly if no cache is set
		_fetch(websites)
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

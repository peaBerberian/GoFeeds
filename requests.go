package main

import "net/http"
import "io/ioutil"
import "time"
import "log"

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
	if cacheTimeout > 0 {
		var websitesToFetch []website

		// Checks cache
		for _, web := range websites {
			var shouldFetch = true
			if cache.has(web.id) {
				webCache, _ := cache.get(web.id)
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

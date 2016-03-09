package main

import "net/http"
import "io/ioutil"
import "encoding/xml"
import "time"

import "fmt"

func fetchWebsitesRss(websites []website, cache *feedCache, cacheTimeout int) []feedFormat {
	var res []feedFormat
	var _fetch = func(webs []website) {
		var c []chan httpResponse

		// performs the requests for every given websites
		for i, web := range webs {
			c = append(c, make(chan httpResponse))

			fmt.Printf("launching for %s\n", web.feedLink)
			go fetchUrl(web.feedLink, c[i])
		}

		// handle each response
		for i, web := range webs {
			response := <-c[i]
			if response.err == nil {
				fmt.Printf("received for %s\n", web.siteName)
				var feedRes feedFormat
				var parsed = false
				switch web.feedFormat {

				// parse RSS Feeds
				case "rss":
					var xmlBody rssFormat
					err := xml.Unmarshal(response.body, &xmlBody)
					if err == nil {
						feedRes = parseRss(xmlBody, web)
						parsed = true
					}

				// parse Atom feeds
				case "atom":
					var xmlBody atomFormat
					err := xml.Unmarshal(response.body, &xmlBody)
					if err == nil {
						feedRes = parseAtom(xmlBody, web)
						parsed = true
					}

				// Try to autodetect Feed type (duck-typing)
				default:
					var err error
					feedRes, err = parseFeed(response.body, web)
					if err == nil {
						parsed = true
					}
				}
				if parsed {
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
				var cacheDate = webCache.Date
				var deltaNano = time.Now().Nanosecond() - cacheDate.Nanosecond()
				if (deltaNano / 1000) < cacheTimeout {
					shouldFetch = false
					res = append(res, webCache.Cache)
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

package main

import "fmt"
import "log"
import "net/http"
import config "github.com/peaberberian/OscarGoGo/config"

func main() {
	log.Printf("starting application")
	start()
}

func start() {
	var cache feedCache
	var cacheTimeout int

	var conf, readErr = config.GetConfig()

	if readErr != nil {
		panic(readErr)
	}
	cacheTimeout = conf.CacheTime

	// Fill cache before starting
	_ = fetchWebsitesRss(conf.Websites, &cache, cacheTimeout)

	log.Printf("launching server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/feeds" && r.URL.Path != "/websites" {
			http.NotFound(w, r)
			return
		}

		if r.Method == "GET" {
			w.Header().Set("content-type", "application/json")
			if r.URL.Path == "/feeds" {
				log.Printf("GET /feeds")
				res := fetchWebsitesRss(conf.Websites, &cache, cacheTimeout)
				jsonRet, err := convertFeedsToJson(res)
				if err != nil {
					fmt.Fprintf(w, "[]")
				} else {
					fmt.Fprintf(w, "%s", jsonRet)
				}
				log.Printf("done")
			} else if r.URL.Path == "/websites" {
				log.Printf("GET /websites")
				ret, err := convertWebsitesToJson(conf.Websites)
				if err != nil {
					fmt.Fprintf(w, "[]")
				} else {
					fmt.Fprintf(w, "%s", ret)
				}
			}
		} else {
			http.Error(w, "Invalid request method.", 405)
		}
	})

	log.Printf("server launched")
	log.Fatal(http.ListenAndServe(":5013", nil))
}

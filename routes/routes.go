package routes

import "fmt"
import "log"
import "net/http"

import "github.com/peaberberian/OscarGoGo/config"
import reqs "github.com/peaberberian/OscarGoGo/requests"
import "github.com/peaberberian/OscarGoGo/format"

func StartServer(conf config.AppConfig) {
	var cache reqs.FeedCache
	var cacheTimeout int
	cacheTimeout = conf.CacheTime

	// Fill cache before starting
	_ = reqs.FetchWebsitesRss(conf.Websites, &cache, cacheTimeout)

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
				res := reqs.FetchWebsitesRss(conf.Websites, &cache, cacheTimeout)
				jsonRet, err := format.ConvertFeedsToJson(res)
				if err != nil {
					fmt.Fprintf(w, "[]")
				} else {
					fmt.Fprintf(w, "%s", jsonRet)
				}
				log.Printf("done")
			} else if r.URL.Path == "/websites" {
				log.Printf("GET /websites")
				ret, err := format.ConvertWebsitesToJson(conf.Websites)
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

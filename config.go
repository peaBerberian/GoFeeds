// This file provides structures definitions for the config file
// (./config.json) accessible through the getConfig function.

package main

import "encoding/json"
import "os"

const CONFIG_FILE_PATH = "./config.json"

type website struct {
	id          int
	siteName    string
	siteLink    string
	feedFormat  string
	feedName    string
	feedLink    string
	description string
}

type appConfig struct {
	websites  []website
	cacheTime int
}

type configFile struct {
	// Websites configuration
	Websites []struct {
		SiteName    string `json:"siteName"`
		SiteLink    string `json:"siteLink"`
		FeedFormat  string `json:"feedFormat"`
		FeedName    string `json:"feedName"`
		FeedLink    string `json:"feedLink"`
		Description string `json:"description"`
	} `json:"websites"`

	// Ammount of time in ns the cache is considered
	Cache int `json:"cache"`
}

// Returns map of current config file
// Can return an error if:
// 	 - the config file could not be read
// 	 - the config file could not be decoded
func getConfig() (appConfig, error) {
	var ret appConfig
	f, err := os.Open(CONFIG_FILE_PATH)
	if err != nil {
		return appConfig{}, err
	}
	var conf configFile
	dec := json.NewDecoder(f)
	if decodeErr := dec.Decode(&conf); decodeErr != nil {
		return appConfig{}, decodeErr
	}
	ret = appConfig{
		cacheTime: conf.Cache,
	}
	for i, item := range conf.Websites {
		ret.websites = append(ret.websites, website{
			id:          i,
			siteName:    item.SiteName,
			siteLink:    item.SiteLink,
			feedFormat:  item.FeedFormat,
			feedName:    item.FeedName,
			feedLink:    item.FeedLink,
			description: item.Description,
		})
	}
	return ret, nil
}

// This file provides structures definitions for the config file
// (./config.json) accessible through the getConfig function.

package config

import "encoding/json"
import "os"

const CONFIG_FILE_PATH = "./config/config.json"

type Website struct {
	Id          int
	SiteName    string
	SiteLink    string
	FeedFormat  string
	FeedName    string
	FeedLink    string
	Description string
}

type appConfig struct {
	Websites  []Website
	CacheTime int
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
func GetConfig() (appConfig, error) {
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
		CacheTime: conf.Cache,
	}
	for i, item := range conf.Websites {
		ret.Websites = append(ret.Websites, Website{
			Id:          i,
			SiteName:    item.SiteName,
			SiteLink:    item.SiteLink,
			FeedFormat:  item.FeedFormat,
			FeedName:    item.FeedName,
			FeedLink:    item.FeedLink,
			Description: item.Description,
		})
	}
	return ret, nil
}

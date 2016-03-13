// This file provides structures definitions for the config file
// (./config.json) accessible through the getConfig function.

package config

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

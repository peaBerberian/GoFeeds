package config

import "encoding/json"
import "os"

const CONFIG_FILE_PATH = "./config/config.json"

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

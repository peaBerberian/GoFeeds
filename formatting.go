package main

import "encoding/json"
import "encoding/xml"
import "time"
import "errors"
import config "github.com/peaberberian/OscarGoGo/config"

type rssFormat struct {
	Channels struct {
		Title         string `xml:"title"`
		Description   string `xml:"description"`
		Link          string `xml:"link"`
		LastBuildDate string `xml:"lastBuildDate"`
		PubDate       string `xml:"pubDate"`
		Ttl           int    `xml:"ttl"`
		Items         []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type atomFormat struct {
	Title    string `xml:"title"`
	Subtitle string `xml:"subtitle"`
	Links    []struct {
		Key string `xml:"href,attr"`
	} `xml:"link"`
	Id      string `xml:"id"`
	Updated string `xml:"updated"`
	Entries []struct {
		Title string `xml:"title"`
		Links []struct {
			Key string `xml:"href,attr"`
		} `xml:"link"`
		Id      string `xml:"id"`
		Updated string `xml:"updated"`
		Summary string `xml:"summary"`
		Content string `xml:"content"`
		Author  struct {
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"author"`
	} `xml:"entry"`
}

// standard feed format used, most notably in the feedCache
type feedFormat struct {
	Id          int
	Title       string
	Description string
	UpdateDate  time.Time
	Link        string
	Entries     []feedEntry
}

type feedEntry struct {
	Title        string
	Link         string
	Description  string
	CreationDate time.Time
	UpdateDate   time.Time
}

// ---- JSON API ----

type jsonItem struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type jsonFormat struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	Link  string     `json:"link"`
	Items []jsonItem `json:"items"`
}

type jsonResponse []jsonFormat

type websiteJSON struct {
	Id          int    `json:"id"`
	SiteName    string `json:"siteName"`
	SiteLink    string `json:"siteLink"`
	FeedFormat  string `json:"feedFormat"`
	FeedName    string `json:"feedName"`
	FeedLink    string `json:"feedLink"`
	Description string `json:"description"`
}

// Convert an rssFormat to a feedFormat
func parseRss(rssMap rssFormat, web config.Website) feedFormat {
	var feedTime = parseRssTime(rssMap.Channels.PubDate)
	var feed = feedFormat{
		Id:          web.Id,
		Title:       rssMap.Channels.Title,
		Link:        web.FeedLink,
		Description: rssMap.Channels.Description,
		UpdateDate:  feedTime,
	}
	for _, item := range rssMap.Channels.Items {
		var date = parseRssTime(item.PubDate)

		feed.Entries = append(feed.Entries, feedEntry{
			Title:        item.Title,
			Link:         item.Link,
			Description:  item.Description,
			CreationDate: date,
			UpdateDate:   date,
		})
	}
	return feed
}

// Convert an atomFormat to a feedFormat
func parseAtom(atomMap atomFormat, web config.Website) feedFormat {
	var feedTime = parseAtomTime(atomMap.Updated)
	var feed = feedFormat{
		Id:          web.Id,
		Title:       atomMap.Title,
		Link:        web.FeedLink,
		Description: atomMap.Subtitle,
		UpdateDate:  feedTime,
	}
	for _, item := range atomMap.Entries {
		var date = parseAtomTime(item.Updated)

		var description string
		if item.Summary != "" {
			description = item.Summary
		} else {
			description = item.Content
		}
		feed.Entries = append(feed.Entries, feedEntry{
			Title:        item.Title,
			Link:         item.Links[0].Key,
			Description:  description,
			CreationDate: date,
			UpdateDate:   date,
		})
	}
	return feed
}

func convertFeedsToJson(feeds []feedFormat) ([]byte, error) {
	var jsobjs = []jsonFormat{}
	for _, feed := range feeds {
		var jsobj = jsonFormat{
			Id:   feed.Id,
			Name: feed.Title,
			Link: feed.Link,
		}

		for _, entry := range feed.Entries {
			jsobj.Items = append(jsobj.Items, jsonItem{
				Title:        entry.Title,
				Link:         entry.Link,
				Description:  entry.Description,
				CreationDate: timeToString(entry.CreationDate),
			})
		}
		jsobjs = append(jsobjs, jsobj)
	}
	res, err := json.Marshal(jsobjs)
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}

// needed or not?
func autoDetectFeedFormat(raw []byte) (string, error) {
	var rssRaw rssFormat
	var atomRaw atomFormat
	errRss := xml.Unmarshal(raw, &rssRaw)
	errAtom := xml.Unmarshal(raw, &atomRaw)
	if errRss == nil && (len(rssRaw.Channels.Items) > 0 ||
		rssRaw.Channels.Title != "") {
		return "rss", nil
	}

	if errAtom == nil && (len(atomRaw.Entries) > 0 ||
		atomRaw.Title != "") {
		return "atom", nil
	}

	if errRss != nil {
		return "", errRss
	}
	if errAtom != nil {
		return "", errAtom
	}

	return "", errors.New("Could not detect your feed format")
}

func parseFeed(raw []byte, web config.Website) (feedFormat, error) {
	var feedRes feedFormat

	switch web.FeedFormat {

	// parse RSS Feeds
	case "rss":
		var xmlBody rssFormat
		err := xml.Unmarshal(raw, &xmlBody)
		if err != nil {
			return feedFormat{}, err
		} else {
			feedRes = parseRss(xmlBody, web)
		}

	// parse Atom feeds
	case "atom":
		var xmlBody atomFormat
		err := xml.Unmarshal(raw, &xmlBody)
		if err != nil {
			return feedFormat{}, err
		} else {
			feedRes = parseAtom(xmlBody, web)
		}

	// Try to autodetect Feed type (duck-typing)
	default:
		var rssRaw rssFormat
		var atomRaw atomFormat
		errRss := xml.Unmarshal(raw, &rssRaw)
		errAtom := xml.Unmarshal(raw, &atomRaw)

		if errRss == nil && (len(rssRaw.Channels.Items) > 0 ||
			rssRaw.Channels.Title != "") {
			ret := parseRss(rssRaw, web)
			return ret, nil
		}

		if errAtom == nil && (len(atomRaw.Entries) > 0 ||
			atomRaw.Title != "") {
			ret := parseAtom(atomRaw, web)
			return ret, nil
		}

		if errRss != nil {
			return feedFormat{}, errRss
		}
		if errAtom != nil {
			return feedFormat{}, errAtom
		}
		return feedFormat{}, errors.New("Could not detect your feed format")
	}
	return feedRes, nil
}

func convertWebsitesToJson(webs []config.Website) ([]byte, error) {
	var websJson []websiteJSON
	for _, web := range webs {
		websJson = append(websJson, websiteJSON{
			Id:          web.Id,
			Description: web.Description,
			FeedLink:    web.FeedLink,
			FeedName:    web.FeedName,
			FeedFormat:  web.FeedFormat,
			SiteLink:    web.SiteLink,
			SiteName:    web.SiteName,
		})
	}
	ret, err := json.Marshal(websJson)
	if err != nil {
		return []byte{}, err
	}
	return ret, nil
}

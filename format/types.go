package format

import "time"

// standard feed format used in the application
type FeedFormat struct {
	Id          int
	Title       string
	Description string
	UpdateDate  time.Time
	Link        string
	Entries     []feedEntry
}

// single feed entry for the FeedFormat
type feedEntry struct {
	Title        string
	Link         string
	Description  string
	CreationDate time.Time
	UpdateDate   time.Time
}

// RSS xml format used on unMarshall
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

// RSS atom format used on unMarshall
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

// ---- JSON API ----

// used on json.marshall for constructing the API response
type jsonFormat struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	Link  string     `json:"link"`
	Items []jsonItem `json:"items"`
}

// the API response for feeds will be an array of jsonFormat
type jsonResponse []jsonFormat

// single entry for the jsonFormat type
type jsonItem struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

// used on json.marshall to describe the websites used in the 'websites'
// api
type websiteJSON struct {
	Id          int    `json:"id"`
	SiteName    string `json:"siteName"`
	SiteLink    string `json:"siteLink"`
	FeedFormat  string `json:"feedFormat"`
	FeedName    string `json:"feedName"`
	FeedLink    string `json:"feedLink"`
	Description string `json:"description"`
}

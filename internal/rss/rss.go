package rss

// RSSFeed represents the top-level feed structure
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Link        string    `xml:"link"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents a single post/item in the feed
type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"` // can parse to time later
}

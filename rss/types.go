package rss

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

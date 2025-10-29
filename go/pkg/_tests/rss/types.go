package rss

import "encoding/xml"

// RSS represents the structure of an RSS feed
type RSS struct {
	XMLName       xml.Name `xml:"rss"`
	Version       string   `xml:"version,attr"`
	XMLNSContent  string   `xml:"xmlns:content,attr"`
	XMLNSDC       string   `xml:"xmlns:dc,attr"`
	XMLNSAtom     string   `xml:"xmlns:atom,attr"`
	Channel       Channel  `xml:"channel"`
}

// Channel represents the channel element in an RSS feed
type Channel struct {
	XMLName       xml.Name  `xml:"channel"`
	Title         string    `xml:"title"`
	Link          string    `xml:"link,omitempty"`
	LinkCData     string    `xml:"link,omitempty,cdata"`
	Description   string    `xml:"description"`
	AtomLink      AtomLink  `xml:"http://www.w3.org/2005/Atom link"`
	LastBuildDate string    `xml:"lastBuildDate"`
	PubDate       string    `xml:"pubDate"`
	Published     string    `xml:"http://www.w3.org/2005/Atom published"`
	Updated       string    `xml:"http://www.w3.org/2005/Atom updated"`
	Categories    []string  `xml:"category"`
	Copyright     string    `xml:"copyright"`
	Image         Image     `xml:"image"`
	Docs          string    `xml:"docs"`
	Generator     string    `xml:"generator"`
	Language      string    `xml:"language"`
	WebMaster     string    `xml:"webMaster"`
	Items         []Item    `xml:"item"`
}

// GetLink returns the channel link, checking both regular and CDATA links
func (c *Channel) GetLink() string {
	if c.Link != "" {
		return c.Link
	}
	return c.LinkCData
}

// AtomLink represents an atom:link element
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// Image represents an image in the RSS feed
type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

// Item represents an item in an RSS feed
type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
	GUID        string   `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
	Published   string   `xml:"atom:published"`
	Creator     string   `xml:"dc:creator"`
	Categories  []string `xml:"category"`
	Content     string   `xml:"content:encoded"`
}

// Enclosure represents an enclosure in an RSS item
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

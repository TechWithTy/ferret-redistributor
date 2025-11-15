package rss

import (
    "encoding/xml"
    "net/url"
    "strings"
    "time"
)

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
    Description   string    `xml:"description"`
    AtomLink      AtomLink  `xml:"http://www.w3.org/2005/Atom link"`
    LastBuildDate string    `xml:"lastBuildDate"`
    PubDate       string    `xml:"pubDate"`
    Published     string    `xml:"http://www.w3.org/2005/Atom published"`
    Updated       string    `xml:"http://www.w3.org/2005/Atom updated"`
    Categories    []Category `xml:"category"`
    Copyright     string    `xml:"copyright"`
    Image         Image     `xml:"image"`
    Docs          string    `xml:"docs"`
    Generator     string    `xml:"generator"`
    Language      string    `xml:"language"`
    WebMaster     string    `xml:"webMaster"`
    Items         []Item    `xml:"item"`

    // Parsed times (not from XML)
    PubDateTime    time.Time `xml:"-"`
    PublishedTime  time.Time `xml:"-"`
    UpdatedTime    time.Time `xml:"-"`
}

// GetLink returns the channel link, checking both regular and CDATA links
func (c Channel) GetLink() string {
    if c.Link != "" {
        return c.Link
    }
    if c.AtomLink.Href != "" {
        return c.AtomLink.Href
    }
    if len(c.Items) > 0 {
        raw := strings.TrimSpace(c.Items[0].Link)
        if raw != "" {
            if u, err := url.Parse(raw); err == nil && u.Scheme != "" && u.Host != "" {
                // Return site root as a sensible fallback
                base := u.Scheme + "://" + u.Host
                if !strings.HasSuffix(base, "/") {
                    base += "/"
                }
                return base
            }
        }
    }
    return ""
}

// AtomLink represents an atom:link element
type AtomLink struct {
    Href string `xml:"href,attr"`
    Rel  string `xml:"rel,attr"`
    Type string `xml:"type,attr"`
    Title string `xml:"title,attr"`
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
    GUID        GUID     `xml:"guid"`
    PubDate     string   `xml:"pubDate"`
    Published   string   `xml:"http://www.w3.org/2005/Atom published"`
    Creator     string   `xml:"http://purl.org/dc/elements/1.1/ creator"`
    Categories  []Category `xml:"category"`
    Content     string   `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`

    // Parsed times (not from XML)
    PubDateTime   time.Time `xml:"-"`
    PublishedTime time.Time `xml:"-"`
}

// GUID captures guid text and attributes like isPermaLink
type GUID struct {
    IsPermaLink bool   `xml:"isPermaLink,attr"`
    Value       string `xml:",chardata"`
}

// Enclosure represents an enclosure in an RSS item
type Enclosure struct {
    URL    string `xml:"url,attr"`
    Length int64  `xml:"length,attr"`
    Type   string `xml:"type,attr"`
}

// Category represents an RSS category with optional domain attribute
type Category struct {
    Domain string `xml:"domain,attr"`
    Value  string `xml:",chardata"`
}

// String returns the category value when printed with %v
func (c Category) String() string { return strings.TrimSpace(c.Value) }

// CanonicalURL trims and lowercases host, preserving path/query/fragment
func CanonicalURL(s string) string {
    s = strings.TrimSpace(s)
    if s == "" {
        return s
    }
    u, err := url.Parse(s)
    if err != nil || u.Scheme == "" || u.Host == "" {
        return s
    }
    u.Host = strings.ToLower(u.Host)
    // leave scheme, path, query as-is; ensure no whitespace leakage
    return u.String()
}

// Normalize trims strings, canonicalizes URLs, and parses dates for channel and items
func (c *Channel) Normalize() {
    c.Title = strings.TrimSpace(c.Title)
    c.Description = strings.TrimSpace(c.Description)
    c.Link = CanonicalURL(c.Link)
    if c.Link == "" {
        // Populate with fallback to ensure a usable link
        c.Link = CanonicalURL(c.GetLink())
    }
    // Parse channel dates
    c.PubDateTime = parseAnyDate(c.PubDate)
    c.PublishedTime = parseAnyDate(c.Published)
    c.UpdatedTime = parseAnyDate(c.Updated)
    // Normalize items
    for i := range c.Items {
        c.Items[i].Normalize()
    }
}

// Normalize trims and parses common fields on an item
func (it *Item) Normalize() {
    it.Title = strings.TrimSpace(it.Title)
    it.Description = strings.TrimSpace(it.Description)
    it.Link = CanonicalURL(it.Link)
    it.PubDateTime = parseAnyDate(it.PubDate)
    it.PublishedTime = parseAnyDate(it.Published)
}

// parseAnyDate tries several common RSS/Atom date formats
func parseAnyDate(s string) time.Time {
    s = strings.TrimSpace(s)
    if s == "" {
        return time.Time{}
    }
    // Try common formats
    formats := []string{
        time.RFC3339Nano,
        time.RFC3339,
        time.RFC1123Z,
        time.RFC1123,
        time.RFC822Z,
        time.RFC822,
    }
    for _, f := range formats {
        if t, err := time.Parse(f, s); err == nil {
            return t
        }
    }
    return time.Time{}
}

package landing

import (
    "bytes"
    _ "embed"
    "encoding/xml"
    "io"
    "net/http"
    "os"
    "strings"
)

//go:embed static-sitemap.xml
var staticSitemap []byte

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
    XMLName xml.Name `xml:"url"`
    Loc     string   `xml:"loc"`
}

func FetchSiteMap() ([]URL, error) {
    sitemapURL := "https://www.dealscale.io/sitemap.xml"
    if envSiteMapURL := os.Getenv("SITEMAP_URL"); envSiteMapURL != "" {
        sitemapURL = envSiteMapURL
    }
    var data []byte

    if strings.EqualFold(sitemapURL, "static") {
        data = staticSitemap
    } else {
        resp, err := http.Get(sitemapURL)
        if err == nil && resp != nil {
            defer resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                b, rerr := io.ReadAll(resp.Body)
                if rerr == nil {
                    ct := resp.Header.Get("Content-Type")
                    lb := bytes.ToLower(b)
                    // Fallback if it looks like HTML instead of XML
                    if !strings.Contains(strings.ToLower(ct), "text/html") && !bytes.Contains(lb, []byte("<!doctype html")) {
                        data = b
                    }
                }
            }
        }

        if len(data) == 0 {
            // Fallback to embedded static sitemap when fetch fails or returns HTML
            data = staticSitemap
        }
    }

    var sitemap URLSet
    if err := xml.Unmarshal(data, &sitemap); err != nil {
        return nil, err
    }

    return sitemap.URLs, nil
}

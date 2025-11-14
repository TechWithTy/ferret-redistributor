package rsshub

// FeedRequest controls how an RSSHub feed is fetched.
type FeedRequest struct {
	// Path is the RSSHub route path, e.g. "/bilibili/fav/2262573".
	Path string
	// Query allows the caller to append query parameters such as ?lang=en or ?limit=10.
	Query map[string]string
	// Format sets RSSHub's format parameter (e.g. "debug.json" or "0.debug.html").
	Format string
}

// ForceRefreshRequest encodes the URL to refresh.
type ForceRefreshRequest struct {
	TargetURL string
}

// RadarSearchRequest describes a lookup against /api/radar/search.
type RadarSearchRequest struct {
	// URL is the page we want RSSHub Radar to inspect.
	URL string
}

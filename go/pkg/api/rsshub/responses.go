package rsshub

// RoutesResponse mirrors /api/routes output.
type RoutesResponse struct {
	Status string                     `json:"status"`
	Data   map[string]RoutesGroupData `json:"data"`
}

// RoutesGroupData groups routes per category.
type RoutesGroupData struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Routes      []RouteDetail `json:"routes"`
}

// RouteDetail contains metadata about a single RSSHub route.
type RouteDetail struct {
	Path        string          `json:"path"`
	Description string          `json:"description"`
	Parameters  []RouteParam    `json:"parameters"`
	Examples    []RouteExample  `json:"examples"`
	Author      string          `json:"author"`
	Source      string          `json:"source"`
	Items       []RouteItemDesc `json:"items"`
}

// RouteParam captures a required/optional parameter.
type RouteParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// RouteExample shows how to use the route.
type RouteExample struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// RouteItemDesc documents feed fields.
type RouteItemDesc struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// VersionResponse returns the running RSSHub version string.
type VersionResponse struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// ForceRefreshResponse returns metadata for a refresh request.
type ForceRefreshResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RadarSearchResponse contains RSSHub Radar matches.
type RadarSearchResponse struct {
	Results []RadarSearchItem `json:"results"`
}

// RadarSearchItem enumerates a route candidate.
type RadarSearchItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

// FeedResult represents a fetched RSS/Atom feed.
type FeedResult struct {
	ContentType string
	Body        []byte
}

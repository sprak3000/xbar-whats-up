package service

// Site holds the data for service status pages
type Site struct {
	URL  string `json:"url"`
	Slug string `json:"slug"`
}

// Sites is a mapping of services to their status page data
type Sites map[string]Site

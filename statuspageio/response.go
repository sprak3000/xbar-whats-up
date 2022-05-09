package statuspageio

import "time"

// Page represents the page information
type Page struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	TimeZone  string    `json:"string"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status represents the status information
type Status struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
}

// Response is the structure returned by statuspage.io powered service status pages
type Response struct {
	Page   Page   `json:"page"`
	Status Status `json:"status"`
}

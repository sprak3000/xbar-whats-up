// Package statuspageio handles communicating with status pages following the statuspage.io format
package statuspageio

import (
	"time"
)

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

// Indicator returns the current status for the service
func (r Response) Indicator() string {
	return r.Status.Indicator
}

// Name returns the name of the service
func (r Response) Name() string {
	return r.Page.Name
}

// UpdatedAt returns the date / time of the most recent status update for the service
func (r Response) UpdatedAt() time.Time {
	return r.Page.UpdatedAt
}

// URL returns the page the service hosts to display more details about current and past status updates
func (r Response) URL() string {
	return r.Page.URL
}

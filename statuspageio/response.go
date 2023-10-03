// Package statuspageio handles communicating with status pages following the statuspage.io format
package statuspageio

import (
	"time"

	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/go-whatsup-client/whatsup"

	"github.com/sprak3000/xbar-whats-up/status"
)

// ServiceType is the name we use for various checks
const ServiceType = "statuspage.io"

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

// ClientReader implements the Reader interface for go-client based reading of a service's status
type ClientReader struct {
	ServiceName string
	PageURL     string
}

// ReadStatus handles communicating with the service to get its status details
func (cr ClientReader) ReadStatus(client whatsup.StatusPageClient) (status.Details, glitch.DataError) {
	return client.StatuspageIoService(cr.ServiceName, cr.PageURL)
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

// Package statuspageio handles communicating with status pages following the statuspage.io format
package statuspageio

import (
	"github.com/sprak3000/go-glitch/glitch"
	whatsupstatus "github.com/sprak3000/go-whatsup-client/status"
	"github.com/sprak3000/go-whatsup-client/whatsup"
)

// ServiceType is the name we use for various checks
const ServiceType = "statuspage.io"

// ClientReader implements the Reader interface for go-client based reading of a service's status
type ClientReader struct {
	ServiceName string
	PageURL     string
}

// ReadStatus handles communicating with the service to get its status details
func (cr ClientReader) ReadStatus(client whatsup.StatusPageClient) (whatsupstatus.Details, glitch.DataError) {
	return client.StatuspageIoService(cr.ServiceName, cr.PageURL)
}

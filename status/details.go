// Package status is an abstraction for handling and displaying status details from various services
package status

import (
	"context"
	"net/http"
	"time"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"
)

// Error codes
const (
	ErrorUnableToMakeClientRequest   = "UNABLE_TO_MAKE_CLIENT_REQUEST"
	ErrorUnableToParseClientResponse = "UNABLE_TO_PARSE_CLIENT_RESPONSE"
)

// Details provides an interface for extracting information from a service's status response
type Details interface {
	Indicator() string
	Name() string
	UpdatedAt() time.Time
	URL() string
}

// Get makes the network request to get a status page
func Get(serviceFinder client.ServiceFinder, serviceName, slug string) ([]byte, glitch.DataError) {
	c := client.NewBaseClient(serviceFinder, serviceName, true, 10*time.Second, nil)
	_, respBytes, err := c.MakeRequest(context.Background(), http.MethodGet, slug, nil, nil, nil)
	return respBytes, err
}

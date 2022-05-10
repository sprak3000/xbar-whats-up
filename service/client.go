package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/xbar-whats-up/statuspageio"
)

// Error codes
const (
	ErrorUnableToMakeClientRequest   = "UNABLE_TO_MAKE_CLIENT_REQUEST"
	ErrorUnableToParseClientResponse = "UNABLE_TO_PARSE_CLIENT_RESPONSE"
)

// Reader provides the requirements for anyone implementing reading a service's status
type Reader interface {
	ReadStatus(serviceFinder client.ServiceFinder, serviceName, slug string) (statuspageio.Response, glitch.DataError)
}

// ClientReader implements the Reader interface for go-client based reading of a service's status
type ClientReader struct {
}

// ReadStatus allows us to use go-client to read a service's status
func (cr ClientReader) ReadStatus(serviceFinder client.ServiceFinder, serviceName, slug string) (statuspageio.Response, glitch.DataError) {
	resp := statuspageio.Response{}

	c := client.NewBaseClient(serviceFinder, serviceName, true, 10*time.Second, nil)
	_, respBytes, err := c.MakeRequest(context.Background(), "GET", slug, nil, nil, nil)
	if err != nil {
		return resp, glitch.NewDataError(err, ErrorUnableToMakeClientRequest, fmt.Sprintf("unable to make client request for %s: %v", serviceName, err))
	}

	uErr := json.Unmarshal(respBytes, &resp)
	if uErr != nil {
		return resp, glitch.NewDataError(err, ErrorUnableToParseClientResponse, fmt.Sprintf("unable to parse client response for %s: %v", serviceName, err))
	}

	return resp, nil
}

// NewClientServiceFinder returns a client.ServiceFinder suitable for use with go-client
func NewClientServiceFinder(sites Sites) client.ServiceFinder {
	return func(serviceName string, useTLS bool) (url.URL, error) {
		u, ok := sites[serviceName]
		if !ok {
			return url.URL{}, errors.New("unable to find " + serviceName + " in the service list")
		}

		return u.URL, nil
	}
}

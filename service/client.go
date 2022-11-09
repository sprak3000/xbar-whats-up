// Package service handles communicating with sites to obtain their current status details
package service

import (
	"errors"
	"net/url"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/xbar-whats-up/slack"
	"github.com/sprak3000/xbar-whats-up/status"
	"github.com/sprak3000/xbar-whats-up/statuspageio"
)

// Reader provides the requirements for anyone implementing reading a service's status
type Reader interface {
	ReadStatus(serviceFinder client.ServiceFinder, serviceName, slug string) (status.Details, glitch.DataError)
}

// ReaderServiceFinder is an alias for the closure that finds a Reader based on the given service type
type ReaderServiceFinder func(serviceType string) (Reader, error)

// NewReaderServiceFinder returns a service.ReaderServiceFinder to provider a Reader based on the given service type
func NewReaderServiceFinder() ReaderServiceFinder {
	return func(serviceType string) (Reader, error) {
		switch serviceType {
		case statuspageio.ServiceType:
			return statuspageio.ClientReader{}, nil
		case slack.ServiceType:
			return slack.ClientReader{}, nil
		default:
			return nil, errors.New("reader not implemented for type " + serviceType)
		}
	}
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

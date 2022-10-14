// Package status is an abstraction for handling and displaying status details from various services
package status

import (
	"time"
)

// Details provides an interface for extracting information from a service's status response
type Details interface {
	Indicator() string
	Name() string
	UpdatedAt() time.Time
	URL() string
}

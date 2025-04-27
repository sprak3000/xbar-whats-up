// Package service handles communicating with sites to obtain their current status details
package service

import (
	"encoding/json"
	"net/url"

	"github.com/sprak3000/go-glitch/glitch"
	whatsupstatus "github.com/sprak3000/go-whatsup-client/status"
	"github.com/sprak3000/go-whatsup-client/whatsup"

	"github.com/sprak3000/xbar-whats-up/configuration"
	"github.com/sprak3000/xbar-whats-up/slack"
	"github.com/sprak3000/xbar-whats-up/status"
	"github.com/sprak3000/xbar-whats-up/statuspageio"
)

// Error codes
const (
	ErrorUnableToWriteDefaultConfiguration = "UNABLE_TO_WRITE_DEFAULT_CONFIGURATION"
	ErrorUnableToParseConfiguration        = "UNABLE_TO_PARSE_CONFIGURATION"
	ErrorUnsupportedServiceType            = "UNSUPPORTED_SERVICE_TYPE"
)

// Reader provides the requirements for anyone implementing reading a service's status
type Reader interface {
	ReadStatus(client whatsup.StatusPageClient) (whatsupstatus.Details, glitch.DataError)
}

// Site holds the data for service status pages
type Site struct {
	URL  url.URL `json:"url,string"`
	Type string  `json:"type"`
}

// UnmarshalJSON handles converting data into the Site type
func (s *Site) UnmarshalJSON(data []byte) error {
	type tmpSite Site

	tmp := struct {
		URL string `json:"url"`
		*tmpSite
	}{
		tmpSite: (*tmpSite)(s),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	u, pErr := url.Parse(tmp.URL)
	if pErr != nil {
		return pErr
	}

	s.URL = *u

	return nil
}

// Sites is a mapping of services to their status page data
type Sites map[string]Site

type readerResult struct {
	details whatsupstatus.Details
	err     glitch.DataError
}

func readStatusPage(c whatsup.StatusPageClient, serviceName string, s Site) readerResult {
	var reader Reader

	switch s.Type {
	case statuspageio.ServiceType:
		reader = statuspageio.ClientReader{
			ServiceName: serviceName,
			PageURL:     s.URL.String(),
		}
	case slack.ServiceType:
		reader = slack.ClientReader{}
	default:
		// Unsupported at this time
		return readerResult{
			details: nil,
			err:     glitch.NewDataError(nil, ErrorUnsupportedServiceType, serviceName+" uses an unsupported service type "+s.Type),
		}
	}

	resp, err := reader.ReadStatus(c)
	return readerResult{
		details: resp,
		err:     err,
	}
}

// GetOverview returns the details about the services monitored
func (sites Sites) GetOverview(client whatsup.StatusPageClient) status.Overview {
	overview := status.Overview{
		OverallStatus: "none",
		List:          map[string][]whatsupstatus.Details{},
		Errors:        []status.OverviewError{},
	}

	c := make(chan readerResult)

	for k, v := range sites {
		serviceName := k
		site := v
		go func() { c <- readStatusPage(client, serviceName, site) }()
	}

	for i := 0; i < len(sites); i++ {
		resp := <-c

		if resp.err != nil {
			overview.Errors = append(overview.Errors, status.OverviewError{
				Details: resp.details,
				Error:   resp.err,
			})
			continue
		}

		nameSize := len(resp.details.Name())
		if nameSize > overview.LargestStringSize {
			overview.LargestStringSize = nameSize
		}

		switch resp.details.Indicator() {
		case "major":
			overview.OverallStatus = "major"
			overview.List["major"] = append(overview.List["major"], resp.details)
		case "minor":
			if overview.OverallStatus != "major" {
				overview.OverallStatus = "minor"
			}
			overview.List["minor"] = append(overview.List["minor"], resp.details)
		default:
			overview.List["none"] = append(overview.List["none"], resp.details)
		}
	}

	return overview
}

// LoadSites reads a JSON file containing a list of sites to monitor
func LoadSites(r configuration.Reader, w configuration.Writer, filename string) (Sites, glitch.DataError) {
	var sites Sites

	if filename == "" {
		filename = "./.whats-up.json"
	}

	config, rErr := r.ReadFile(filename)
	if rErr != nil {
		// Create an empty configuration file
		config = []byte("{\n}")
		wErr := w.WriteFile(filename, config, 0644)
		if wErr != nil {
			return sites, glitch.NewDataError(wErr, ErrorUnableToWriteDefaultConfiguration, "unable to create default What's Up configuration")
		}
	}

	uErr := json.Unmarshal(config, &sites)
	if uErr != nil {
		return sites, glitch.NewDataError(uErr, ErrorUnableToParseConfiguration, "error parsing What's Up configuration")
	}

	return sites, nil
}

// Package service handles communicating with sites to obtain their current status details
package service

import (
	"encoding/json"
	"net/url"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/xbar-whats-up/configuration"
	"github.com/sprak3000/xbar-whats-up/status"
)

// Error codes
const (
	ErrorUnableToWriteDefaultConfiguration = "UNABLE_TO_WRITE_DEFAULT_CONFIGURATION"
	ErrorUnableToParseConfiguration        = "UNABLE_TO_PARSE_CONFIGURATION"
)

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

// GetOverview returns the details about the services monitored
func (sites Sites) GetOverview(serviceFinder client.ServiceFinder, readerFinder ReaderServiceFinder) status.Overview {
	overview := status.Overview{
		OverallStatus: "none",
		List:          map[string][]status.Details{},
		Errors:        []string{},
	}

	for k, v := range sites {
		reader, rfErr := readerFinder(v.Type)
		if rfErr != nil {
			overview.Errors = append(overview.Errors, k+" uses an unsupported service type "+v.Type)
			// Unsupported at this time
			continue
		}

		resp, rErr := reader.ReadStatus(serviceFinder, k, v.URL.Path)
		if rErr != nil {
			overview.Errors = append(overview.Errors, rErr.Error())
			continue
		}

		switch resp.Indicator() {
		case "major":
			overview.OverallStatus = "major"
			overview.List["major"] = append(overview.List["major"], resp)
		case "minor":
			if overview.OverallStatus != "major" {
				overview.OverallStatus = "minor"
			}
			overview.List["minor"] = append(overview.List["minor"], resp)
		default:
			overview.List["none"] = append(overview.List["none"], resp)
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

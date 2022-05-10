package service

import (
	"encoding/json"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/xbar-whats-up/configuration"
	"github.com/sprak3000/xbar-whats-up/status"
	"github.com/sprak3000/xbar-whats-up/statuspageio"
)

// Error codes
const (
	ErrorUnableToWriteDefaultConfiguration = "UNABLE_TO_WRITE_DEFAULT_CONFIGURATION"
	ErrorUnableToParseConfiguration        = "UNABLE_TO_PARSE_CONFIGURATION"
)

// Site holds the data for service status pages
type Site struct {
	URL  string `json:"url"`
	Slug string `json:"slug"`
}

// Sites is a mapping of services to their status page data
type Sites map[string]Site

// GetOverview returns the details about the services monitored
func (s Sites) GetOverview(serviceFinder client.ServiceFinder, reader Reader) status.Overview {
	overview := status.Overview{
		OverallStatus: "🟢",
		List:          map[string][]statuspageio.Response{},
		Errors:        []string{},
	}

	for k, v := range s {
		resp, rErr := reader.ReadStatus(serviceFinder, k, v.Slug)
		if rErr != nil {
			overview.Errors = append(overview.Errors, rErr.Error())
			continue
		}

		switch resp.Status.Indicator {
		case "major":
			overview.OverallStatus = "🔴"
			overview.List["major"] = append(overview.List["major"], resp)
		case "minor":
			overview.OverallStatus = "🟠"
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

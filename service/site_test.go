package service

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/xbar-whats-up/configuration"
	"github.com/sprak3000/xbar-whats-up/status"
	"github.com/sprak3000/xbar-whats-up/statuspageio"
)

func TestUnit_Site_UnmarshalJSON(t *testing.T) {
	codeClimateURL, err := url.Parse("https://status.codeclimate.com/api/v2/status.json")
	require.NoError(t, err)

	tests := map[string]struct {
		siteJSON     []byte
		expectedSite Site
		expectedErr  error
		validate     func(t *testing.T, expectedSite, actualSite Site, expectedErr, actualErr error)
	}{
		"base path": {
			siteJSON: []byte(`{"url":"https://status.codeclimate.com/api/v2/status.json","type":"statuspage.io"}`),
			expectedSite: Site{
				URL:  *codeClimateURL,
				Type: "statuspage.io",
			},
			validate: func(t *testing.T, expectedSite, actualSite Site, expectedErr, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedSite, actualSite)
			},
		},
		"exceptional path- parse URL error": {
			siteJSON:    []byte(`{"url":":","type":"statuspage.io"}`),
			expectedErr: errors.New(`parse ":": missing protocol scheme`),
			validate: func(t *testing.T, expectedSite, actualSite Site, expectedErr, actualErr error) {
				require.Error(t, actualErr)
				require.Equal(t, expectedErr.Error(), actualErr.Error())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var s Site
			err := json.Unmarshal(tc.siteJSON, &s)
			tc.validate(t, tc.expectedSite, s, tc.expectedErr, err)
		})
	}
}

func TestUnit_GetOverview(t *testing.T) {
	codeClimateURL, err := url.Parse("https://status.codeclimate.com/api/v2/status.json")
	require.NoError(t, err)

	circleciURL, err := url.Parse("https://status.circleci.com/api/v2/status.json")
	require.NoError(t, err)

	tests := map[string]struct {
		sites            Sites
		finder           client.ServiceFinder
		reader           Reader
		expectedOverview status.Overview
		validate         func(t *testing.T, expectedOverview, actualOverview status.Overview)
	}{
		"base path- highest severity is major": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: "statuspage.io",
				},
				"CircleCI": {
					URL:  *circleciURL,
					Type: "statuspage.io",
				},
			},
			reader: clientReaderSeverityMajor{},
			expectedOverview: status.Overview{
				OverallStatus: "major",
				List: map[string][]status.Details{
					"major": {
						statuspageio.Response{
							Page: statuspageio.Page{
								ID:   "circle-ci",
								Name: "CircleCI",
								URL:  "https://status.circleci.com/api/v2/status.json",
							},
							Status: statuspageio.Status{
								Indicator:   "major",
								Description: "major outage",
							},
						},
					},
					"minor": {
						statuspageio.Response{
							Page: statuspageio.Page{
								ID:   "code-climate",
								Name: "CodeClimate",
								URL:  "https://status.codeclimate.com/api/v2/status.json",
							},
							Status: statuspageio.Status{
								Indicator:   "minor",
								Description: "minor outage",
							},
						},
					},
				},
				Errors: []string{},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
		"base path- highest severity is minor": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: "statuspage.io",
				},
				"CircleCI": {
					URL:  *circleciURL,
					Type: "statuspage.io",
				},
			},
			reader: clientReaderSeverityMinor{},
			expectedOverview: status.Overview{
				OverallStatus: "minor",
				List: map[string][]status.Details{
					"minor": {
						statuspageio.Response{
							Page: statuspageio.Page{
								ID:   "circle-ci",
								Name: "CircleCI",
								URL:  "https://status.circleci.com/api/v2/status.json",
							},
							Status: statuspageio.Status{
								Indicator:   "minor",
								Description: "minor outage",
							},
						},
					},
					"none": {
						statuspageio.Response{
							Page: statuspageio.Page{
								ID:   "code-climate",
								Name: "CodeClimate",
								URL:  "https://status.codeclimate.com/api/v2/status.json",
							},
							Status: statuspageio.Status{
								Indicator:   "none",
								Description: "none",
							},
						},
					},
				},
				Errors: []string{},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
		"base path- has error fetching a service status": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: "statuspage.io",
				},
				"CircleCI": {
					URL:  *circleciURL,
					Type: "statuspage.io",
				},
			},
			reader: clientReaderHasError{},
			expectedOverview: status.Overview{
				OverallStatus: "none",
				List: map[string][]status.Details{
					"none": {
						statuspageio.Response{
							Page: statuspageio.Page{
								ID:   "code-climate",
								Name: "CodeClimate",
								URL:  "https://status.codeclimate.com/api/v2/status.json",
							},
							Status: statuspageio.Status{
								Indicator:   "none",
								Description: "none",
							},
						},
					},
				},
				Errors: []string{
					"Code: [UNABLE_TO_MAKE_CLIENT_REQUEST] Message: [test err] Inner error: [%!s(<nil>)]",
				},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			o := tc.sites.GetOverview(tc.finder, tc.reader)
			tc.validate(t, tc.expectedOverview, o)
		})
	}
}

type clientReaderSeverityMajor struct {
}

func (crsMajor clientReaderSeverityMajor) ReadStatus(_ client.ServiceFinder, serviceName, _ string) (statuspageio.Response, glitch.DataError) {
	if serviceName == "CircleCI" {
		return statuspageio.Response{
			Page: statuspageio.Page{
				ID:   "circle-ci",
				Name: "CircleCI",
				URL:  "https://status.circleci.com/api/v2/status.json",
			},
			Status: statuspageio.Status{
				Indicator:   "major",
				Description: "major outage",
			},
		}, nil
	}

	return statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "code-climate",
			Name: "CodeClimate",
			URL:  "https://status.codeclimate.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "minor",
			Description: "minor outage",
		},
	}, nil
}

type clientReaderSeverityMinor struct {
}

func (crsMinor clientReaderSeverityMinor) ReadStatus(_ client.ServiceFinder, serviceName, _ string) (statuspageio.Response, glitch.DataError) {
	if serviceName == "CircleCI" {
		return statuspageio.Response{
			Page: statuspageio.Page{
				ID:   "circle-ci",
				Name: "CircleCI",
				URL:  "https://status.circleci.com/api/v2/status.json",
			},
			Status: statuspageio.Status{
				Indicator:   "minor",
				Description: "minor outage",
			},
		}, nil
	}

	return statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "code-climate",
			Name: "CodeClimate",
			URL:  "https://status.codeclimate.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "none",
			Description: "none",
		},
	}, nil
}

type clientReaderHasError struct {
}

func (crhe clientReaderHasError) ReadStatus(_ client.ServiceFinder, serviceName, _ string) (statuspageio.Response, glitch.DataError) {
	if serviceName == "CircleCI" {
		return statuspageio.Response{}, glitch.NewDataError(nil, ErrorUnableToMakeClientRequest, "test err")
	}

	return statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "code-climate",
			Name: "CodeClimate",
			URL:  "https://status.codeclimate.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "none",
			Description: "none",
		},
	}, nil
}

func TestUnit_LoadSites(t *testing.T) {
	codeClimateURL, err := url.Parse("https://status.codeclimate.com/api/v2/status.json")
	require.NoError(t, err)

	tests := map[string]struct {
		reader        configuration.Reader
		writer        configuration.Writer
		filename      string
		expectedSites Sites
		expectedErr   glitch.DataError
		validate      func(t *testing.T, expectedSites, actualSites Sites, expectedErr, actualErr glitch.DataError)
	}{
		"base path- filename given": {
			reader:   fileReaderWithFilename{},
			filename: "test-config.json",
			expectedSites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: "statuspage.io",
				},
			},
			validate: func(t *testing.T, expectedSites, actualSites Sites, expectedErr, actualErr glitch.DataError) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedSites, actualSites)
			},
		},
		"base path- no filename given": {
			reader:        fileReaderWithoutFilename{},
			writer:        fileWriterSuccess{},
			expectedSites: Sites{},
			validate: func(t *testing.T, expectedSites, actualSites Sites, expectedErr, actualErr glitch.DataError) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedSites, actualSites)
			},
		},
		"exceptional path- cannot write configuration file": {
			reader:      fileReaderWithoutFilename{},
			writer:      fileWriterErr{},
			expectedErr: glitch.NewDataError(errors.New("write err"), ErrorUnableToWriteDefaultConfiguration, "unable to create default What's Up configuration"),
			validate: func(t *testing.T, expectedSites, actualSites Sites, expectedErr, actualErr glitch.DataError) {
				require.Error(t, actualErr)
				require.Equal(t, expectedErr.Code(), actualErr.Code())
			},
		},
		"exceptional path- cannot unmarshal configuration file": {
			reader:      fileReaderReturnsInvalidContents{},
			filename:    "test-config.json",
			expectedErr: glitch.NewDataError(nil, ErrorUnableToParseConfiguration, "error parsing What's Up configuration"),
			validate: func(t *testing.T, expectedSites, actualSites Sites, expectedErr, actualErr glitch.DataError) {
				require.Error(t, actualErr)
				require.Equal(t, expectedErr.Code(), actualErr.Code())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, err := LoadSites(tc.reader, tc.writer, tc.filename)
			tc.validate(t, tc.expectedSites, s, tc.expectedErr, err)
		})
	}
}

type fileReaderWithFilename struct {
}

func (fr fileReaderWithFilename) ReadFile(_ string) ([]byte, error) {
	return []byte(`{"CodeClimate":{"url":"https://status.codeclimate.com/api/v2/status.json","type":"statuspage.io"}}`), nil
}

type fileReaderReturnsInvalidContents struct {
}

func (fr fileReaderReturnsInvalidContents) ReadFile(_ string) ([]byte, error) {
	return []byte(`{`), nil
}

type fileReaderWithoutFilename struct {
}

func (fr fileReaderWithoutFilename) ReadFile(_ string) ([]byte, error) {
	return nil, errors.New("read err")
}

type fileWriterSuccess struct {
}

func (fw fileWriterSuccess) WriteFile(_ string, _ []byte, _ fs.FileMode) error {
	return nil
}

type fileWriterErr struct {
}

func (fw fileWriterErr) WriteFile(_ string, _ []byte, _ fs.FileMode) error {
	return errors.New("write err")
}

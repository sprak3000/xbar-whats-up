package service

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/go-whatsup-client/slack"
	whatsupstatus "github.com/sprak3000/go-whatsup-client/status"
	"github.com/sprak3000/go-whatsup-client/statuspageio"
	"github.com/sprak3000/go-whatsup-client/whatsup"
	"github.com/sprak3000/go-whatsup-client/whatsup/clientmock"
	"github.com/stretchr/testify/require"

	"github.com/sprak3000/xbar-whats-up/configuration"
	"github.com/sprak3000/xbar-whats-up/status"
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
				Type: statuspageio.ServiceType,
			},
			validate: func(t *testing.T, expectedSite, actualSite Site, _, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedSite, actualSite)
			},
		},
		"exceptional path- parse URL error": {
			siteJSON:    []byte(`{"url":":","type":"statuspage.io"}`),
			expectedErr: errors.New(`parse ":": missing protocol scheme`),
			validate: func(t *testing.T, _, _ Site, expectedErr, actualErr error) {
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	codeClimateURL, err := url.Parse("https://status.codeclimate.com/api/v2/status.json")
	require.NoError(t, err)

	codeClimateMinorOutageResp := statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "code-climate",
			Name: "CodeClimate",
			URL:  "https://status.codeclimate.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "minor",
			Description: "minor outage",
		},
	}

	codeClimateNoOutageResp := statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "code-climate",
			Name: "CodeClimate",
			URL:  "https://status.codeclimate.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "none",
			Description: "none",
		},
	}

	circleciURL, err := url.Parse("https://status.circleci.com/api/v2/status.json")
	require.NoError(t, err)

	circleciMajorOutageResp := statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "circle-ci",
			Name: "CircleCI",
			URL:  "https://status.circleci.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "major",
			Description: "major outage",
		},
	}

	circleciMinorOutageResp := statuspageio.Response{
		Page: statuspageio.Page{
			ID:   "circle-ci",
			Name: "CircleCI",
			URL:  "https://status.circleci.com/api/v2/status.json",
		},
		Status: statuspageio.Status{
			Indicator:   "minor",
			Description: "minor outage",
		},
	}

	slackURL, err := url.Parse("https://status.slack.com/api/v2.0.0/current")
	require.NoError(t, err)

	slackNoOutageResp := slack.Response{
		Status:          "none",
		DateCreated:     time.Time{},
		DateUpdated:     time.Time{},
		ActiveIncidents: nil,
	}

	tests := map[string]struct {
		sites                 Sites
		setupStatusPageClient func(t *testing.T, expectedErr glitch.DataError) whatsup.StatusPageClient
		expectedOverview      status.Overview
		expectedClientErr     glitch.DataError
		validate              func(t *testing.T, expectedOverview, actualOverview status.Overview)
	}{
		"base path- highest severity is major": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: statuspageio.ServiceType,
				},
				"CircleCI": {
					URL:  *circleciURL,
					Type: statuspageio.ServiceType,
				},
				"Slack": {
					URL:  *slackURL,
					Type: slack.ServiceType,
				},
			},
			setupStatusPageClient: func(_ *testing.T, _ glitch.DataError) whatsup.StatusPageClient {
				c := clientmock.NewMockStatusPageClient(ctrl)
				c.EXPECT().StatuspageIoService("CodeClimate", codeClimateURL.String()).Times(1).Return(codeClimateMinorOutageResp, nil)
				c.EXPECT().StatuspageIoService("CircleCI", circleciURL.String()).Times(1).Return(circleciMajorOutageResp, nil)
				c.EXPECT().Slack().Times(1).Return(slackNoOutageResp, nil)
				return c
			},
			expectedOverview: status.Overview{
				OverallStatus:     "major",
				LargestStringSize: 11,
				List: map[string][]whatsupstatus.Details{
					"major": {
						circleciMajorOutageResp,
					},
					"minor": {
						codeClimateMinorOutageResp,
					},
					"none": {
						slackNoOutageResp,
					},
				},
				Errors: []status.OverviewError{},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
		"base path- highest severity is minor": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: statuspageio.ServiceType,
				},
				"CircleCI": {
					URL:  *circleciURL,
					Type: statuspageio.ServiceType,
				},
			},
			setupStatusPageClient: func(_ *testing.T, _ glitch.DataError) whatsup.StatusPageClient {
				c := clientmock.NewMockStatusPageClient(ctrl)
				c.EXPECT().StatuspageIoService("CodeClimate", codeClimateURL.String()).Times(1).Return(codeClimateNoOutageResp, nil)
				c.EXPECT().StatuspageIoService("CircleCI", circleciURL.String()).Times(1).Return(circleciMinorOutageResp, nil)
				return c
			},
			expectedOverview: status.Overview{
				OverallStatus:     "minor",
				LargestStringSize: 11,
				List: map[string][]whatsupstatus.Details{
					"minor": {
						circleciMinorOutageResp,
					},
					"none": {
						codeClimateNoOutageResp,
					},
				},
				Errors: []status.OverviewError{},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
		"base path- has error fetching a service status": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: statuspageio.ServiceType,
				},
				"CircleCI": {
					URL:  *circleciURL,
					Type: statuspageio.ServiceType,
				},
			},
			setupStatusPageClient: func(_ *testing.T, _ glitch.DataError) whatsup.StatusPageClient {
				c := clientmock.NewMockStatusPageClient(ctrl)
				c.EXPECT().StatuspageIoService("CodeClimate", codeClimateURL.String()).AnyTimes().Return(codeClimateNoOutageResp, nil)
				c.EXPECT().StatuspageIoService("CircleCI", circleciURL.String()).AnyTimes().Return(nil, glitch.NewDataError(nil, whatsupstatus.ErrorUnableToMakeClientRequest, "test err"))
				return c
			},
			expectedOverview: status.Overview{
				OverallStatus:     "none",
				LargestStringSize: 11,
				List: map[string][]whatsupstatus.Details{
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
				Errors: []status.OverviewError{
					{
						Details: nil,
						Error:   glitch.NewDataError(nil, "UNABLE_TO_MAKE_CLIENT_REQUEST", "test err"),
					},
				},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
		"base path- service not supported": {
			sites: Sites{
				"CodeClimate": {
					URL:  *codeClimateURL,
					Type: "not-a-finger",
				},
			},
			setupStatusPageClient: func(_ *testing.T, _ glitch.DataError) whatsup.StatusPageClient {
				c := clientmock.NewMockStatusPageClient(ctrl)
				return c
			},
			expectedOverview: status.Overview{
				OverallStatus: "none",
				List:          map[string][]whatsupstatus.Details{},
				Errors: []status.OverviewError{
					{
						Details: nil,
						Error:   glitch.NewDataError(nil, ErrorUnsupportedServiceType, "CodeClimate uses an unsupported service type not-a-finger"),
					},
				},
			},
			validate: func(t *testing.T, expectedOverview, actualOverview status.Overview) {
				require.Equal(t, expectedOverview, actualOverview)
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			o := tc.sites.GetOverview(tc.setupStatusPageClient(t, tc.expectedClientErr))
			tc.validate(t, tc.expectedOverview, o)
		})
	}
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
					Type: statuspageio.ServiceType,
				},
			},
			validate: func(t *testing.T, expectedSites, actualSites Sites, _, actualErr glitch.DataError) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedSites, actualSites)
			},
		},
		"base path- no filename given": {
			reader:        fileReaderWithoutFilename{},
			writer:        fileWriterSuccess{},
			expectedSites: Sites{},
			validate: func(t *testing.T, expectedSites, actualSites Sites, _, actualErr glitch.DataError) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedSites, actualSites)
			},
		},
		"exceptional path- cannot write configuration file": {
			reader:      fileReaderWithoutFilename{},
			writer:      fileWriterErr{},
			expectedErr: glitch.NewDataError(errors.New("write err"), ErrorUnableToWriteDefaultConfiguration, "unable to create default What's Up configuration"),
			validate: func(t *testing.T, _, _ Sites, expectedErr, actualErr glitch.DataError) {
				require.Error(t, actualErr)
				require.Equal(t, expectedErr.Code(), actualErr.Code())
			},
		},
		"exceptional path- cannot unmarshal configuration file": {
			reader:      fileReaderReturnsInvalidContents{},
			filename:    "test-config.json",
			expectedErr: glitch.NewDataError(nil, ErrorUnableToParseConfiguration, "error parsing What's Up configuration"),
			validate: func(t *testing.T, _, _ Sites, expectedErr, actualErr glitch.DataError) {
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

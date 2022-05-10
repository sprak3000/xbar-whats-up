package service

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_NewClientServiceFinder(t *testing.T) {
	codeClimateURL, err := url.Parse("https://status.codeclimate.com/api/v2/status.json")
	require.NoError(t, err)

	circleciURL, err := url.Parse("https://status.circleci.com/api/v2/status.json")
	require.NoError(t, err)

	tests := map[string]struct {
		sites       Sites
		serviceName string
		expectedURL url.URL
		expectedErr error
		validate    func(t *testing.T, expectedURL, actualURL url.URL, expectedErr, actualErr error)
	}{
		"base path": {
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
			serviceName: "CodeClimate",
			expectedURL: *codeClimateURL,
			validate: func(t *testing.T, expectedURL, actualURL url.URL, expectedErr, actualErr error) {
				require.NoError(t, actualErr)
				require.Equal(t, expectedURL, actualURL)
			},
		},
		"exceptional path- unable to find service": {
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
			serviceName: "GitHub",
			expectedErr: errors.New("unable to find GitHub in the service list"),
			validate: func(t *testing.T, expectedURL, actualURL url.URL, expectedErr, actualErr error) {
				require.Error(t, actualErr)
				require.Equal(t, expectedErr, actualErr)
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := NewClientServiceFinder(tc.sites)
			u, err := f(tc.serviceName, true)
			tc.validate(t, tc.expectedURL, u, tc.expectedErr, err)
		})
	}
}

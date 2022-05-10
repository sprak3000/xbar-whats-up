package service

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprak3000/go-glitch/glitch"
	"github.com/sprak3000/xbar-whats-up/configuration"
)

type fileReaderWithFilename struct {
}

func (fr fileReaderWithFilename) ReadFile(_ string) ([]byte, error) {
	return []byte(`{"CodeClimate":{"url":"https://status.circleci.com","slug":"/api/v2/status.json"}}`), nil
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

func TestUnit_LoadSites(t *testing.T) {
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
					URL:  "https://status.circleci.com",
					Slug: "/api/v2/status.json",
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

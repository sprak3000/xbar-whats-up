package status

import (
	"bytes"
	"testing"
	"time"

	"github.com/sprak3000/go-glitch/glitch"
	whatsupstatus "github.com/sprak3000/go-whatsup-client/status"
	"github.com/stretchr/testify/require"
)

func TestUnit_Overview_Display(t *testing.T) {
	var (
		now          = time.Now()
		nowFormatted = now.Format("2006 Jan 02")
	)

	tests := map[string]struct {
		validate func(t *testing.T)
	}{
		"base path- overall status major": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus:     "major",
					LargestStringSize: 12,
					List: map[string][]whatsupstatus.Details{
						"major": {
							testResponse{
								updatedAt: now,
							},
						},
						"minor": {
							testResponse{
								updatedAt: now,
							},
						},
						"none": {
							testResponse{
								updatedAt: now,
							},
						},
					},
				}

				var buf bytes.Buffer
				o.Display(&buf)

				require.Equal(t, "🔴\n---\n\x1b[31;1mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n---\n\x1b[38;5;208mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n---\n\x1b[32;1mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n", buf.String())
			},
		},
		"base path- overall status minor": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus:     "minor",
					LargestStringSize: 12,
					List: map[string][]whatsupstatus.Details{
						"minor": {
							testResponse{
								updatedAt: now,
							},
						},
						"none": {
							testResponse{
								updatedAt: now,
							},
						},
					},
				}

				var buf bytes.Buffer
				o.Display(&buf)

				require.Equal(t, "🟠\n---\n\x1b[38;5;208mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n---\n\x1b[32;1mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n", buf.String())
			},
		},
		"base path- overall status none": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus:     "none",
					LargestStringSize: 12,
					List: map[string][]whatsupstatus.Details{
						"none": {
							testResponse{
								updatedAt: now,
							},
						},
					},
				}

				var buf bytes.Buffer
				o.Display(&buf)

				require.Equal(t, "🟢\n---\n\x1b[32;1mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n", buf.String())
			},
		},
		"base path- has error": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus:     "none",
					LargestStringSize: 12,
					List: map[string][]whatsupstatus.Details{
						"none": {
							testResponse{
								updatedAt: now,
							},
						},
					},
					Errors: []OverviewError{
						{
							Details: testResponse{updatedAt: now},
							Error:   glitch.NewDataError(nil, "WRONG", "Something went wrong with a test service."),
						},
					},
				}

				var buf bytes.Buffer
				o.Display(&buf)

				require.Equal(t, "🟢\n---\n\x1b[32;1mTest Service     \x1b[0m\x1b[30m "+nowFormatted+" | font=Monaco href=https://test.service/\n---\n\x1b[31;1mTest Service     \x1b[0m\x1b[30m 2025 Apr 26 | font=Monaco href=https://test.service/\n----\nCode: [WRONG] Message: [Something went wrong with a test service.] Inner error: [%!s(<nil>)]", buf.String())
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.validate(t)
		})
	}
}

type testResponse struct {
	updatedAt time.Time
}

func (tr testResponse) Indicator() string {
	return ""
}

func (tr testResponse) Name() string {
	return "Test Service"
}

func (tr testResponse) UpdatedAt() time.Time {
	return tr.updatedAt
}

func (tr testResponse) URL() string {
	return "https://test.service/"
}

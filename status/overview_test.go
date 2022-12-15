package status

import (
	"bytes"
	"testing"
	"time"

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
					OverallStatus: "major",
					List: map[string][]Details{
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

				require.Equal(t, "🔴\n---\n\x1b[31;1mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\n\x1b[38;5;208mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\n\x1b[32;1mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\n", buf.String())
			},
		},
		"base path- overall status minor": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus: "minor",
					List: map[string][]Details{
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

				require.Equal(t, "🟠\n---\n---\n\x1b[38;5;208mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\n\x1b[32;1mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\n", buf.String())
			},
		},
		"base path- overall status none": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus: "none",
					List: map[string][]Details{
						"none": {
							testResponse{
								updatedAt: now,
							},
						},
					},
				}

				var buf bytes.Buffer
				o.Display(&buf)

				require.Equal(t, "🟢\n---\n---\n---\n\x1b[32;1mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\n", buf.String())
			},
		},
		"base path- has error": {
			validate: func(t *testing.T) {
				o := Overview{
					OverallStatus: "none",
					List: map[string][]Details{
						"none": {
							testResponse{
								updatedAt: now,
							},
						},
					},
					Errors: []string{
						"Something went wrong with a test service.",
					},
				}

				var buf bytes.Buffer
				o.Display(&buf)

				require.Equal(t, "🟢\n---\n---\n---\n\x1b[32;1mTest Service\x1b[0m\x1b[30m ("+nowFormatted+") | href=https://test.service/\n---\nSomething went wrong with a test service.\n", buf.String())
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

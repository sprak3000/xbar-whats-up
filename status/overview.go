// Package status is an abstraction for handling and displaying status details from various services
package status

import (
	"fmt"
	"io"

	"github.com/sprak3000/go-glitch/glitch"
	whatsupstatus "github.com/sprak3000/go-whatsup-client/status"
)

// List is a mapping of status codes to services reporting that status code
type List map[string][]whatsupstatus.Details

// OverviewError bundles the details of a failed overview request
type OverviewError struct {
	Details whatsupstatus.Details
	Error   glitch.DataError
}

// Overview provides an overall status for all services monitored -- most severe status wins -- along with all the
// services categorized by status
type Overview struct {
	OverallStatus     string
	LargestStringSize int
	List
	Errors []OverviewError
}

// Display outputs the data in the xbar format
func (o Overview) Display(w io.Writer) {
	switch o.OverallStatus {
	case "major":
		_, _ = fmt.Fprintln(w, "🔴")
	case "minor":
		_, _ = fmt.Fprintln(w, "🟠")
	default:
		_, _ = fmt.Fprintln(w, "🟢")
	}

	displayDetails(w, o.LargestStringSize, o.List["major"], "\u001B[31;1m")
	displayDetails(w, o.LargestStringSize, o.List["minor"], "\u001b[38;5;208m")
	displayDetails(w, o.LargestStringSize, o.List["none"], "\u001B[32;1m")

	if len(o.Errors) > 0 {
		_, _ = fmt.Fprintln(w, "---")
		for _, v := range o.Errors {
			_, _ = fmt.Fprintf(w, "%s%-*s%s%s %s | font=Monaco href=%s\n", "\u001B[31;1m", o.LargestStringSize+5, v.Details.Name(), "\u001b[0m", "\u001b[30m", v.Details.UpdatedAt().Format("2006 Jan 02"), v.Details.URL())
			_, _ = fmt.Fprintln(w, "----")
			_, _ = fmt.Fprintf(w, "%s", v.Error.Error())
		}
	}
}

func displayDetails(w io.Writer, largestStringSize int, details []whatsupstatus.Details, detailColor string) {
	if len(details) > 0 {
		_, _ = fmt.Fprintln(w, "---")
		for _, v := range details {
			_, _ = fmt.Fprintf(w, "%s%-*s%s%s %s | font=Monaco href=%s\n", detailColor, largestStringSize+5, v.Name(), "\u001b[0m", "\u001b[30m", v.UpdatedAt().Format("2006 Jan 02"), v.URL())
		}
	}
}

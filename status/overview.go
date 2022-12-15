// Package status is an abstraction for handling and displaying status details from various services
package status

import (
	"fmt"
	"io"
)

// List is a mapping of status codes to services reporting that status code
type List map[string][]Details

// Overview provides an overall status for all services monitored -- most severe status wins -- along with all the
// services categorized by status
type Overview struct {
	OverallStatus string
	List
	Errors []string
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

	displayDetails(w, o.List["major"], "\u001B[31;1m")
	displayDetails(w, o.List["minor"], "\u001b[38;5;208m")
	displayDetails(w, o.List["none"], "\u001B[32;1m")

	_, _ = fmt.Fprintln(w, "---")
	if len(o.Errors) > 0 {
		for _, v := range o.Errors {
			_, _ = fmt.Fprintln(w, v)
		}
	}
}

func displayDetails(w io.Writer, details []Details, detailColor string) {
	_, _ = fmt.Fprintln(w, "---")
	if len(details) > 0 {
		for _, v := range details {
			_, _ = fmt.Fprintln(w, detailColor+v.Name()+"\u001b[0m"+"\u001b[30m"+" ("+v.UpdatedAt().Format("2006 Jan 02")+") | href="+v.URL())
		}
	}
}

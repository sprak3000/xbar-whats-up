// Package status is an abstraction for handling and displaying status details from various services
package status

import (
	"fmt"
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
func (o Overview) Display() {
	switch o.OverallStatus {
	case "major":
		fmt.Println("🔴")
	case "minor":
		fmt.Println("🟠")
	default:
		fmt.Println("🟢")
	}

	fmt.Println("---")
	if len(o.List["major"]) > 0 {
		for _, v := range o.List["major"] {
			fmt.Println("\u001B[31;1m" + v.Name() + "\u001b[0m" + "\u001b[30m" + " (" + v.UpdatedAt().Format("2006 Jan 02") + ") | href=" + v.URL())
		}
	}

	fmt.Println("---")
	if len(o.List["minor"]) > 0 {
		for _, v := range o.List["minor"] {
			fmt.Println("\u001b[38;5;208m" + v.Name() + "\u001b[0m" + "\u001b[30m" + " (" + v.UpdatedAt().Format("2006 Jan 02") + ") | href=" + v.URL())
		}
	}

	fmt.Println("---")
	if len(o.List["none"]) > 0 {
		for _, v := range o.List["none"] {
			fmt.Println("\u001B[32;1m" + v.Name() + "\u001b[0m" + "\u001b[30m" + " (" + v.UpdatedAt().Format("2006 Jan 02") + ") | href=" + v.URL())
		}
	}

	fmt.Println("---")
	if len(o.Errors) > 0 {
		for _, v := range o.Errors {
			fmt.Println(v)
		}
	}
}

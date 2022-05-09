// <xbar.title>What's Up?</xbar.title>
// <xbar.version>v1.0</xbar.version>
// <xbar.author>Luis A. Cruz</xbar.author>
// <xbar.author.github>sprak3000</xbar.author.github>
// <xbar.desc>Tracks if services are reporting any outages.</xbar.desc>
// <xbar.dependencies>golang</xbar.dependencies>

package main

import (
	"fmt"
	"os"

	"github.com/sprak3000/whats-up/service"
)

func main() {
	sites, lErr := service.LoadSites("./.whats-up.json")
	if lErr != nil {
		fmt.Println("What's Up Error")
		fmt.Println("---")
		fmt.Printf("%v\n", lErr.Error())
		os.Exit(1)
	}

	sites.GetOverview().Display()

	//overview := sites.GetOverview()
	//
	//fmt.Println(overview.OverallStatus)
	//fmt.Println("---")
	//if len(overview.List["major"]) > 0 {
	//	for _, v := range overview.List["major"] {
	//		fmt.Println("\u001B[31;1m" + v.Page.Name + "\u001b[0m" + "\u001b[30m" + " (" + v.Page.UpdatedAt.Format("2006 Jan 02") + ") | href=" + v.Page.URL)
	//	}
	//}
	//fmt.Println("---")
	//if len(overview.List["minor"]) > 0 {
	//	for _, v := range overview.List["minor"] {
	//		fmt.Println("\u001b[38;5;208m" + v.Page.Name + "\u001b[0m" + "\u001b[30m" + " (" + v.Page.UpdatedAt.Format("2006 Jan 02") + ") | href=" + v.Page.URL)
	//	}
	//}
	//fmt.Println("---")
	//if len(overview.List["none"]) > 0 {
	//	for _, v := range overview.List["none"] {
	//		fmt.Println("\u001B[32;1m" + v.Page.Name + "\u001b[0m" + "\u001b[30m" + " (" + v.Page.UpdatedAt.Format("2006 Jan 02") + ") | href=" + v.Page.URL)
	//	}
	//}
}

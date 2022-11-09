// <xbar.title>What's Up?</xbar.title>
// <xbar.version>v0.12.0</xbar.version>
// <xbar.author>Luis A. Cruz</xbar.author>
// <xbar.author.github>sprak3000</xbar.author.github>
// <xbar.desc>Tracks if services are reporting any outages.</xbar.desc>
// <xbar.dependencies>golang</xbar.dependencies>

// Package main is the entry point for running the plugin
package main

import (
	"fmt"
	"os"

	"github.com/sprak3000/xbar-whats-up/configuration"
	"github.com/sprak3000/xbar-whats-up/service"
)

func main() {
	sites, lErr := service.LoadSites(configuration.FileReader{}, configuration.FileWriter{}, "./.whats-up.json")
	if lErr != nil {
		fmt.Println("What's Up Error")
		fmt.Println("---")
		fmt.Printf("%v\n", lErr.Error())
		os.Exit(1)
	}

	sites.GetOverview(service.NewClientServiceFinder(sites), service.NewReaderServiceFinder()).Display()
}

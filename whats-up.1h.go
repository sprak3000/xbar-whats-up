// <xbar.title>What's Up?</xbar.title>
// <xbar.version>v1.0</xbar.version>
// <xbar.author>Luis A. Cruz</xbar.author>
// <xbar.author.github>sprak3000</xbar.author.github>
// <xbar.desc>Tracks if services are reporting any outages.</xbar.desc>
// <xbar.dependencies>golang</xbar.dependencies>

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/sprak3000/go-client/client"
	"github.com/sprak3000/whats-up/service"
	"github.com/sprak3000/whats-up/statuspageio"
)

func main() {
	config, rErr := ioutil.ReadFile("./whats-up.json")
	if rErr != nil {
		// Create an empty configuration file
		config = []byte("{\n}")
		wErr := ioutil.WriteFile("./whats-up.json", config, 0644)
		if wErr != nil {
			fmt.Println("Unable to create default What's Up configuration")
			fmt.Println("---")
			fmt.Printf("Error: %v\n", wErr)
			os.Exit(1)
		}
	}

	var sites service.Sites
	uErr := json.Unmarshal(config, &sites)
	if uErr != nil {
		fmt.Println("Error parsing What's Up configuration")
		fmt.Println("---")
		fmt.Printf("Error: %v\n", uErr)
		os.Exit(1)
	}

	s := statuspageio.Response{}

	f := func(serviceName string, useTLS bool) (url.URL, error) {
		u, err := url.Parse(sites[serviceName].URL)
		return *u, err
	}

	statusMap := map[string][]statuspageio.Response{}
	overallStatus := "🟢"

	for k, v := range sites {
		c := client.NewBaseClient(f, k, true, 10*time.Second, nil)
		_, respBytes, err := c.MakeRequest(context.Background(), "GET", v.Slug, nil, nil, nil)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		uErr := json.Unmarshal(respBytes, &s)
		if uErr != nil {
			fmt.Printf("Error: %v\n", uErr)
			os.Exit(1)
		}

		switch s.Status.Indicator {
		case "major":
			overallStatus = "🔴"
			statusMap["major"] = append(statusMap["major"], s)
		case "minor":
			overallStatus = "🟠"
			statusMap["minor"] = append(statusMap["minor"], s)
		default:
			statusMap["none"] = append(statusMap["none"], s)
		}

	}

	fmt.Println(overallStatus)
	fmt.Println("---")
	if len(statusMap["major"]) > 0 {
		for _, v := range statusMap["major"] {
			fmt.Println("\u001B[31;1m" + v.Page.Name + "\u001b[0m" + "\u001b[30m" + " (" + v.Page.UpdatedAt.Format("2006 Jan 02") + ") | href=" + v.Page.URL)
		}
	}
	if len(statusMap["minor"]) > 0 {
		for _, v := range statusMap["minor"] {
			fmt.Println("\u001b[38;5;208m" + v.Page.Name + "\u001b[0m" + "\u001b[30m" + " (" + v.Page.UpdatedAt.Format("2006 Jan 02") + ") | href=" + v.Page.URL)
		}
	}
	if len(statusMap["none"]) > 0 {
		for _, v := range statusMap["none"] {
			fmt.Println("\u001B[32;1m" + v.Page.Name + "\u001b[0m" + "\u001b[30m" + " (" + v.Page.UpdatedAt.Format("2006 Jan 02") + ") | href=" + v.Page.URL)
		}
	}
}

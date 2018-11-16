/******************************************************************************
* Author:            jbreyer
* Date:              15 NOV 2018
* Purpose:           Presents a hashmap of dotgov.gov registrations as a web
*                    service.  Registrations are updated every 14 days,
*                    coinciding with the registrations being updated.
*
* License:           GPLv2
******************************************************************************/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type registration struct {
	domainName   string
	domainType   string
	agency       string
	organization string
	city         string
	state        string
	isStateLcl   bool
	createdDate  time.Time
	lastUpdate   time.Time
}

func main() {
	url := "https://raw.githubusercontent.com/GSA/data/master/dotgov-domains/current-full.csv"

	registrations := processCSV(fetchList(url))

	fmt.Println(registrations[strings.ToUpper("lbl.gov")].isStateLcl)

}

func fetchList(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	return string(body)
}

func processCSV(file string) map[string]registration {

	registrations := make(map[string]registration)

	//remove leading/trailing whitespaces from CSV and split by lines
	lines := strings.Split(strings.TrimSpace(file), "\n")

	for i := 1; i < len(lines); i++ {
		reg := strings.Split(lines[i], ",")

		//most organizations aren't compounds
		if len(reg) == 6 {
			registrations[string(reg[0])] = registration{
				string(reg[0]),
				string(reg[1]),
				string(reg[2]),
				string(reg[3]),
				string(reg[4]),
				string(reg[5]),
				string(reg[2]) == "Non-Federal Agency",
				time.Now(),
				time.Now(),
			}
		}

		//orgs that are compounds are wrapped in double quotes but strings.split doesn't care about that
		//so we have to separately split the string on double quotes
		if len(reg) != 6 {
			org := strings.Split(lines[i], "\"")
			registrations[string(reg[0])] = registration{
				string(reg[0]),
				string(reg[1]),
				string(reg[2]),
				org[1],
				string(reg[len(reg)-2]),
				string(reg[len(reg)-1]),
				string(reg[2]) == "Non-Federal Agency",
				time.Now(),
				time.Now(),
			}
		}

	}

	return registrations
}

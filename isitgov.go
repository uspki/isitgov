/******************************************************************************
* Author:            jbreyer
* Date:              15 NOV 2018
* Purpose:           Presents a map of dotgov.gov registrations as a web
*                    service.  Registrations are updated every 14 days,
*                    matching the frequency of registration file updates.
******************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type registration struct {
	DomainName   string
	DomainType   string
	Agency       string
	Organization string
	City         string
	State        string
	IsStateLcl   bool
	CreatedDate  time.Time
	LastUpdate   time.Time
}

var regs = make(map[string]registration)

//timers for book keeping only - will be replaced by prometheus hooks
var lastUpdate = time.Time{}
var nextUpdate = time.Time{}

func main() {
	url := "https://raw.githubusercontent.com/GSA/data/master/dotgov-domains/current-full.csv"

	//run goroutine to poll and process registrations every 14 days
	go pollDotGov(url)

	//set up REST endpoints
	router := mux.NewRouter()
	router.HandleFunc("/registrations", getRegs).Methods("GET")
	router.HandleFunc("/registrations/{domain}", getRegDomain).Methods("GET")
	router.HandleFunc("/isStateLocal/{domain}", isStateLocal).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func pollDotGov(url string) {
	for {
		regs = processCSV(fetchList(url))
		lastUpdate = time.Now()
		nextUpdate = lastUpdate.AddDate(0, 0, 14)
		time.Sleep(24 * 14 * time.Hour)
	}
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

	//remove leading/trailing whitespaces and CRs from CSV and split by lines
	lines := strings.Split(strings.TrimSpace(strings.Replace(file, "\r", "", -1)), "\n")

	for i := 1; i < len(lines); i++ {
		reg := strings.Split(lines[i], ",")
		var regUpdate = time.Time{}

		//check if we've seen this registration before and if so retain the CreatedDate
		if regs[reg[0]].CreatedDate.IsZero() {
			regUpdate = time.Now()
		} else {
			regUpdate = regs[reg[0]].CreatedDate
		}

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
				regUpdate,
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
				regUpdate,
				time.Now(),
			}
		}
	}

	return registrations
}

func isStateLocal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(regs[strings.ToUpper(params["domain"])].IsStateLcl)
}

func getRegs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(regs)
}

func getRegDomain(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(regs[strings.ToUpper(params["domain"])])
}

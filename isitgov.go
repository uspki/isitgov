/************************************************************
* Author: jbreyer
* Date:

*/

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
	domainName		string
	domainType      string
	agency          string
	organization    string
	city            string
	state           string
	createdDate     time.Time
	lastUpdate      time.Time
}

func main(){
	url := "https://raw.githubusercontent.com/GSA/data/master/dotgov-domains/current-full.csv"

	resp, err := http.Get(url)
	if err != nil{
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil{
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

    lines := strings.Split(string(body), "\n")

    for i, e := range lines {
    	if i > 0{
    		splitLines := strings.Split(e, ",")

    		if splitLines[2] != "Non-Federal Agency" {
				fmt.Println(splitLines)
			}
		}
	}
}


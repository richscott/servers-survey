package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type surveyResult struct {
	site           string
	serverSoftware string
}

func surveyWorker(id int, siteJobs <-chan string, results chan<- surveyResult) {
	serverSW := "<unknown>"

	timeout := time.Duration(5 * time.Second)
	httpClient := http.Client{
		Timeout: timeout,
	}

	for hostname := range siteJobs {
		// fmt.Printf("Worker %3d = surveying %s\n", id, hostname)

		// TODO set an aggressive short timeout on this
		resp, err := httpClient.Get(fmt.Sprintf("http://%s", hostname))
		if err != nil {
			log.Print(err)
			results <- surveyResult{site: hostname, serverSoftware: serverSW}
			return
		}
		defer resp.Body.Close()

		if len(resp.Header.Get("Server")) > 0 {
			serverSW = resp.Header.Get("Server")
		}
		results <- surveyResult{site: hostname, serverSoftware: serverSW}
	}
}

func main() {
	siteListFile, err := os.Open("fortune500-2014.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer siteListFile.Close()

	sites := make([]string, 0)
	r := csv.NewReader(siteListFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		siteDomain := record[2]
		sites = append(sites, siteDomain)
	}

	siteServer := make(map[string]string)
	jobs := make(chan string, len(sites))
	results := make(chan surveyResult, len(sites))

	// Start up the pool of workers.
	// They will block until they have something to do,
	// which arrives to them via the jobs channel.
	for w := 1; w <= 50; w++ {
		go surveyWorker(w, jobs, results)
	}

	for _, site := range sites {
		jobs <- site
	}

	for n := 1; n <= len(sites); n++ {
		res := <-results
		siteServer[res.site] = res.serverSoftware
	}

	serverCount := make(map[string]int)

	for _, software := range siteServer {
		if count, exists := serverCount[software]; exists {
			serverCount[software] = count + 1
		} else {
			serverCount[software] = 1
		}
	}

	for software, count := range serverCount {
		fmt.Printf("%3d  %-50s\n", count, software)
	}

}

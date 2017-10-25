package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type surveyResult struct {
	site           string
	serverSoftware string
}

var surveyNumber = 0

func surveyWorker(id int, siteJobs <-chan string, results chan<- surveyResult) {
	serverSW := "<unknown>"

	timeout := time.Duration(2 * time.Second)
	httpClient := http.Client{Timeout: timeout}

	for hostname := range siteJobs {
		surveyNumber += 1
		fmt.Printf("Worker %3d survey %3d: surveying %s\n", id, surveyNumber, hostname)

		resp, err := httpClient.Get(fmt.Sprintf("http://%s", hostname))
		if err != nil {
			log.Print(err)
			results <- surveyResult{site: hostname, serverSoftware: serverSW}
			continue
		}
		defer resp.Body.Close()

		if len(resp.Header.Get("Server")) > 0 {
			serverSW = resp.Header.Get("Server")
		}
		results <- surveyResult{site: hostname, serverSoftware: serverSW}
	}
}

func readTabDelimList(tabListFile string) []string {

	file, err := os.Open(tabListFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sites := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "#") {
			siteDomain := strings.Split(scanner.Text(), "\t")[6]
			sites = append(sites, siteDomain)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return sites
}

func readCSVlist(csvFile string) []string {
	siteListFile, err := os.Open(csvFile)
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

	return sites
}

func main() {
	// sites := readCSVlist("fortune500-2014.csv")
	sites := readTabDelimList("fortune1000_companies.tab")

	siteServer := make(map[string]string)
	jobs := make(chan string, len(sites))
	results := make(chan surveyResult, len(sites))

	for _, site := range sites {
		jobs <- site
	}

	// Start up the pool of workers. They will block until they have
	// something to do, which arrives to them via the jobs channel.
	for w := 1; w <= 30; w++ {
		go surveyWorker(w, jobs, results)
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

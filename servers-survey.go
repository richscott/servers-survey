package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
)

func surveySite(siteDomain string, results chan<- string) {
  serverSoftware := "<unknown>"

  resp, err := http.Get(fmt.Sprintf("http://%s", siteDomain))
  if err != nil {
    log.Print(err)
    results <- serverSoftware
    return
  }
  defer resp.Body.Close()

  if len(resp.Header.Get("Server")) > 0 {
    serverSoftware = resp.Header.Get("Server")
  }
  fmt.Printf("%s -> %s\n", siteDomain, serverSoftware)
  results <- serverSoftware
}

func main() {
  siteListFile, err := os.Open("fortune500-2014.csv")
  if err != nil {
    log.Fatal(err)
  }
  defer siteListFile.Close()

	sites := make(map[string]string)

	serverChan := make(chan string)

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

    go surveySite(siteDomain, serverChan)
    serverSoftware := <-serverChan

    sites[siteDomain] = serverSoftware
  }

  //for site, software := range sites {
  //  fmt.Printf("%s -> %s\n", site, software)
  //}
}

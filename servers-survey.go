package main

import (
  "bufio"
  "encoding/csv"
  "fmt"
  "log"
  "os"
)

// Here's the worker, of which we'll run several concurrent instances.
// These workers will receive work on the `jobs` channel and send the
// corresponding results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func surveySite(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("surveySite", id, "processing job", j)
		results <- j * 2
	}
}

func main() {
  siteListFile, err := os.Open("fortune500-2014.csv")
  if err != nil {
    log.Fatal(err)
  }
  defer siteListFile.Close()

  scanner := bufio.NewScanner(siteListFile)
  for scanner.Scan() {
    lineStr := scanner.Text()
    fmt.Println(lineStr)
  }

/*

  // In order to use our pool of surveySites we need to send them work and
  // collect their results. We make 2 channels for this.
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// This starts up 3 surveySites, initially blocked because there are no
  // jobs yet.
	for w := 1; w <= 3; w++ {
		go surveySite(w, jobs, results)
	}

	// Here we send 9 `jobs` and then `close` that channel to indicate
  // that's all the work we have.
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 1; a <= 9; a++ {
		<-results
	}
*/
}

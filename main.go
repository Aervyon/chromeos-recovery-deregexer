package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"
)

func main() {
	fmt.Println("Opening file recovery2.json")
	file, err := os.Open("recovery2.json")

	if err != nil {
		fmt.Println("Error occurred")
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	var recovery []Recovery

	json.Unmarshal(bytes, &recovery)
	now := time.Now()
	start := regexp.MustCompile(`^\^|^\(`)
	end := regexp.MustCompile(`-.{0,2}.$|\..*|\\.*$| \(.*$|\s$`)
	end2 := regexp.MustCompile(`\?$|\[.*$`)
	fuckparens := regexp.MustCompile(`^\(`)

	for i := 0; i < len(recovery); i++ {
		recovery[i].HwIDMatch = start.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end2.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = fuckparens.ReplaceAllString(recovery[i].HwIDMatch, "")
	}

	fmt.Println("Modified", len(recovery), "chromebook models in", (time.Since(now)))

	// TURN INTO CSV
	csvFile, err := os.Create("recovery2-deregexd.csv")
	if err != nil {
		fmt.Println("Ran into error. See error below.")
		fmt.Println(err)
	}
	defer csvFile.Close()

	// hwidmatch manufacturer md5 model url chrome_version
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	header := []string{"hwidmatch", "manufacturer", "md5", "model", "url", "chrome_version"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Ran into error. See error below.")
		println(err)
	}

	for _, r := range recovery {
		var csvRow []string
		csvRow = append(csvRow, r.HwIDMatch, r.Manufacturer, r.MD5, r.Model, r.URL, r.Chromeversion)
		if err := writer.Write(csvRow); err != nil {
			fmt.Println("Ran into error. See error below")
			fmt.Println(err)
		}
	}

	fmt.Println("Finished in", time.Since(now))
}

type Recovery struct {
	HwIDMatch     string `json:"hwidmatch"`
	Model         string `json:"model"`
	URL           string `json:"url"`
	Chromeversion string `json:"chrome_version"`
	Manufacturer  string `json:"manufacturer"`
	Version       string `json:"version"`
	MD5           string `json:"md5"`
}

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
	now := time.Now()
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

	// spaces := regexp.MustCompile(`\s.*$`)
	start := regexp.MustCompile(`^\^|^\(`)
	end := regexp.MustCompile(`-.{0,2}.$|\..*|\\.*$| \(.*$|\s$`)
	end2 := regexp.MustCompile(`\?$|\[.*$`)
	fuckparens := regexp.MustCompile(`^\(`)

	json.Unmarshal(bytes, &recovery)
	for i := 0; i < len(recovery); i++ {
		// recovery[i].HwIDMatch = spaces.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = start.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end2.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = fuckparens.ReplaceAllString(recovery[i].HwIDMatch, "")
		fmt.Println(recovery[i].HwIDMatch)
	}

	fmt.Println("Modified", len(recovery), "models in", (time.Since(now)))

	// TURN INTO CSV
	csvFile, err := os.Create("recovery2-deregexd.csv")
	if err != nil {
		fmt.Println("Failed to create CSV. Will continue to save to JSON")
	}
	defer csvFile.Close()

	// hwidmatch manufacturer md5 model url chrome_version
	if err == nil {
		writer := csv.NewWriter(csvFile)
		defer writer.Flush()

		header := []string{"hwidmatch", "manufacturer", "md5", "model", "url", "chrome_version"}
		if err := writer.Write(header); err != nil {
			fmt.Println("Fuckin failed to write CSV header. 10/10 would give head again. See below")
			fmt.Println(err)
		} else {
			for _, r := range recovery {
				var csvRow []string
				csvRow = append(csvRow, r.HwIDMatch, r.Manufacturer, r.MD5, r.Model, r.URL, r.Chromeversion)
				if err := writer.Write(csvRow); err != nil {
					fmt.Println("Encountered Error. See below")
					fmt.Println(err)
				}
			}
		}
	}

	fmt.Println("Saving to recovery2-deregexd.json")

	endBytes, err := json.MarshalIndent(recovery, "", "  ")
	if err != nil {
		fmt.Println("Ran into an error while marshalling. Will print error and close.")
		fmt.Println(err)
		os.Exit(1)
	}

	os.WriteFile("recovery2-deregex.json", endBytes, 0666)
	fmt.Println("Saving to recovery2-deregexd.csv")
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

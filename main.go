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
	// spaces := regexp.MustCompile(`\s.*$`)
	start := regexp.MustCompile(`^\^|^\(`)
	end := regexp.MustCompile(`-.{0,2}.$|\..*|\\.*$| \(.*$|\s$`)
	end2 := regexp.MustCompile(`\?$|\[.*$`)
	fuckparens := regexp.MustCompile(`^\(`)

	for i := 0; i < len(recovery); i++ {
		// recovery[i].HwIDMatch = spaces.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = start.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end2.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = fuckparens.ReplaceAllString(recovery[i].HwIDMatch, "")
		// fmt.Println(recovery[i].HwIDMatch)
	}

	fmt.Println("Modified", len(recovery), "chromebook models in", (time.Since(now)))

	success, err := toCSV(recovery, "recovery2-deregexd.csv")

	if err != nil {
		fmt.Println("Something went wrong while converting to CSV. See error below. (Will try to save as JSON)")
		fmt.Println(err)
	}
	if success {
		fmt.Println("Saved to recovery2-deregexd.csv")
	}

	fmt.Println("Saving to recovery2-deregexd.json")

	endBytes, err := json.MarshalIndent(recovery, "", "  ")
	if err != nil {
		fmt.Println("Ran into an error while marshalling. Will print error and close.")
		fmt.Println(err)
		os.Exit(1)
	}

	os.WriteFile("recovery2-deregex.json", endBytes, 0666)
	fmt.Println("Finished in", time.Since(now))
}

func toCSV(data []Recovery, destination string) (success bool, err error) {
	// TURN INTO CSV
	csvFile, err := os.Create(destination)
	if err != nil {
		return false, fmt.Errorf("failed to create CSV.")
	}
	defer csvFile.Close()

	// hwidmatch manufacturer md5 model url chrome_version
	if err == nil {
		writer := csv.NewWriter(csvFile)
		defer writer.Flush()

		header := []string{"hwidmatch", "manufacturer", "md5", "model", "url", "chrome_version"}
		if err := writer.Write(header); err != nil {
			return false, err
		}

		for _, r := range data {
			var csvRow []string
			csvRow = append(csvRow, r.HwIDMatch, r.Manufacturer, r.MD5, r.Model, r.URL, r.Chromeversion)
			if err := writer.Write(csvRow); err != nil {
				return false, err
			}
		}
	}
	return true, nil
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

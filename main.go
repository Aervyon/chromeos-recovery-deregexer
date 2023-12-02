package main

import (
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

	spaces := regexp.MustCompile(`\s.*$`)
	start := regexp.MustCompile(`^\^|^\(`)
	end := regexp.MustCompile(`\..*$|-.*$`)
	end2 := regexp.MustCompile(`\?$|\[$`)

	json.Unmarshal(bytes, &recovery)
	for i := 0; i < len(recovery); i++ {
		recovery[i].HwIDMatch = spaces.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = start.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end.ReplaceAllString(recovery[i].HwIDMatch, "")
		recovery[i].HwIDMatch = end2.ReplaceAllString(recovery[i].HwIDMatch, "")
		fmt.Println(recovery[i].HwIDMatch)
	}

	fmt.Println("Modified", len(recovery), "models in", (time.Since(now)))
	fmt.Println("Saving to recovery2-deregexd.json")

	endBytes, err := json.MarshalIndent(recovery, "", "  ")
	if err != nil {
		fmt.Println("Ran into an error while marshalling. Will print error and close.")
		fmt.Println(err)
		os.Exit(1)
	}

	os.WriteFile("recovery2-deregex.json", endBytes, 0666)
}

type Recovery struct {
	SHA1          string `json:"sha1"`
	HwIDMatch     string `json:"hwidmatch"`
	Model         string `json:"model"`
	URL           string `json:"url"`
	Chromeversion string `json:"chrome_version"`
	Manufacturer  string `json:"manufacturer"`
	Version       string `json:"version"`
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"scraper/pkg/github"
)

func main() {
	username := flag.String("u", "", "The username to get issues and PRs")
	flag.Parse()

	if *username == "" {
		fmt.Println("Username cannot be empty")
		os.Exit(1)
	}

	issues, prs, err := github.SearchIssuesAndPrs(*username)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}

	response := map[string]interface{}{
		"issues": issues,
		"prs":    prs,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonResponse))
}

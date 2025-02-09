package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"scraper/pkg/github"
	"text/template"
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
		"Issues": issues,
		"Prs":    prs,
	}

	prsTmpl := template.Must(template.ParseFiles("templates/prs.html"))
	issuesTmpl := template.Must(template.ParseFiles("templates/issues.html"))

	var prsRendered bytes.Buffer
	err = prsTmpl.Execute(&prsRendered, response)
	if err != nil {
		log.Fatalf("Error rendering PRs template: %v", err)
	}

	var issuesRendered bytes.Buffer
	err = issuesTmpl.Execute(&issuesRendered, response)
	if err != nil {
		log.Fatalf("Error rendering Issues template: %v", err)
	}

	data := map[string]interface{}{
		"Issues": issuesRendered.String(),
		"Prs":    prsRendered.String(),
	}

	var baseRendered bytes.Buffer
	tmpl := template.Must(template.ParseFiles("templates/base.html"))
	err = tmpl.Execute(&baseRendered, data)
	if err != nil {
		log.Fatalf("Error rendering base template: %v", err)
	}

	// Github Markdown wasn't rendering HTML tags with empty lines in between
	content := removeEmptyLines(&baseRendered)
	err = os.WriteFile("../README.md", content, 0644)
	if err != nil {
		log.Fatalf("Error writing to README.md: %v", err)
	}

	fmt.Println("README.md updated successfully")
}

func removeEmptyLines(buffer *bytes.Buffer) []byte {
	input := buffer.Bytes()
	lines := bytes.Split(input, []byte("\n"))
	var nonEmptyLines [][]byte
	for _, line := range lines {
		if len(bytes.TrimSpace(line)) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	output := bytes.Join(nonEmptyLines, []byte("\n"))
	return output
}

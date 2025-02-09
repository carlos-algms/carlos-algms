package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Issue struct {
	Title          string `json:"title"`
	URL            string `json:"html_url"`
	RepositoryURL  string `json:"repository_url"`
	RepositoryPath string
	Number         int    `json:"number"`
	State          string `json:"state"`
	Reason         string `json:"state_reason"`
}

func SearchIssuesAndPrs(username string) ([]Issue, []Issue, error) {
	var wg sync.WaitGroup
	var issues, prs interface{}
	var issuesErr, prsErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		issues, issuesErr = SearchGitHub(username, IssueType)
	}()

	go func() {
		defer wg.Done()
		prs, prsErr = SearchGitHub(username, PullRequestType)
	}()

	wg.Wait()

	if issuesErr != nil {
		return nil, nil, errors.New("failed to fetch issues")
	}

	if prsErr != nil {
		return nil, nil, errors.New("failed to fetch pull requests")
	}

	return issues.([]Issue), prs.([]Issue), nil
}

type SearchType string

const (
	IssueType       SearchType = "issue"
	PullRequestType SearchType = "pr"
)

func SearchGitHub(username string, searchType SearchType) ([]Issue, error) {
	if searchType != IssueType && searchType != PullRequestType {
		return nil, fmt.Errorf("invalid search type: %s", searchType)
	}

	filters := []string{
		// "is:open",
		fmt.Sprintf("author:%s", username),
		fmt.Sprintf("is:%s", searchType),
		fmt.Sprintf("-org:%s", username),
		fmt.Sprintf("-user:%s", username),
		fmt.Sprintf("-user:%s", "webdev-tools"), // I don't want my own repos
		fmt.Sprintf("-user:%s", "talesprates"),
	}

	url := fmt.Sprintf(
		"https://api.github.com/search/issues?q=%s&sort=created&order=desc&per_page=10",
		strings.Join(filters, "+"),
	)

	log.Printf("Fetching %s from: %s", searchType, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status code %d: %s, body: %s", resp.StatusCode, resp.Status, string(body))
	}

	var result struct {
		Items []Issue `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`/pull/.*`)
	for i := range result.Items {
		result.Items[i].RepositoryURL = re.ReplaceAllString(result.Items[i].URL, "")
		result.Items[i].RepositoryPath = strings.ReplaceAll(result.Items[i].RepositoryURL, "https://github.com/", "")
	}

	log.Printf("Fetched %d %ss", len(result.Items), searchType)

	return result.Items, nil
}

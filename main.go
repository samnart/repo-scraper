// SPDX-FileCopyrightText: 2025 Awesome
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Repository struct {
	ID			int			`json:"id"`
	Name		string		`json:"name"`
	FullName	string		`json:"full_name"`
	HTMLURL		string		`json:"html_url"`
	CloneURL	string		`json:"clone_url"`
	SSHURL		string		`json:"ssh_url"`
	Description	string		`json:"description"`
	Language	string		`json:"language"`
	Stars		int			`json:"stargazers_count"`
	Forks		int			`json:"forks_count"`
	CreatedAt	string		`json:"created_at"`
	UpdatedAt	string		`json:"updated_at"`
	Private		bool		`json:"private"`
	Fork		bool		`json:"fork"`
	Archived	bool		`json:"archived"`
}

type GitHubScraper struct {
	token		string
	baseUrl		string
	httpClient	*http.Client
}

func NewGitHubScraper(token string) *GitHubScraper {
	return &GitHubScraper{
		token: 		token,
		baseUrl: 	"https://api.github.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

//perform http request with auth
func (gs *GitHubScraper) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return  nil, err
	}

	if gs.token != "" {
		req.Header.Set("Authorization", "Bearer " + gs.token)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "GitHub-Repo-Scraper/1.0")

	return gs.httpClient.Do(req)
}

//scrape all repos for a user
func (gs *GitHubScraper) ScrapeUserRepos(username string) ([]Repository, error) {
	var allRepos []Repository
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s/users/%s/repos?page=%d&per_page=%d&sort=updated&direction=desc",
			gs.baseUrl, username, page, perPage)

		fmt.Printf("Fetching page %d...\n", page)
		repos, hasMore, err := gs.fetchReposPage(url)
		if err != nil {
			return  nil, err
		}

		fmt.Printf("Found %d repositories on page %d\n", len(repos), page)
		allRepos = append(allRepos, repos...)

		if !hasMore || len(repos) == 0 {
			break
		}
		page++
	}

	return allRepos, nil
}

func (gs *GitHubScraper) ScraperOrgRepos(orgName string) ([]Repository, error) {
	var allRepos []Repository
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s/orgs/%s/repos?page=%d&per_page=%dsort=updated&direction=desc",
			gs.baseUrl, orgName, page, perPage)

		fmt.Printf("Fetching page %d...\n", page)
		repos, hasMore, err := gs.fetchReposPage(url)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Found %d repositories on page %d\n", len(repos), page)
		allRepos = append(allRepos, repos...)

		if !hasMore || len(repos) == 0 {
			break
		}
		page++
	}

	return allRepos, nil
}

func (gs *GitHubScraper) fetchReposPage(url string) ([]Repository, bool, error) {
	resp, err := gs.makeRequest(url)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, false, err
	}

	//check for more apges
	linkHeader := resp.Header.Get("Link")
	// hasMore := strings.Contains(linkHeader, `rel=next`)
	hasMore := parseLinkHeader(linkHeader)

	return  repos, hasMore, nil
}

func parseLinkHeader(linkHeader string) bool {
	if linkHeader == "" {
		return false
	}

	links := strings.Split(linkHeader, ",")

	for _, link := range links {
		link = strings.TrimSpace(link)

		if strings.Contains(link, `rel="next"`) {
			return true
		}
	}

	return  false
}

func FilterRepos(repos []Repository, includePrivate, includeForks, includeArchived bool) []Repository {
	var filtered []Repository

	for _, repo := range repos {
		if !includePrivate && repo.Private {
			continue
		}
		if !includeForks && repo.Fork {
			continue
		}
		if !includeArchived && repo.Archived {
			continue
		}
		filtered = append(filtered, repo)
	}

	return filtered
}

func SaveToJson(repos []Repository, filename string) error {
	data, err := json.MarshalIndent(repos, "", " ")
	if err != nil {
		return  err
	}

	return os.WriteFile(filename, data, 0644)
}

func SaveToCSV(repos []Repository, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return  err
	}
	defer file.Close()

	header := "Name,Full Name,Description,Language,Stars,Focks,Clone URL,HTML URL,Created At,Updated At,Private,Fork,Archived\n"
	if _, err := file.WriteString(header); err != nil {
		return  err
	}

	for _, repo := range repos {
		row := fmt.Sprintf("%s,%s,\"%s\",%s,%d,%d,%s,%s,%s,%s,%t,%t,%t\n",
			repo.Name, repo.FullName, strings.ReplaceAll(repo.Description, "\"", "\"\""), repo.Language,
			repo.Stars, repo.Forks, repo.CloneURL, repo.HTMLURL, repo.CreatedAt, repo.UpdatedAt, repo.Private,
			repo.Fork, repo.Archived)

		if _, err := file.WriteString(row); err != nil {
			return err
		}
	}

	return nil
}

func PrintSummary(repos []Repository, targetName string) {
	fmt.Printf("\n=== Repository Summary for %s ===\n", targetName)
	fmt.Printf("Total repositories found: %d\n", len(repos))

	if len(repos) == 0 {
		return
	}

	languages := make(map[string]int)
	totalStars := 0
	totalForks := 0
	privateCount := 0
	forkCount := 0
	archivedCount := 0

	for _, repo := range repos {
		if repo.Language != "" {
			languages[repo.Language]++
		}
		totalStars += repo.Stars
		totalForks += repo.Forks
		if repo.Private {
			privateCount++
		}
		if repo.Fork {
			forkCount++
		}
		if repo.Archived {
			archivedCount++
		}
	}

	fmt.Printf("Total stars: %d\n", totalStars)
	fmt.Printf("Total forks: %d\n", totalForks)
	fmt.Printf("Private repositories: %d\n", privateCount)
	fmt.Printf("Forked repositories: %d\n", forkCount)
	fmt.Printf("Archived repositories: %d\n", archivedCount)

	fmt.Printf("\nTop languages:\n")
	for lang, count := range languages {
		fmt.Printf("	%s: %d repositories\n", lang, count)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <user|org> <name> [options]")
		fmt.Println("	user|org: Specify whether to scrape a user or organization")
		fmt.Println("	name: GitHub username or organization name")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("	--token <token>: GitHub personal access token (recommended)")
		fmt.Println("	--output <format>: Output format (json|csv|both) [default: json]")
		fmt.Println("	--include-private: Include private repositories")
		fmt.Println("	--include-forked: Include forked repositories")
		fmt.Println("	--include-archived: Include archived repositories")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("	go run main.go org noi-techpark")
		fmt.Println("	go run main.go user samnart1 --token ghp_xxxx --output both")
		fmt.Println()
		os.Exit(1)
	}

	targetType := os.Args[1]
	targetName := os.Args[2]

	var token string
	outputFormat := "json"
	includePrivate := false
	includeArchived := false
	includeForks := true

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--token":
			if i+1 < len(os.Args) {
				token = os.Args[i+1]
				i++
			}

		case "--ouput":
			if i+1 < len(os.Args) {
				outputFormat = os.Args[i+1]
				i++
			}

		case "--include-private":
			includePrivate = true

		case "--include-forks":
			includeForks = false

		case "--include-archived":
			includeArchived = true
		}
	}

	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	scraper := NewGitHubScraper(token)

	fmt.Printf("Scraping repositories for %s: %s...\n", targetType, targetName)

	var repos []Repository
	var err error

	if targetType == "org" {
		repos, err = scraper.ScraperOrgRepos(targetName)
	} else if targetType == "user" {
		repos, err = scraper.ScrapeUserRepos(targetName)
	} else {
		fmt.Println("Error: First argument must be 'user' or 'org'")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error scraping repositories: %v\n", err)
		os.Exit(1)
	}

	repos = FilterRepos(repos, includePrivate, includeForks, includeArchived)

	PrintSummary(repos, targetName)

	timestamp := time.Now()
	baseFilename := fmt.Sprintf("%s_%s_%s", targetType, targetName, timestamp)

	switch outputFormat {
	case "json":
		filename := baseFilename + ".json"
		if err := SaveToJson(repos, filename); err != nil {
			fmt.Printf("Error saving JSON: %v\n", err)
		} else {
			fmt.Printf("\nRepositories saved to %s\n", filename)
		}
	case "csv":
		filename := baseFilename + ".csv"
		if err := SaveToCSV(repos, filename); err != nil {
			fmt.Printf("Error saving CSV: %v\n", err)
		} else {
			fmt.Printf("\nRepositories saved to %s\n", filename)
		}
	case "both":
		jsonFile := baseFilename + ".json"
		csvFile := baseFilename + ".csv"
		
		if err := SaveToJson(repos, jsonFile); err != nil {
			fmt.Printf("Error saving to JSON: %v\n", err)
		} else {
			fmt.Printf("\nRepositories saved to %s\n", jsonFile)
		}
		
		if err := SaveToCSV(repos, csvFile); err != nil {
			fmt.Printf("Error saving to CSV: %v\n", err)
		} else {
			fmt.Printf("\nRepositories saved to %s\n", csvFile)
		}
		
	default:
		fmt.Printf("Error: unknown output format '%s'. Use 'json', 'csv', or 'both'\n", outputFormat)
	}
}
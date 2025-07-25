# GitHub Repository Scraper

[![REUSE status](https://api.reuse.software/badge/github.com/samnart/repo-scraper)](https://api.reuse.software/info/github.com/samnart/repo-scraper)
A flexible and powerful command-line tool written in GoLang to scrape public repositories from GitHub organizations and users. Originally built for Open Data Hub but works with any GitHub account (org, user).

## Features
- **Flexible target:** Works with both orgs and users.
- **Auth support:** Optional access token for higher rate limits
- **Pagination handling:** Automatically fetches all repos across muliple pages.
- **Multi output format:** export data as json, csv or both
- **Rich metadata:** captures comprehensive repo info like stars, forks, languages and timestamps (createdAt, last updated)

## Quick Start
### Prerequisities
- Go 1.19 or higher
- Optional: github access token

### Installation
1. Clone the repo:
```
git clone https:/github.com/samnart/repo-scraper.git
cd repo-scraper
```

2. Initialize module:
```
go mod init repo-scraper
go mod tidy
```

3. Build the application:
```
go build -o repo-scraper main.go
```

4. Basic Usage
```
#scrape Open Data Hub organization
./repo-scraper org noi-techpark

#scrape a user's repos
./repo-scraper user samnart1

#with auth token
./repo-scraper org noi-techpark --token ghp_your_special_token

#include all repo types and save both as csv and json
./repo-scraper org noi-techpark --include-archived --output both
```

## Support

If you encounter any problems or have suggestions, please [open an issue](https://github.com/samnart/repo-scraper/issues/new).

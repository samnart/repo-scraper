# Repository Scraper

[![REUSE status](https://api.reuse.software/badge/github.com/samnart/repo-scraper)](https://api.reuse.software/info/github.com/samnart/repo-scraper)

A powerful Go-based tool that scrapes all public repositories from any GitHub organization. Built with performance and ease of use in mind, this tool helps you analyze and collect information about GitHub organizations' public repositories efficiently.

## Features

- 🚀 Fast and concurrent repository scraping
- 📦 Fetch all public repositories from any GitHub organization
- 🔄 Rate limit handling and automatic retries
- 💾 Multiple output formats support (JSON, CSV - coming soon)
- 🛠 Easy-to-use command line interface
- ✨ Written in Go for high performance

## Installation

### Prerequisites

- Go 1.16 or higher
- GitHub account (for API access)

### Installing from source

```bash
# Clone the repository
git clone https://github.com/samnart/repo-scraper.git

# Change to the project directory
cd repo-scraper

# Build the project
go build

# Or install directly using go get
go get github.com/samnart/repo-scraper
```

## Usage

### Basic Usage

```bash
# Basic usage with organization name
repo-scraper -org="organization-name"

# Specify output format (coming soon)
repo-scraper -org="organization-name" -output="json"

# Filter repositories by language (coming soon)
repo-scraper -org="organization-name" -filter="go"
```

### Environment Variables

You can configure the tool using the following environment variables:

- `GITHUB_TOKEN`: Your GitHub personal access token (recommended for higher rate limits)
- `OUTPUT_DIR`: Directory where the results will be saved (optional)

### Output

The tool generates output containing information about each repository, including:
- Repository name
- Description
- Stars count
- Fork count
- Last update time
- Primary language
- Topics
- And more...

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Make your changes
4. Run tests and ensure they pass
5. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
6. Push to the branch (`git push origin feature/AmazingFeature`)
7. Open a Pull Request

Please make sure to update tests as appropriate and follow the existing code style.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
All files are REUSE compliant with appropriate licensing and copyright information.

## Acknowledgments

- Thanks to the Go community for the excellent standard library
- GitHub API for making this tool possible
- All contributors who help improve this project

## Support

If you encounter any problems or have suggestions, please [open an issue](https://github.com/samnart/repo-scraper/issues/new).
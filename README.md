# GitCrawler

Go REST API that clones GitHub repositories and extracts file data for ingestion by the RetrospectiveAI service.

## How it works

When a request arrives, the API clones the target repository locally, scans it for files matching the requested extensions and directories, and returns the results. For business model summary requests, it extracts relevant files and sends them to the LLM, returning the generated summary.

## Tech Stack

- **Go 1.22**
- Standard library `net/http`
- Static binary built on Alpine Linux

## Running locally

Requirements: Go 1.22+, `git` installed.

```bash
cd app
go mod download
go build -o gitcrawler ./main/
./gitcrawler
```

Server starts on **port 8080**.

## Environment variables

| Variable | Description |
|---|---|
| `API_KEY` | LLM API key (used for business model summaries) |
| `AI_RESUME_PROMPT` | Custom prompt for repository summaries (optional) |

## API Endpoints

### `POST /getRepoData`
Clones a repository and returns a list of matching files.

```bash
curl -X POST http://localhost:8080/getRepoData \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/org/repo.git",
    "dirs": ["src", "pkg"],
    "extensions": [".go", ".md"]
  }'
```

### `POST /saveRepoData`
Clones a repository and exports matching files to disk.

```bash
curl -X POST http://localhost:8080/saveRepoData \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/org/repo.git",
    "dirs": ["src"],
    "extensions": [".go"],
    "option": "json"
  }'
```

### `GET /getBusinessRepoResume?url=<repo-url>`
Returns an AI-generated summary of the repository's purpose and tech stack.

```bash
curl "http://localhost:8080/getBusinessRepoResume?url=https://github.com/org/repo.git"
```

## Running with Docker

```bash
# From the root retrospective/ directory
docker compose up gitcrawler
```

## Project structure

```
app/
├── main/
│   └── main.go                  # Entry point and dependency injection
└── impl/
    ├── adapters/
    │   ├── facade/              # RepositoryFacade, AIResumeGenerateFacade
    │   └── register/            # HTTP route registration
    ├── core/
    │   └── service/             # CloneService, CrawlerService, etc.
    └── external/
        ├── integration/         # LLM integration
        └── rest/                # HTTP controller
```

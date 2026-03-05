## GitCrawler

Rest API that trims repositories from GitHub and return only the necessary data or resumes of the repositories for data ingestion

## How It Works

When a request is received, the API first clones the target GitHub repository locally.  
For file extraction requests, it scans the cloned repository and selects files based on the requested extensions and directories.
With the
For business model summary requests, the API extracts relevant information from the repository, such as the README, documentation, or main code files, and sends it to the AI service. The generated summary is then returned directly to the user.

## Diagram

<Project C4 or sequencial diagram>
  
## Tech Stack

Golang

## Use Cases
### 1. Get Files in GitHub Repos
- **Input:** GitHub URLs + file extensions/directories  
- **Output:** List of matching files (path, repo, type)  

### 2. Save Files of GitHub Repos
- **Input:** GitHub URLs + file extensions/packages + option to save (json, csv, etc)
- **Output:** Saves in the downloads directory of your machine and archive with all the matching files, in the format you asked

### 3. Generate Business Model Summary via AI
- **Input:** GitHub URL + API_KEY  
- **Output:** Business model summary

## Set-up

```bash
git clone https://github.com/Gabriel-Gerhardt/GitCrawler.git
cd GitCrawler
docker compose up
```

Access:

    localhost:8080/getRepoData:
    example: curl -X POST http://localhost:8080/getRepoData \
    -H "Content-Type: application/json" \
    -d '{
    "url": "https://github.com/Gabriel-Gerhardt/GitCrawler.git",
    "dirs": ["pkg", "service"],
    "extensions": [".go", ".md"]
    }'

    localhost:8080/saveRepoData:
    example: curl -X POST http://localhost:8080/saveRepoData \
    -H "Content-Type: application/json" \
    -d '{
    "url": "https://github.com/Gabriel-Gerhardt/GitCrawler.git",
    "dirs": ["pkg", "service"],
    "extensions": [".go", ".md"],
    "option": "csv"
    }'

    get -> localhost:8080/getBusinessRepoResume/{repositoryUrl}
    example: curl "localhost:8080/getBusinessRepoResume?url=https://github.com/Gabriel-Gerhardt/GitCrawler.git"

## Contact
[LinkedIn](https://www.linkedin.com/in/gabriel-gerhardt-0a8b852b9/)

[Gmail](mailto:gabrielgerhardt27@gmail.com)

[GitHub](https://github.com/Gabriel-Gerhardt)


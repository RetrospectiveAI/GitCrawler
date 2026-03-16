package rest

import (
	"encoding/json"
	"gitcrawler/app/impl/adapters/dto/request"
	"gitcrawler/app/impl/adapters/dto/response"
	"gitcrawler/app/impl/adapters/facade"
	"net/http"
)

type CrawlerController struct {
	repositoryFacade *facade.RepositoryFacade
}

func NewCrawlerController(repositoryFacade *facade.RepositoryFacade) *CrawlerController {
	return &CrawlerController{repositoryFacade: repositoryFacade}
}

func (c *CrawlerController) GetRepositoryFiles(w http.ResponseWriter, r *http.Request) {
	var req request.GetRepositoryFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Url == "" {
		http.Error(w, "Url must contain something", http.StatusBadRequest)
		return
	}
	data, err := c.repositoryFacade.GetRepositoryFiles(req.Url, req.Extensions, req.Dirs, req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *CrawlerController) SaveRepositoryFile(w http.ResponseWriter, r *http.Request) {
	var req request.RepositoryFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Url == "" {
		http.Error(w, "Url must contain something", http.StatusBadRequest)
		return
	}
	err := c.repositoryFacade.SaveRepositoryFiles(req.Url, req.Extensions, req.Dirs, req.Option)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The csv is ready"))

}

func (c *CrawlerController) GetBusinessRepoResume(w http.ResponseWriter, r *http.Request) {
	var resp response.ResumeResponse
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Url must contain something", http.StatusBadRequest)
		return
	}
	// Optional PAT for private repositories; empty string for public repos.
	token := r.URL.Query().Get("token")

	aiResponse, err := c.repositoryFacade.GenerateBusinessResume(url, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.AiResponse = aiResponse

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.AiResponse))
}

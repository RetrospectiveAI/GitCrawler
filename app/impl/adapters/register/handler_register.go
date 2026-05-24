package register

import (
	"gitcrawler/app/impl/external/rest"
	"net/http"
)

func GetHandlers(crawlerController *rest.CrawlerController, internalApiKey string) {
	auth := func(h http.HandlerFunc) http.Handler {
		return rest.RequireApiKey(internalApiKey, h)
	}
	http.Handle("/getRepoData", auth(crawlerController.GetRepositoryFiles))
	http.Handle("/saveRepoData", auth(crawlerController.SaveRepositoryFile))
	http.Handle("/getBusinessRepoResume", auth(crawlerController.GetBusinessRepoResume))
}

package main

import (
	"fmt"
	"gitcrawler/app/impl/adapters/facade"
	"gitcrawler/app/impl/adapters/register"
	"gitcrawler/app/impl/core/service"
	"gitcrawler/app/impl/external/integration"
	"gitcrawler/app/impl/external/rest"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	llm := integration.NewLlmIntegration(os.Getenv("API_KEY"))
	fileWriter := integration.NewFileWriter()

	resumeService := service.NewResumeGenerateService(llm)
	cloneService := service.NewCloneService()
	crawlerService := service.NewCrawlerService()

	repositoryLoader := service.NewRepositoryLoaderService(crawlerService, cloneService)

	resumeFacade := facade.NewAIResumeGenerateFacade(repositoryLoader, resumeService)
	repoFacade := facade.NewRepositoryFacade(repositoryLoader, fileWriter)
	crawlerController := rest.NewCrawlerController(repoFacade, resumeFacade)

	server := http.Server{}
	register.GetHandlers(crawlerController)

	server.Addr = ":8080"
	fmt.Println("Running server on " + server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

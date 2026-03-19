package service

import (
	"bytes"
	"gitcrawler/app/impl/core/contract"
	"os"
)

type ResumeGenerateService struct {
	llmIntegration contract.LlmIntegrationContract
}

func NewResumeGenerateService(llmIntegration contract.LlmIntegrationContract) *ResumeGenerateService {
	return &ResumeGenerateService{llmIntegration: llmIntegration}
}

func (s *ResumeGenerateService) GenerateBusinessResume(data string) (text string, err error) {
	resume := os.Getenv("AI_RESUME_PROMPT") + data
	jsonResume := escapeJSON(resume)
	response, err := s.llmIntegration.ReturnAIResponse(jsonResume)
	if err != nil {
		return "", err
	}
	return response, nil
}

func escapeJSON(str string) string {
	b := []byte(str)
	b = bytes.ReplaceAll(b, []byte(`\`), []byte(`\\`))
	b = bytes.ReplaceAll(b, []byte(`"`), []byte(`\"`))
	b = bytes.ReplaceAll(b, []byte("\n"), []byte(`\n`))
	return string(b)
}

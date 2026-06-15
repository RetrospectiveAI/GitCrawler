package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const openRouterUrl = "https://openrouter.ai/api/v1/chat/completions"

type LlmIntegration struct {
	key string
}

func NewLlmIntegration(key string) *LlmIntegration {
	return &LlmIntegration{key: key}
}

func (s *LlmIntegration) ReturnAIResponse(prompt string) (aiResponse string, err error) {
	request, err := s.buildRequest(prompt)
	if err != nil {
		return "", err
	}
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return "", errors.New("unable to complete request")

	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return s.getResponse(respBody)
}

func (s *LlmIntegration) buildRequest(prompt string) (request *http.Request, err error) {
	body := []byte(`{
		"model": "openai/gpt-oss-120b",
		"messages": [{"role": "user", "content": "` + prompt + `"}]
	}`)

	request, err = http.NewRequest("POST", openRouterUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("unable to mount request")
	}
	request.Header.Set("Authorization", "Bearer "+s.key)

	return request, nil
}

func (s *LlmIntegration) getResponse(respBody []byte) (aiResponse string, err error) {
	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err = json.Unmarshal(respBody, &aiResp); err != nil {
		return "", errors.New("could not parse AI response JSON")
	}

	if len(aiResp.Choices) == 0 {
		return "", errors.New("AI has no response, the key may be wrong")
	}

	return aiResp.Choices[0].Message.Content, nil
}

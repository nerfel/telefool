package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"telefool/configs"
	"telefool/pkg/memory"
	"time"
)

var ErrUnauthorized = errors.New("401 Unauthorized")

type YandexGetIAMTokenResponse struct {
	IamToken  string    `json:"iamToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type ModelRequestPayload struct {
	ModelUri          string                          `json:"modelUri"`
	CompletionOptions *configs.ModelCompletionOptions `json:"completionOptions"`
	Messages          []*NeuroMessage                 `json:"messages"`
}

type NeuroMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type NeuroResponse struct {
	Result       NeuroResult `json:"result"`
	ModelVersion string      `json:"modelVersion"`
}

type NeuroResult struct {
	Alternatives []ResponseMessage `json:"alternatives"`
}

type ResponseMessage struct {
	Message NeuroMessage `json:"message"`
	Status  string       `json:"status"`
}

func GetYandexGPTOauthToken(config *configs.Config) error {
	postBody, _ := json.Marshal(map[string]string{
		"yandexPassportOauthToken": config.YandexCloudConfig.Token,
	})

	response, err := http.Post(
		config.YandexCloudConfig.GetIamTokenUrl,
		"application/json",
		bytes.NewBuffer(postBody),
	)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status + "_ydx")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var yandexResponse YandexGetIAMTokenResponse
	json.Unmarshal(body, &yandexResponse)

	config.YandexCloudConfig.IamToken = yandexResponse.IamToken

	return nil
}

func BuildModelRequestPayload(messages []memory.Message, config *configs.Config) (*ModelRequestPayload, error) {
	systemMessage := &NeuroMessage{
		Role: "system",
		Text: config.YandexCloudConfig.SystemPrompt,
	}

	neuroMessages := make([]*NeuroMessage, 0, len(messages))

	neuroMessages = append(neuroMessages, systemMessage)

	var role string
	for _, message := range messages {
		if message.FromCurrentBot {
			role = "assistant"
		} else {
			role = "user"
		}
		neuroMessages = append(neuroMessages, &NeuroMessage{
			Role: role,
			Text: message.Text,
		})
	}

	return &ModelRequestPayload{
		ModelUri:          config.YandexCloudConfig.ModelUri,
		CompletionOptions: config.YandexCloudConfig.ModelCompletionOptions,
		Messages:          neuroMessages,
	}, nil
}

func RequestModel(payload *ModelRequestPayload, config *configs.Config) (*NeuroResponse, error) {
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", config.YandexCloudConfig.GptModelRequestUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.YandexCloudConfig.IamToken)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("An error occurred while sending the message history to the yandex gpt")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	var neuroResponse NeuroResponse
	err = json.NewDecoder(response.Body).Decode(&neuroResponse)
	if err != nil {
		log.Println("An error occurred while decoding response Yandex gpt for struct")
		return nil, err
	}

	return &neuroResponse, nil
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"personal-diary/models"
	"strconv"
	"strings"
	"time"
)

// Implement the chat GPT

func RefineTextWithChatGPT(text, context, tone string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	prompt := buildRefinePrompt(text, context, tone)

	chatGPTReq := models.ChatGPTRequest{
		Model: "gpt-4o-mini",
		Messages: []models.Message{
			{Role: "system", Content: "You are a helpful writing assistant..."},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
		MaxTokens:   500,
	}

	jsonData, err := json.Marshal(chatGPTReq)
	if err != nil {
		return "", fmt.Errorf("marshal request error: %v", err)
	}

	var lastErr error
	for i := 0; i < 3; i++ { // Retry up to 3 times
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("create request error: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("send request error: %v", err)
			time.Sleep(time.Duration(2<<i) * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("read response error: %v", err)
			time.Sleep(time.Duration(2<<i) * time.Second)
			continue
		}

		if resp.StatusCode == 429 {
			lastErr = fmt.Errorf("rate limited: %s", string(body))

			// Optional: check for Retry-After header
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					time.Sleep(time.Duration(seconds) * time.Second)
				}
			} else {
				time.Sleep(time.Duration(2<<i) * time.Second) // exponential backoff
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("non-200 response: %d\n%s", resp.StatusCode, string(body))
		}

		var chatGPTResp models.ChatGPTResponse
		if err := json.Unmarshal(body, &chatGPTResp); err != nil {
			return "", fmt.Errorf("parse response error: %v", err)
		}
		if chatGPTResp.Error != nil {
			return "", fmt.Errorf("ChatGPT error: %s", chatGPTResp.Error.Message)
		}
		if len(chatGPTResp.Choices) == 0 {
			return "", fmt.Errorf("empty response from ChatGPT")
		}
		return strings.TrimSpace(chatGPTResp.Choices[0].Message.Content), nil
	}
	return "", lastErr
}

// buildRefinePrompt creates an appropriate prompt
func buildRefinePrompt(text, context, tone string) string {
	var prompt strings.Builder
	prompt.WriteString("Please refine the following text to improve its clarity, grammar, and overall quality")

	switch strings.ToLower(context) {
	case "diary", "diary_entry":
		prompt.WriteString(" for a personal diary entry. Maintain the personal and reflective nature")
	case "email":
		prompt.WriteString(" for an email. Make it clear and professional")
	case "academic":
		prompt.WriteString(" for academic writing. Use formal language and ensure precision")
	case "creative":
		prompt.WriteString(" for creative writing. Enhance the literary quality and flow")
	default:
		prompt.WriteString(" while maintaining the original meaning and style")
	}

	switch strings.ToLower(tone) {
	case "professional":
		prompt.WriteString(", using a professional tone")
	case "casual":
		prompt.WriteString(", keeping a casual and friendly tone")
	case "formal":
		prompt.WriteString(", using formal language")
	case "friendly":
		prompt.WriteString(", maintaining a warm and friendly tone")
	}

	prompt.WriteString(":\n\n")
	prompt.WriteString(text)
	prompt.WriteString("\n\nPlease provide only the refined text without any explanations or additional comments.")

	return prompt.String()
}

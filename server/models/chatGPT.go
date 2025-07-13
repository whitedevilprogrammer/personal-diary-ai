package models


// RefineRequest represents the request structure for text refinement
type RefineRequest struct {
	Text    string `json:"text"`
	Context string `json:"context,omitempty"`
	Tone    string `json:"tone,omitempty"`
}

// RefineResponse represents the response structure
type RefineResponse struct {
	RefinedText string `json:"refinedText"`
	Message     string `json:"message,omitempty"`
}

// Message for ChatGPT
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTRequest format
type ChatGPTRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// Choice from ChatGPT response
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// APIError struct
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// ChatGPTResponse full response
type ChatGPTResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

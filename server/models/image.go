// ======================
// models/gemini_models.go (add these to your existing models)
// ======================

package models

// Enhanced response structure that includes both full path and URL
type BackgroundImageData struct {
	ImageURL    string `json:"imageUrl"`  // URL path for web access
	ImagePath   string `json:"imagePath"` // Full file system path
	Prompt      string `json:"prompt"`
	GeneratedAt string `json:"generatedAt"`
}

type BackgroundImageResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    *BackgroundImageData `json:"data,omitempty"`
}

type BackgroundImageRequest struct {
	Content string `json:"content"`
	Title   string `json:"title,omitempty"`
}

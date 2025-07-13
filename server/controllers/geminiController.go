// ======================
// controllers/geminiController.go (Enhanced version)
// ======================
package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"personal-diary/models"
	"personal-diary/services"
)

type BackgroundImageController struct {
	ImageService *services.ImageGenerationService
}

func NewBackgroundImageController(imageService *services.ImageGenerationService) *BackgroundImageController {
	return &BackgroundImageController{
		ImageService: imageService,
	}
}

func (c *BackgroundImageController) GenerateBackground(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received background generation request")

	ctx := context.Background()

	payload := models.NewPayload()
	result := models.NewResponse()
	var req models.BackgroundImageRequest

	if err := payload.DecodePayload(r, &req); err != nil {
		log.Printf("Error decoding payload: %v", err)
		result.ErrorResponse(w, "Invalid request format")
		return
	}

	log.Printf("Request content: %s, title: %s", req.Content, req.Title)

	if strings.TrimSpace(req.Content) == "" {
		log.Printf("Empty content provided")
		result.ErrorResponse(w, "Content is required")
		return
	}

	if req.Title == "" {
		req.Title = "Diary Entry"
	}

	// Generate image prompt
	log.Printf("Generating image prompt...")
	prompt, err := c.ImageService.GenerateImagePrompt(ctx, req.Content, req.Title)
	if err != nil {
		log.Printf("Error generating prompt: %v", err)
		c.sendErrorResponse(w, "Failed to generate image prompt", http.StatusInternalServerError)
		return
	}

	log.Printf("Generated prompt: %s", prompt)

	// Generate image with both paths
	log.Printf("Generating image...")
	imageResult, err := c.ImageService.GenerateImageWithPaths(ctx, prompt)
	if err != nil {
		log.Printf("Error generating image: %v", err)
		c.sendErrorResponse(w, "Failed to generate background image", http.StatusInternalServerError)
		return
	}

	log.Printf("Generated image file path: %s", imageResult.FilePath)
	log.Printf("Generated image URL: %s", imageResult.URL)

	response := models.BackgroundImageResponse{
		Status: "success",
		Data: &models.BackgroundImageData{
			ImageURL:    imageResult.URL,      // Web-accessible URL
			ImagePath:   imageResult.FilePath, // Full file system path
			Prompt:      prompt,
			GeneratedAt: time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

	log.Printf("Successfully sent response with both paths")
}

func (c *BackgroundImageController) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Sending error response: %s (status: %d)", message, statusCode)

	response := models.BackgroundImageResponse{
		Status:  "error",
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

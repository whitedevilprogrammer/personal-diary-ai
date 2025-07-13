// ======================
// services/image_service.go (Fixed version)
// ======================
package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type ImageGenerationService struct {
	Client    *genai.Client
	UploadDir string
}

// ImageResult contains both file path and URL
type ImageResult struct {
	FilePath string // Full file system path
	URL      string // Web-accessible URL path
	Size     int64  // File size in bytes
}

func NewImageGenerationService(apiKey, uploadDir string) (*ImageGenerationService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Ensure the upload directory exists with proper permissions
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for upload directory: %w", err)
	}

	if err := os.MkdirAll(absUploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Test write permissions
	testFile := filepath.Join(absUploadDir, ".test_write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return nil, fmt.Errorf("upload directory is not writable: %w", err)
	}
	os.Remove(testFile) // Clean up test file

	log.Printf("Upload directory created/verified: %s", absUploadDir)

	return &ImageGenerationService{
		Client:    client,
		UploadDir: absUploadDir,
	}, nil
}

func (s *ImageGenerationService) GenerateImagePrompt(ctx context.Context, content, title string) (string, error) {
	model := s.Client.GenerativeModel("gemini-1.5-flash")

	prompt := fmt.Sprintf(`Based on the following diary entry content, create a detailed, artistic image prompt that would make a beautiful background image for this diary entry.

Guidelines:
- Focus on creating atmospheric, aesthetic backgrounds
- Include artistic styles (watercolor, digital art, photography, etc.)
- Mention colors, lighting, and mood
- Keep it suitable as a subtle background (not too busy)
- Make it emotionally resonant with the diary content
- Include keywords like "soft", "dreamy", "atmospheric" for background suitability

Diary Title: %s
Diary Content: %s

Generate only the image prompt (no explanations):`, title, content)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate image prompt: %w", err)
	}

	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no prompt generated")
	}

	generatedPrompt := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return strings.TrimSpace(generatedPrompt), nil
}

// GenerateImage returns full file system path (backward compatibility)
func (s *ImageGenerationService) GenerateImage(ctx context.Context, prompt string) (string, error) {
	result, err := s.GenerateImageWithPaths(ctx, prompt)
	if err != nil {
		return "", err
	}
	return result.FilePath, nil
}

// GenerateImageWithPaths returns both file path and URL with proper image generation
func (s *ImageGenerationService) GenerateImageWithPaths(ctx context.Context, prompt string) (*ImageResult, error) {
	// Generate unique filename with timestamp
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("bg_%d.jpg", timestamp)
	filePath := filepath.Join(s.UploadDir, filename)

	log.Printf("Generating image file: %s", filePath)
	log.Printf("Using prompt: %s", prompt)

	// Try to generate actual image using an API or create a proper placeholder
	imageData, err := s.generateActualImage(ctx, prompt)
	if err != nil {
		log.Printf("Failed to generate actual image, creating placeholder: %v", err)
		// Create a proper image placeholder instead of text
		imageData, err = s.createPlaceholderImage(prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to create placeholder image: %w", err)
		}
	}

	// Write the image data to file
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return nil, fmt.Errorf("failed to save image file: %w", err)
	}

	// Verify file was created and get its size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("image file was not created successfully: %w", err)
	}

	// Get absolute path
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create URL path for web access
	urlPath := "/uploads/" + filename
	
	log.Printf("Image generated successfully at: %s (size: %d bytes)", absolutePath, fileInfo.Size())
	log.Printf("Web URL: %s", urlPath)
	
	return &ImageResult{
		FilePath: absolutePath,
		URL:      urlPath,
		Size:     fileInfo.Size(),
	}, nil
}

// generateActualImage attempts to generate a real image using an external API
func (s *ImageGenerationService) generateActualImage(ctx context.Context, prompt string) ([]byte, error) {
	// Option 1: Use DALL-E API (if you have OpenAI API key)
	// Replace with your actual image generation API
	
	// Option 2: Use Stability AI API
	// Replace with your actual image generation API
	
	// Option 3: Use any other image generation service
	// For now, return an error to fall back to placeholder
	
	return nil, fmt.Errorf("actual image generation not implemented - add your preferred image API here")
}

// createPlaceholderImage creates a proper image file instead of text
func (s *ImageGenerationService) createPlaceholderImage(prompt string) ([]byte, error) {
	// Create a simple colored image as placeholder
	width, height := 1024, 768
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Create a gradient background based on prompt keywords
	bgColor := s.getColorFromPrompt(prompt)
	
	// Fill with gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a subtle gradient
			alpha := uint8(255 - (y*50/height))
			img.Set(x, y, color.RGBA{bgColor.R, bgColor.G, bgColor.B, alpha})
		}
	}
	
	// Encode as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("failed to encode placeholder image: %w", err)
	}
	
	return buf.Bytes(), nil
}

// getColorFromPrompt returns a color based on prompt keywords
func (s *ImageGenerationService) getColorFromPrompt(prompt string) color.RGBA {
	promptLower := strings.ToLower(prompt)
	
	// Define color themes based on keywords
	switch {
	case strings.Contains(promptLower, "sunset") || strings.Contains(promptLower, "warm"):
		return color.RGBA{255, 165, 0, 255} // Orange
	case strings.Contains(promptLower, "ocean") || strings.Contains(promptLower, "blue"):
		return color.RGBA{70, 130, 180, 255} // Steel Blue
	case strings.Contains(promptLower, "forest") || strings.Contains(promptLower, "green"):
		return color.RGBA{34, 139, 34, 255} // Forest Green
	case strings.Contains(promptLower, "purple") || strings.Contains(promptLower, "lavender"):
		return color.RGBA{147, 112, 219, 255} // Medium Purple
	case strings.Contains(promptLower, "pink") || strings.Contains(promptLower, "rose"):
		return color.RGBA{255, 182, 193, 255} // Light Pink
	default:
		return color.RGBA{173, 216, 230, 255} // Light Blue
	}
}

// DownloadAndSaveImage downloads an image from URL and saves it locally
func (s *ImageGenerationService) DownloadAndSaveImage(ctx context.Context, imageURL string) (*ImageResult, error) {
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Download the image
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Determine file extension from content type
	ext := ".jpg"
	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("bg_%d%s", timestamp, ext)
	filePath := filepath.Join(s.UploadDir, filename)

	// Save the image
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Get absolute path
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create URL path
	urlPath := "/uploads/" + filename

	log.Printf("Image downloaded and saved: %s (size: %d bytes)", absolutePath, fileInfo.Size())

	return &ImageResult{
		FilePath: absolutePath,
		URL:      urlPath,
		Size:     fileInfo.Size(),
	}, nil
}

// ValidateImageFile checks if a file is a valid image
func (s *ImageGenerationService) ValidateImageFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Try to decode as image
	_, _, err = image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("file is not a valid image: %w", err)
	}

	return nil
}

// GetImageDimensions returns the dimensions of an image file
func (s *ImageGenerationService) GetImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image config: %w", err)
	}

	return config.Width, config.Height, nil
}

// ConvertToPNG converts an image to PNG format
func (s *ImageGenerationService) ConvertToPNG(inputPath string) (string, error) {
	// Open source image
	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source image: %w", err)
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Create output filename
	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimSuffix(inputPath, ext) + ".png"

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Encode as PNG
	if err := png.Encode(outFile, img); err != nil {
		return "", fmt.Errorf("failed to encode PNG: %w", err)
	}

	return outputPath, nil
}

// ListGeneratedImages returns a list of all generated images
func (s *ImageGenerationService) ListGeneratedImages() ([]ImageResult, error) {
	var results []ImageResult

	entries, err := os.ReadDir(s.UploadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read upload directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		// Check if it's an image file
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			continue
		}

		filePath := filepath.Join(s.UploadDir, filename)
		info, err := entry.Info()
		if err != nil {
			continue
		}

		results = append(results, ImageResult{
			FilePath: filePath,
			URL:      "/uploads/" + filename,
			Size:     info.Size(),
		})
	}

	return results, nil
}

// CleanupOldImages removes images older than the specified duration
func (s *ImageGenerationService) CleanupOldImages(olderThan time.Duration) error {
	entries, err := os.ReadDir(s.UploadDir)
	if err != nil {
		return fmt.Errorf("failed to read upload directory: %w", err)
	}

	cutoff := time.Now().Add(-olderThan)
	deletedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(s.UploadDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				log.Printf("Failed to delete old image %s: %v", filePath, err)
			} else {
				deletedCount++
				log.Printf("Deleted old image: %s", filePath)
			}
		}
	}

	log.Printf("Cleanup completed: deleted %d old images", deletedCount)
	return nil
}

func (s *ImageGenerationService) Close() error {
	if s.Client != nil {
		return s.Client.Close()
	}
	return nil
}
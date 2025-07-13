// ======================
// routers/geminiRouters.go
// ======================
package routers

import (
	"personal-diary/controllers"
	"personal-diary/services"

	"github.com/gorilla/mux"
)

func GeminiRouters(router *mux.Router, apiKey string, uploadDir string) error {
	// Initialize the image service
	imageService, err := services.NewImageGenerationService(apiKey, uploadDir)
	if err != nil {
		return err
	}

	// Initialize the controller
	imageController := controllers.NewBackgroundImageController(imageService)

	// Create subrouter
	geminiRouter := router.PathPrefix("/gemini").Subrouter()

	// Bind handler
	geminiRouter.HandleFunc("/generate-background", imageController.GenerateBackground).Methods("POST")

	return nil
}

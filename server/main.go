package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath" // Make sure this is imported

	"personal-diary/middleware"
	"personal-diary/routers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// ... (rest of your init)
}

// main function to start the server
// It sets up the router, initializes routes, and serves static files.
func main() {
	// config.ConnectDB()

	r := mux.NewRouter()

	// Setup existing routes
	routers.AuthRouters(r)
	routers.DiaryRouters(r)

	// Define the upload directory relative to the server's execution path
	// This path should point to: your_project_root/personal-diary-frontend/public/uploads
	uploadDir := "../personal-diary-frontend/public/uploads" // MODIFIED

	// Ensure upload directory exists before starting server
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path for upload directory: %v", err)
	}

	if err := os.MkdirAll(absUploadDir, 0755); err != nil { // Use 0755 for directory permissions
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	log.Printf("Upload directory successfully set to: %s", absUploadDir)

	// Initialize Gemini routers
	// The routers.GeminiRouters function and services.NewImageGenerationService
	// will receive this updated absUploadDir (or the relative one, depending on their implementation,
	// but services.NewImageGenerationService already does filepath.Abs on its received uploadDir)
	if err := routers.GeminiRouters(r, os.Getenv("GEMINI_API_KEY"), absUploadDir); err != nil { // Pass absolute path
		log.Fatalf("Failed to initialize Gemini routers: %v", err)
	}

	// Serve static files from the new uploads directory
	// This allows the Go server to serve these files directly if needed.
	// During Vue development, the Vite dev server will also serve files from `public`.
	fileServer := http.FileServer(http.Dir(absUploadDir))

	// You have two handlers for "/uploads/", let's keep the logging one for clarity:
	// r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fileServer)) // This one would be overwritten

	r.PathPrefix("/uploads/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Go server serving static file request: %s from %s", req.URL.Path, absUploadDir)
		http.StripPrefix("/uploads/", fileServer).ServeHTTP(w, req)
	})

	// Wrap with CORS middleware
	handler := middleware.CORSMiddleware(r)

	log.Println("Server running on http://localhost:8080")
	log.Printf("Static files will be served by Go server from: %s via /uploads/", absUploadDir)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

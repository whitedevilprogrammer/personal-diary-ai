package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"personal-diary/config"
	"personal-diary/middleware"
	"personal-diary/models"
	"personal-diary/services"
	"personal-diary/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = config.GetCollection("users")

type TokenRequest struct {
	Token string `json:"token"`
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a more secure state token
	state := utils.GenerateRandomString(32) // You'll need to implement this utility function

	// Store state in session or cache for validation in callback
	// For simplicity, using a basic approach here
	url := services.GetGoogleOAuthConfig().AuthCodeURL(state)

	log.Printf("OAuth Config Client ID: %s", services.GetGoogleOAuthConfig().ClientID)
	log.Printf("Redirecting to Google OAuth: %s", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Validate state parameter for CSRF protection
	if state == "" {
		log.Println("State parameter missing")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	if code == "" {
		log.Println("Authorization code missing")
		http.Error(w, "Authorization code missing", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := services.GetGoogleOAuthConfig().Exchange(ctx, code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	client := services.GetGoogleOAuthConfig().Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Google API returned status: %d", resp.StatusCode)
		http.Error(w, "Failed to get user info from Google", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		http.Error(w, "Failed to read user info", http.StatusInternalServerError)
		return
	}

	var user models.GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("JSON decode failed: %v", err)
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	// Validate required user fields
	if user.Email == "" {
		log.Println("Email not provided by Google")
		http.Error(w, "Email not provided by Google", http.StatusBadRequest)
		return
	}

	// Save or update user in database
	models.LoginWithGoogleOAuth(user)
	// if err != nil {
	// 	log.Printf("Failed to save user: %v", err)
	// 	http.Error(w, "Failed to save user", http.StatusInternalServerError)
	// 	return
	// }

	// Generate JWT token
	tokenString, err := utils.GenerateToken(user.Email)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Redirect to frontend with token and Dynamic Host
	redirectURL := fmt.Sprintf("http://%s:5173/login?token=%s", os.Getenv("HOST"), tokenString)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	payload := models.NewPayload()
	result := models.NewResponse()

	if err := payload.DecodePayload(r, &user); err != nil {
		log.Printf("Invalid request payload: %v", err)
		result.ErrorResponse(w, "Invalid request payload")
		return
	}

	// Sanitize input
	user.Email = utils.SanitizeString(user.Email)

	// Sanitize input
	user.Email = utils.SanitizeString(user.Email)

	// Validate required fields
	if user.Email == "" || user.Password == "" {
		result.ErrorResponse(w, "Email and password are required")
		return
	}

	// Validate email format (basic validation)
	if !utils.IsValidEmail(user.Email) { // You'll need to implement this utility function
		result.ErrorResponse(w, "Invalid email format")
		return
	}

	// Validate password strength
	if isValid, message := utils.ValidatePassword(user.Password); !isValid {
		result.ErrorResponse(w, message)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		result.ErrorResponse(w, "User already exists")
		return
	} else if err != mongo.ErrNoDocuments {
		log.Printf("Database error: %v", err)
		result.ErrorResponse(w, "Database error")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		result.ErrorResponse(w, "Error processing password")
		return
	}
	user.Password = string(hashedPassword)

	// Set creation timestamp
	// user.CreatedAt = time.Now()
	// user.UpdatedAt = time.Now()

	// Insert the new user
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		result.ErrorResponse(w, "Error creating user")
		return
	}

	result.SuccessResponse(w, "User registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	payload := models.NewPayload()
	result := models.NewResponse()

	if err := payload.DecodePayload(r, &user); err != nil {
		log.Printf("Invalid request payload: %v", err)
		result.ErrorResponse(w, "Invalid request payload")
		return
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" {
		result.ErrorResponse(w, "Email and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user by email
	var foundUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			result.ErrorResponse(w, "Before Login, Kindly register your email ID")
		} else {
			log.Printf("Database error: %v", err)
			result.ErrorResponse(w, "Database error")
		}
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		result.ErrorResponse(w, "Invalid Password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(foundUser.Email)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		result.ErrorResponse(w, "Error generating token")
		return
	}

	// // Update last login timestamp
	// updateData := bson.M{
	// 	"$set": bson.M{
	// 		"last_login": time.Now(),
	// 		"updated_at": time.Now(),
	// 	},
	// }
	// _, err = userCollection.UpdateOne(ctx, bson.M{"email": foundUser.Email}, updateData)
	// if err != nil {
	// 	log.Printf("Error updating last login: %v", err)
	// 	// Don't fail the login for this error
	// }

	// Return success response with token
	data := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"email": foundUser.Email,
			"name":  foundUser.Name,
		},
	}

	result.SetData(data)
	result.SuccessResponse(w, "Login successful")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	result := models.NewResponse()

	token := r.Header.Get("Authorization")
	if token == "" {
		result.ErrorResponse(w, "No token provided")
		return
	}

	// Strip "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Validate token format
	if token == "" {
		result.ErrorResponse(w, "Invalid token format")
		return
	}

	// Add token to blacklist
	middleware.BlacklistedTokens[token] = true

	log.Printf("User logged out, token blacklisted")
	result.SuccessResponse(w, "Logout successful")
}

func RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	payload := models.NewPayload()
	result := models.NewResponse()

	// Create a struct to decode JSON payload
	var input struct {
		Email string `json:"email"`
	}

	if err := payload.DecodePayload(r, &input); err != nil {
		log.Printf("Invalid request payload: %v", err)
		result.ErrorResponse(w, "Invalid request payload")
		return
	}

	email := input.Email
	if email == "" {
		result.ErrorResponse(w, "Email is required")
		return
	}

	// Generate reset token
	token, err := utils.GenerateToken(email)
	if err != nil {
		log.Println("Token generation failed:", err)
		result.ErrorResponse(w, "Internal server error")
		return
	}

	// 👉 TODO: Save the token to DB with expiry

	// Send email
	sender := services.NewEmailSender()
	if err := sender.SendPasswordResetEmail(email, token); err != nil {
		log.Println("Email send failed:", err)
		result.ErrorResponse(w, "Failed to send email")
		return
	}

	result.SuccessResponse(w, "Successfully requested password reset")
}

func ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	payload := models.NewPayload()
	result := models.NewResponse()
	var req TokenRequest
	if err := payload.DecodePayload(r, &req); err != nil {
		log.Printf("Invalid request payload: %v", err)
		result.ErrorResponse(w, "Invalid request payload")
		return
	}

	// Check if token is blacklisted
	if middleware.IsBlacklisted(req.Token) {
		http.Error(w, "Token has been revoked", http.StatusUnauthorized)
		return
	}

	// Parse token
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		return middleware.JwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("Token parsing error:", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"status":  "success",
	// 	"message": "Valid reset token",
	// })

	result.SuccessResponse(w, "Valid reset token")
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	payload := models.NewPayload()
	result := models.NewResponse()

	// Create a struct to decode JSON payload
	var input struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := payload.DecodePayload(r, &input); err != nil {
		log.Printf("Invalid request payload: %v", err)
		result.ErrorResponse(w, "Invalid request payload")
		return
	}

	// Validate token
	if middleware.IsBlacklisted(input.Token) {
		result.ErrorResponse(w, "Token has been revoked")
		return
	}

	// Parse the token and extract claims
	token, err := jwt.Parse(input.Token, func(t *jwt.Token) (interface{}, error) {
		return middleware.JwtKey, nil
	})
	if err != nil || !token.Valid {
		log.Println("Token parsing error:", err)
		result.ErrorResponse(w, "Invalid or expired token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["email"] == nil {
		result.ErrorResponse(w, "Invalid token claims")
		return
	}
	email := claims["email"].(string)

	// Validate passwords
	if input.Password == "" || input.ConfirmPassword == "" {
		result.ErrorResponse(w, "Both password fields are required")
		return
	}

	if input.Password != input.ConfirmPassword {
		result.ErrorResponse(w, "Passwords do not match")
		return
	}

	if isValid, msg := utils.ValidatePassword(input.Password); !isValid {
		result.ErrorResponse(w, msg)
		return
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		result.ErrorResponse(w, "Error processing password")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the user's password
	updateData := bson.M{
		"$set": bson.M{
			"password": hashedPassword,
		},
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"email": email}, updateData)
	if err != nil {
		log.Printf("Error updating new password: %v", err)
		result.ErrorResponse(w, "Failed to update password")
		return
	}

	// Blacklist the token after use
	middleware.RevokeToken(input.Token)

	result.SuccessResponse(w, "Password has been reset successfully")
}

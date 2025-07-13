package controllers

import (
	"context"
	"fmt"
	"net/http"
	"personal-diary/config"
	"personal-diary/models"
	"personal-diary/services"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var diaryCollection *mongo.Collection = config.GetCollection("diaries")

func getEmailFromHeader(r *http.Request) string {
	tokenStr := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	return claims["email"].(string)
}

func CreateDiary(w http.ResponseWriter, r *http.Request) {
	var entry models.DiaryEntry
	payload := models.NewPayload()
	result := models.NewResponse()
	if err := payload.DecodePayload(r, &entry); err != nil {
		result.ErrorResponse(w, "Invalid request payload")
		return
	}
	// json.NewDecoder(r.Body).Decode(&entry)

	entry.ID = primitive.NewObjectID().Hex()
	entry.CreatedAt = time.Now()
	entry.Email = getEmailFromHeader(r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := diaryCollection.InsertOne(ctx, entry)
	if err != nil {
		result.ErrorResponse(w, "Failed to create diary entry")
		// http.Error(w, "Failed to create diary entry", http.StatusInternalServerError)
		return
	}

	result.SuccessResponse(w, "Diary entry created successfully")
	// json.NewEncoder(w).Encode(entry)

}

func GetAllDiaries(w http.ResponseWriter, r *http.Request) {
	email := getEmailFromHeader(r)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := diaryCollection.Find(ctx, bson.M{"email": email}, options.Find().SetSort(bson.D{{"createdAt", -1}}))
	if err != nil {
		http.Error(w, "Failed to fetch diary entries", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var entries []models.DiaryEntry
	for cursor.Next(ctx) {
		var entry models.DiaryEntry
		if err := cursor.Decode(&entry); err != nil {
			http.Error(w, "Failed to decode diary entry", http.StatusInternalServerError)
			return
		}
		entries = append(entries, entry)
	}

	result := models.NewResponse()
	result.SetData(entries)
	result.SuccessResponse(w, "Diary entries fetched successfully")
	// json.NewEncoder(w).Encode(entries)
}

func UpdateDiary(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updateData models.DiaryEntry
	payload := models.NewPayload()
	result := models.NewResponse()
	if err := payload.DecodePayload(r, &updateData); err != nil {
		result.ErrorResponse(w, "Invalid request payload")
		return
	}
	// json.NewDecoder(r.Body).Decode(&updateData)

	// Validate the ID
	if id == "" {
		result.ErrorResponse(w, "ID is required")
		// http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "email": getEmailFromHeader(r)}
	update := bson.M{
		"$set": bson.M{
			"title":   updateData.Title,
			"content": updateData.Content,
		},
	}

	_, err := diaryCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		result.ErrorResponse(w, "Failed to update diary entry")
		// http.Error(w, "Failed to update diary entry", http.StatusInternalServerError)
		return
	}

	// json.NewEncoder(w).Encode(entry)
	result.SuccessResponse(w, "Diary entry updated successfully")
	// json.NewEncoder(w).Encode(map[string]string{"message": "Updated"})
}

func DeleteDiary(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result := models.NewResponse()
	// Validate the ID
	if id == "" {
		result.ErrorResponse(w, "ID is required")
		// http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := diaryCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		result.ErrorResponse(w, "Failed to delete diary entry")
		// http.Error(w, "Failed to delete diary entry", http.StatusInternalServerError)
		return
	}

	result.SuccessResponse(w, "Diary entry deleted successfully")
	// json.NewEncoder(w).Encode(map[string]string{"message": "Deleted"})
}

// RefineTextHandler handles text refinement requests
func RefineTextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := models.NewPayload()
	result := models.NewResponse()
	var req models.RefineRequest

	if err := payload.DecodePayload(r, &req); err != nil {
		result.ErrorResponse(w, "IInvalid JSON format")
		return
	}

	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, `{"status":"error","message":"Invalid JSON format"}`, http.StatusBadRequest)
	// 	return
	// }

	// Validate
	text := strings.TrimSpace(req.Text)
	if len(text) == 0 {
		result.ErrorResponse(w, "Text cannot be empty")
		// http.Error(w, `{"status":"error","message":"Text cannot be empty"}`, http.StatusBadRequest)
		return
	}
	if len(text) > 5000 {
		result.ErrorResponse(w, "Text is too long. Maximum 5000 characters allowed.")
		// http.Error(w, `{"status":"error","message":"Text is too long. Maximum 5000 characters allowed."}`, http.StatusBadRequest)
		return
	}

	// Process
	refinedText, err := services.RefineTextWithChatGPT(text, req.Context, req.Tone)
	if err != nil {
		result.ErrorResponse(w, fmt.Sprintf(`Failed to refine text,error:"%s"`, err.Error()))
		// http.Error(w, fmt.Sprintf(`{"status":"error","message":"Failed to refine text","error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"status": "success",
	// 	"data": RefineResponse{
	// 		RefinedText: refinedText,
	// 		Message:     "Text refined successfully",
	// 	},
	// })

	result.SetData(models.RefineResponse{RefinedText: refinedText, Message: "Text refined successfully"})
	result.SuccessResponse(w, "Text refined successfully")
}

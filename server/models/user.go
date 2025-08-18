package models

import (
	"context"
	"log"
	"personal-diary/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type GoogleUser struct {
	ID      string `json:"id" bson:"id"`
	Email   string `json:"email" bson:"email"`
	Name    string `json:"name" bson:"name"`
	Picture string `json:"picture" bson:"picture"`
}

func SaveOrUpdateUser(user GoogleUser) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": user}
	opts := options.Update().SetUpsert(true)

	userCollection.UpdateOne(ctx, filter, update, opts)
}

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) EmailExists() bool {
	userCollection := config.GetCollection(u.CollectionName())
	filter := bson.M{"email": u.Email}

	var result bson.M
	err := userCollection.FindOne(context.TODO(), filter).Decode(&result)
	return err == nil // true if found, false if not found
}

func (u *GoogleUser) CollectionName() string {
	return "users"
}

func LoginWithGoogleOAuth(user GoogleUser) {
	newUser := User{
		Name:  user.Name,
		Email: user.Email,
	}

	var update bson.M

	if !newUser.EmailExists() {
		// Hash the Google ID as the "password" substitute
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.ID), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return
		}
		newUser.Password = string(hashedPassword)

		// Include password in update
		update = bson.M{"$set": bson.M{
			"name":     newUser.Name,
			"email":    newUser.Email,
			"password": newUser.Password,
		}}
	} else {
		// Do not update password
		update = bson.M{"$set": bson.M{
			"name":  newUser.Name,
			"email": newUser.Email,
		}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": newUser.Email}
	opts := options.Update().SetUpsert(true)
	userCollection := config.GetCollection(newUser.CollectionName())
	_, err := userCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error updating/inserting user: %v", err)
	}
}

func init() {
	// Reuse the DB connection from config package
	userCollection = config.DB.Database("oauthdb").Collection("users")
}

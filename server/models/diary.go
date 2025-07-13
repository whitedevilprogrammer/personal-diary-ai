package models

import "time"

type DiaryEntry struct {
	ID        string    `json:"_id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Email     string    `json:"email" bson:"email"` // owner
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Thumbnail   Thumbnail          `json:"thumbnail,omitempty"`
	PublishedAt string             `json:"publishedAt,omitempty"`
}

type Thumbnail struct {
	Default string `json:"default,omitempty"`
	Medium  string `json:"medium,omitempty"`
	High    string `json:"high,omitempty"`
}

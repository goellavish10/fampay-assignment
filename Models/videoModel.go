package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	Id          primitive.ObjectID `json:"_id,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Thumbnail   Thumbnail          `json:"thumbnail,omitempty"`
	PublishedAt string             `json:"publishedAt,omitempty"`
	VideoId     string             `json:"videoId,omitempty"`
}

type Thumbnail struct {
	Default ThumbnailProperties `json:"default,omitempty"`
	Medium  ThumbnailProperties `json:"medium,omitempty"`
	High    ThumbnailProperties `json:"high,omitempty"`
}

type ThumbnailProperties struct {
	Url    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

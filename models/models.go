package models

import "time"

// Struct to represent a YouTube video
type Video struct {
	Id          string    `json: "id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Thumbnails  []string  `json:"thumbnails"`
	URL         string    `json:"url"`
}

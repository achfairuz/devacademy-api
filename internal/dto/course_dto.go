package dto

import "github.com/google/uuid"

type CreateCourseRequest struct {
	MentorID    uuid.UUID `json:"mentor_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail" binding:"omitempty,url"`
	Price       float64   `json:"price"`
	Level       string    `json:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Duration    int       `json:"duration"`
	Status      string    `json:"status" binding:"omitempty,oneof=draft published"`
}

type UpdateCourseRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail" binding:"omitempty,url"`
	Price       float64 `json:"price"`
	Level       string  `json:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Duration    int     `json:"duration"`
	Status      string  `json:"status" binding:"omitempty,oneof=draft published"`
}

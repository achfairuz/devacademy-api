package course

import "github.com/google/uuid"

type CreateCourseRequest struct {
	MentorID    uuid.UUID `json:"mentor_id" form:"mentor_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" form:"category_id" binding:"required"`
	Title       string    `json:"title" form:"title" binding:"required"`
	Description string    `json:"description" form:"description"`
	Thumbnail   string    `json:"thumbnail" form:"thumbnail"`
	Price       float64   `json:"price" form:"price"`
	Level       string    `json:"level" form:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Duration    int       `json:"duration" form:"duration"`
	Status      string    `json:"status" form:"status" binding:"omitempty,oneof=draft published"`
}

type UpdateCourseRequest struct {
	Title       string  `json:"title" form:"title"`
	Description string  `json:"description" form:"description"`
	Thumbnail   string  `json:"thumbnail" form:"thumbnail"`
	Price       float64 `json:"price" form:"price"`
	Level       string  `json:"level" form:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Duration    int     `json:"duration" form:"duration"`
	Status      string  `json:"status" form:"status" binding:"omitempty,oneof=draft published"`
}

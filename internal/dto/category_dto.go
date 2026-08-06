package dto

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

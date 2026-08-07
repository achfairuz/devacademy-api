package dto

type CreateCategoryRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
	Icon string `json:"icon" form:"icon"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" form:"name"`
	Icon string `json:"icon" form:"icon"`
}

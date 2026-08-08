package level

type CreateLevelRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type UpdateLevelRequest struct {
	Name string `json:"name" form:"name"`
}

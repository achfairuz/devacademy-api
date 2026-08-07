package coursesection

type CourseSectionDto struct {
	Title       string `json:"title" form:"title" binding:"required"`
	OrderNumber int    `json:"order_number" form:"order_number"`
}

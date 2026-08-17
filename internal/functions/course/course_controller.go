package course

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/utils"
	"github.com/iyuz/devacademy-api/pkg/response"
)

const maxThumbSize = 5 << 20

type CourseController struct {
	service CourseService
}

func NewCourseController(service CourseService) *CourseController {
	return &CourseController{service: service}
}

func saveThumbnail(c *gin.Context) (string, error) {
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return "", nil
	}

	path, err := utils.SaveUploadedImage(file, "./uploads/courses", maxThumbSize)
	if err != nil {
		return "", err
	}

	compressed, err := utils.Compress(path)
	if err != nil {
		return "", err
	}

	return "/uploads/courses/" + filepath.Base(compressed), nil
}

// Create godoc
//
//	@Summary		Create a new course
//	@Description	Create course with optional thumbnail upload (multipart/form-data)
//	@Tags			Courses
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			mentor_id	formData	string	true	"Mentor ID"
//	@Param			category_id	formData	string	true	"Category ID"
//	@Param			title		formData	string	true	"Course title"
//	@Param			description	formData	string	false	"Course description"
//	@Param			thumbnail	formData	file	false	"Course thumbnail (jpg/jpeg/png/gif/webp/svg/avif, max 5MB)"
//	@Param			price		formData	number	false	"Course price"
//	@Param			level_id	formData	string	false	"Course level ID"
//	@Param			duration	formData	int		false	"Course duration (hours)"
//	@Param			status		formData	string	false	"Course status (draft|published)"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Router			/courses [post]
func (ctr *CourseController) Create(c *gin.Context) {
	mentorID, err := uuid.Parse(c.PostForm("mentor_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", "invalid mentor_id")
		return
	}

	categoryID, err := uuid.Parse(c.PostForm("category_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", "invalid category_id")
		return
	}

	title := c.PostForm("title")
	if title == "" {
		response.Error(c, http.StatusBadRequest, "invalid request", "title is required")
		return
	}

	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	duration, _ := strconv.Atoi(c.PostForm("duration"))

	levelID, _ := uuid.Parse(c.PostForm("level_id"))

	req := &CreateCourseRequest{
		MentorID:    mentorID,
		CategoryID:  categoryID,
		LevelID:     levelID,
		Title:       title,
		Description: c.PostForm("description"),
		Price:       price,
		Duration:    duration,
		Status:      c.PostForm("status"),
	}

	thumbnail, err := saveThumbnail(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid thumbnail", err.Error())
		return
	}
	req.Thumbnail = thumbnail

	result, err := ctr.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "course created", result)
}

// GetByID godoc
//
//	@Summary		Get course by ID
//	@Description	Retrieve a single course by UUID
//	@Tags			Courses
//	@Produce		json
//	@Param			id	path		string	true	"Course ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/courses/{id} [get]
func (ctr *CourseController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	course, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "course not found", err.Error())
		return
	}

	response.Success(c, "course retrieved", course)
}

// GetDetailBySlug godoc
//
//	@Summary		Get course detail by slug
//	@Description	Retrieve a single course with all its sections and lessons by slug
//	@Tags			Courses
//	@Produce		json
//	@Param			slug	path	string	true	"Course slug"
//	@Success		200		{object}	response.Response{data=CourseDetail}
//	@Failure		404		{object}	response.Response
//	@Router			/courses/slug/{slug}/detail [get]
func (ctr *CourseController) GetDetailBySlug(c *gin.Context) {
	slug := c.Param("slug")
	course, err := ctr.service.GetDetailBySlug(c.Request.Context(), slug)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get course", err.Error())
		return
	}
	if course == nil {
		response.Error(c, http.StatusNotFound, "course not found", "")
		return
	}
	response.Success(c, "course retrieved", course)
}

// GetBySlug godoc
//
//	@Summary		Get course by slug
//	@Description	Retrieve a single course by slug
//	@Tags			Courses
//	@Produce		json
//	@Param			slug	path	string	true	"Course slug"
//	@Success		200		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/courses/slug/{slug} [get]
func (ctr *CourseController) GetBySlug(c *gin.Context) {
	course, err := ctr.service.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		response.Error(c, http.StatusNotFound, "course not found", err.Error())
		return
	}

	response.Success(c, "course retrieved", course)
}

// GetAll godoc
//
//	@Summary		List courses
//	@Description	Retrieve paginated list of courses
//	@Tags			Courses
//	@Produce		json
//	@Param			page		query	int	false	"Page number"
//	@Param			page_size	query	int	false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/courses [get]
func (ctr *CourseController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	courses, err := ctr.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get courses", err.Error())
		return
	}

	response.Success(c, "courses retrieved", courses)
}

// GetCards godoc
//
//	@Summary		List course cards
//	@Description	Retrieve paginated list of published courses as compact cards (optional JWT to include user progress)
//	@Tags			Courses
//	@Produce		json
//	@Param			page		query	int	false	"Page number"
//	@Param			page_size	query	int	false	"Items per page"
//	@Success		200			{object}	response.Response{data=[]CourseCard}
//	@Router			/courses/cards [get]
func (ctr *CourseController) GetCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var userID *uuid.UUID
	if raw, exists := c.Get("user_id"); exists {
		if uid, ok := raw.(uuid.UUID); ok {
			userID = &uid
		}
	}

	cards, err := ctr.service.GetCards(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get course cards", err.Error())
		return
	}

	response.Success(c, "course cards retrieved", cards)
}

// GetByMentor godoc
//
//	@Summary		List courses by mentor
//	@Description	Retrieve paginated list of courses by mentor
//	@Tags			Courses
//	@Produce		json
//	@Param			mentor_id	path	string	true	"Mentor ID"
//	@Param			page		query	int		false	"Page number"
//	@Param			page_size	query	int		false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/courses/mentor/{mentor_id} [get]
func (ctr *CourseController) GetByMentor(c *gin.Context) {
	mentorID, err := uuid.Parse(c.Param("mentor_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid mentor id", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	courses, err := ctr.service.GetByMentor(c.Request.Context(), mentorID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get courses", err.Error())
		return
	}

	response.Success(c, "courses retrieved", courses)
}

// GetByCategory godoc
//
//	@Summary		List courses by category
//	@Description	Retrieve paginated list of courses by category
//	@Tags			Courses
//	@Produce		json
//	@Param			category_id	path	string	true	"Category ID"
//	@Param			page		query	int		false	"Page number"
//	@Param			page_size	query	int		false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/courses/category/{category_id} [get]
func (ctr *CourseController) GetByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid category id", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	courses, err := ctr.service.GetByCategory(c.Request.Context(), categoryID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get courses", err.Error())
		return
	}

	response.Success(c, "courses retrieved", courses)
}

// GetByLevel godoc
//
//	@Summary		List courses by level
//	@Description	Retrieve paginated list of courses by level
//	@Tags			Courses
//	@Produce		json
//	@Param			level_id	path	string	true	"Level ID"
//	@Param			page		query	int		false	"Page number"
//	@Param			page_size	query	int		false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/courses/level/{level_id} [get]
func (ctr *CourseController) GetByLevel(c *gin.Context) {
	levelID, err := uuid.Parse(c.Param("level_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid level id", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	courses, err := ctr.service.GetByLevel(c.Request.Context(), levelID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get courses", err.Error())
		return
	}

	response.Success(c, "courses retrieved", courses)
}

// Update godoc
//
//	@Summary		Update course
//	@Description	Update course fields (multipart/form-data)
//	@Tags			Courses
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id			path		string	true	"Course ID"
//	@Param			title		formData	string	false	"Course title"
//	@Param			description	formData	string	false	"Course description"
//	@Param			thumbnail	formData	file	false	"Course thumbnail"
//	@Param			price		formData	number	false	"Course price"
//	@Param			level		formData	string	false	"Course level (beginner|intermediate|advanced)"
//	@Param			duration	formData	int		false	"Course duration (hours)"
//	@Param			status		formData	string	false	"Course status (draft|published)"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id} [put]
func (ctr *CourseController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	req := &UpdateCourseRequest{}
	if err := c.ShouldBind(req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	thumbnail, err := saveThumbnail(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid thumbnail", err.Error())
		return
	}
	if thumbnail != "" {
		req.Thumbnail = thumbnail
	}

	course, err := ctr.service.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "course updated", course)
}

// UpdateStatus godoc
//
//	@Summary		Update course status
//	@Description	Update course status by slug (draft|published)
//	@Tags			Courses
//	@Produce		json
//	@Param			slug	path	string	true	"Course slug"
//	@Param			status	query	string	true	"Course status (draft|published)"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/courses/slug/{slug}/status [patch]
func (ctr *CourseController) UpdateStatus(c *gin.Context) {
	slug := c.Param("slug")
	status := c.Query("status")
	if status == "" {
		response.Error(c, http.StatusBadRequest, "invalid request", "status is required")
		return
	}
	if status != "draft" && status != "published" {
		response.Error(c, http.StatusBadRequest, "invalid request", "status must be draft or published")
		return
	}

	if err := ctr.service.UpdateStatus(c.Request.Context(), slug, status); err != nil {
		response.Error(c, http.StatusNotFound, "update status failed", err.Error())
		return
	}

	response.Success(c, "course status updated", nil)
}

// Delete godoc
//
//	@Summary		Delete course
//	@Description	Soft delete a course by ID
//	@Tags			Courses
//	@Produce		json
//	@Param			id	path	string	true	"Course ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/courses/{id} [delete]
func (ctr *CourseController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "course deleted", nil)
}

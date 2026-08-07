package coursesection

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/functions/course"
	"github.com/iyuz/devacademy-api/pkg/response"
)

type CourseSectionController struct {
	service       CourseSectionService
	CourseService course.CourseService
}

func NewCourseSectionController(service CourseSectionService, courseService course.CourseService) *CourseSectionController {
	return &CourseSectionController{service: service, CourseService: courseService}
}

// Create godoc
//
//	@Summary		Create course section
//	@Description	Add a new section to a course, validating the course exists
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			course_id	path	string				true	"Course ID"
//	@Param			request		body	CourseSectionDto	true	"Section payload"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{course_id}/sections [post]
func (ctr *CourseSectionController) Create(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("course_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid course id", err.Error())
		return
	}

	if _, err := ctr.CourseService.GetByID(c.Request.Context(), courseID); err != nil {
		response.Error(c, http.StatusNotFound, "course not found", err.Error())
		return
	}

	var req CourseSectionDto
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	section, err := ctr.service.Create(c.Request.Context(), courseID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "section created", section)
}

// GetByCourse godoc
//
//	@Summary		List course sections
//	@Description	Retrieve paginated list of sections for a course
//	@Tags			Courses
//	@Produce		json
//	@Param			course_id	path	int	true	"Course ID"
//	@Param			page		query	int	false	"Page number"
//	@Param			page_size	query	int	false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/courses/{course_id}/sections [get]
func (ctr *CourseSectionController) GetByCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("course_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid course id", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	sections, err := ctr.service.GetByCourse(c.Request.Context(), courseID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get sections", err.Error())
		return
	}

	response.Success(c, "sections retrieved", sections)
}

// GetByID godoc
//
//	@Summary		Get course section by ID
//	@Description	Retrieve a single course section by UUID
//	@Tags			Courses
//	@Produce		json
//	@Param			id	path	string	true	"Section ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/sections/{id} [get]
func (ctr *CourseSectionController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	section, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "section not found", err.Error())
		return
	}

	response.Success(c, "section retrieved", section)
}

// Update godoc
//
//	@Summary		Update course section
//	@Description	Update section title and/or order number
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	string				true	"Section ID"
//	@Param			request	body	CourseSectionDto	true	"Section payload"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/sections/{id} [put]
func (ctr *CourseSectionController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var req CourseSectionDto
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	section, err := ctr.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "section updated", section)
}

// Delete godoc
//
//	@Summary		Delete course section
//	@Description	Delete a course section by ID
//	@Tags			Courses
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Section ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/sections/{id} [delete]
func (ctr *CourseSectionController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "section deleted", nil)
}

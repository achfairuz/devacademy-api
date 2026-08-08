package lessons

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	coursesection "github.com/iyuz/devacademy-api/internal/functions/course/course_section"
	"github.com/iyuz/devacademy-api/pkg/response"
)

type LessonController struct {
	service        LessonService
	sectionService coursesection.CourseSectionService
}

func NewLessonController(service LessonService, sectionService coursesection.CourseSectionService) *LessonController {
	return &LessonController{service: service, sectionService: sectionService}
}

// Store godoc
//
//	@Summary		Create lesson
//	@Description	Add a new lesson to a section, validating the section exists
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string			true	"Course ID"
//	@Param			section_id	path	string			true	"Section ID"
//	@Param			request		body	LessonRequest	true	"Lesson payload"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons [post]
func (ctr *LessonController) Store(c *gin.Context) {
	sectionUUID, err := uuid.Parse(c.Param("section_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid section_id", err.Error())
		return
	}

	if _, err := ctr.sectionService.GetByID(c.Request.Context(), sectionUUID); err != nil {
		response.Error(c, http.StatusNotFound, "section not found", err.Error())
		return
	}

	var req LessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	lesson, err := ctr.service.Create(c.Request.Context(), &req, sectionUUID.String())
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "lesson created", lesson)
}

// Index godoc
//
//	@Summary		List lessons
//	@Description	Retrieve list of lessons for a section
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Success		200			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons [get]
func (ctr *LessonController) Index(c *gin.Context) {
	sectionUUID, err := uuid.Parse(c.Param("section_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid section_id", err.Error())
		return
	}

	if _, err := ctr.sectionService.GetByID(c.Request.Context(), sectionUUID); err != nil {
		response.Error(c, http.StatusNotFound, "section not found", err.Error())
		return
	}

	lessons, err := ctr.service.GetBySectionID(c.Request.Context(), sectionUUID.String())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get lessons", err.Error())
		return
	}

	response.Success(c, "lessons retrieved", lessons)
}

// Show godoc
//
//	@Summary		Get lesson by ID
//	@Description	Retrieve a single lesson by UUID
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id} [get]
func (ctr *LessonController) Show(c *gin.Context) {
	id, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson id", err.Error())
		return
	}

	lesson, err := ctr.service.GetByID(c.Request.Context(), id.String())
	if err != nil {
		response.Error(c, http.StatusNotFound, "lesson not found", err.Error())
		return
	}

	response.Success(c, "lesson retrieved", lesson)
}

// Update godoc
//
//	@Summary		Update lesson
//	@Description	Update lesson title, description and/or video url
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string			true	"Course ID"
//	@Param			section_id	path	string			true	"Section ID"
//	@Param			lesson_id	path	string			true	"Lesson ID"
//	@Param			request		body	LessonRequest	true	"Lesson payload"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id} [put]
func (ctr *LessonController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson id", err.Error())
		return
	}

	var req LessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	lesson, err := ctr.service.Update(c.Request.Context(), id.String(), &req, c.Param("section_id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "lesson updated", lesson)
}

// Delete godoc
//
//	@Summary		Delete lesson
//	@Description	Delete a lesson by ID
//	@Tags			Courses
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id} [delete]
func (ctr *LessonController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id.String()); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "lesson deleted", nil)
}

package assignment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/functions/course/lessons"
	"github.com/iyuz/devacademy-api/pkg/response"
)

type AssignmentController struct {
	service       AssignmentService
	lessonService lessons.LessonService
}

func NewAssignmentController(service AssignmentService, lessonService lessons.LessonService) *AssignmentController {
	return &AssignmentController{service: service, lessonService: lessonService}
}

// Store godoc
//
//	@Summary		Create assignment
//	@Description	Create an assignment for a lesson, validating the lesson exists
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string				true	"Course ID"
//	@Param			section_id	path	string				true	"Section ID"
//	@Param			lesson_id	path	string				true	"Lesson ID"
//	@Param			request		body	AssignmentRequest	true	"Assignment payload"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/assignments [post]
func (ctr *AssignmentController) Store(c *gin.Context) {
	lessonUUID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson_id", err.Error())
		return
	}

	if _, err := ctr.lessonService.GetByID(c.Request.Context(), lessonUUID.String()); err != nil {
		response.Error(c, http.StatusNotFound, "lesson not found", err.Error())
		return
	}

	var req AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	assignment, err := ctr.service.Create(c.Request.Context(), lessonUUID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "assignment created", assignment)
}

// Index godoc
//
//	@Summary		List assignments
//	@Description	Retrieve list of assignments for a lesson
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Success		200			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/assignments [get]
func (ctr *AssignmentController) Index(c *gin.Context) {
	lessonUUID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson_id", err.Error())
		return
	}

	assignments, err := ctr.service.GetByLessonID(c.Request.Context(), lessonUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get assignments", err.Error())
		return
	}

	response.Success(c, "assignments retrieved", assignments)
}

// Show godoc
//
//	@Summary		Get assignment by ID
//	@Description	Retrieve a single assignment by UUID
//	@Tags			Courses
//	@Produce		json
//	@Param			id				path	string	true	"Course ID"
//	@Param			section_id		path	string	true	"Section ID"
//	@Param			lesson_id		path	string	true	"Lesson ID"
//	@Param			assignment_id	path	string	true	"Assignment ID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/assignments/{assignment_id} [get]
func (ctr *AssignmentController) Show(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assignment_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid assignment id", err.Error())
		return
	}

	assignment, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "assignment not found", err.Error())
		return
	}

	response.Success(c, "assignment retrieved", assignment)
}

// Update godoc
//
//	@Summary		Update assignment
//	@Description	Update assignment title, description and/or due date
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path	string				true	"Course ID"
//	@Param			section_id		path	string				true	"Section ID"
//	@Param			lesson_id		path	string				true	"Lesson ID"
//	@Param			assignment_id	path	string				true	"Assignment ID"
//	@Param			request			body	AssignmentRequest	true	"Assignment payload"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/assignments/{assignment_id} [put]
func (ctr *AssignmentController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assignment_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid assignment id", err.Error())
		return
	}

	var req AssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	assignment, err := ctr.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "assignment updated", assignment)
}

// Delete godoc
//
//	@Summary		Delete assignment
//	@Description	Delete an assignment by ID
//	@Tags			Courses
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path	string	true	"Course ID"
//	@Param			section_id		path	string	true	"Section ID"
//	@Param			lesson_id		path	string	true	"Lesson ID"
//	@Param			assignment_id	path	string	true	"Assignment ID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/assignments/{assignment_id} [delete]
func (ctr *AssignmentController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assignment_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid assignment id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "assignment deleted", nil)
}

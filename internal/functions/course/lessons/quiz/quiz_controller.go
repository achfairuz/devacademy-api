package quiz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/functions/course/lessons"
	"github.com/iyuz/devacademy-api/pkg/response"
)

type QuizController struct {
	service       QuizService
	lessonService lessons.LessonService
}

func NewQuizController(service QuizService, lessonService lessons.LessonService) *QuizController {
	return &QuizController{service: service, lessonService: lessonService}
}

// Store godoc
//
//	@Summary		Create quiz
//	@Description	Create a quiz with questions and options for a lesson, validating the lesson exists
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string			true	"Course ID"
//	@Param			section_id	path	string			true	"Section ID"
//	@Param			lesson_id	path	string			true	"Lesson ID"
//	@Param			request		body	QuizRequest		true	"Quiz payload"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/quizzes [post]
func (ctr *QuizController) Store(c *gin.Context) {
	lessonUUID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson_id", err.Error())
		return
	}

	if _, err := ctr.lessonService.GetByID(c.Request.Context(), lessonUUID.String()); err != nil {
		response.Error(c, http.StatusNotFound, "lesson not found", err.Error())
		return
	}

	var req QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	quiz, err := ctr.service.Create(c.Request.Context(), lessonUUID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "quiz created", quiz)
}

// Index godoc
//
//	@Summary		List quizzes
//	@Description	Retrieve list of quizzes for a lesson
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Success		200			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/quizzes [get]
func (ctr *QuizController) Index(c *gin.Context) {
	lessonUUID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson_id", err.Error())
		return
	}

	quizzes, err := ctr.service.GetByLessonID(c.Request.Context(), lessonUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get quizzes", err.Error())
		return
	}

	response.Success(c, "quizzes retrieved", quizzes)
}

// Show godoc
//
//	@Summary		Get quiz by ID
//	@Description	Retrieve a single quiz with questions and options by UUID
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Param			quiz_id		path	string	true	"Quiz ID"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/quizzes/{quiz_id} [get]
func (ctr *QuizController) Show(c *gin.Context) {
	id, err := uuid.Parse(c.Param("quiz_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid quiz id", err.Error())
		return
	}

	quiz, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "quiz not found", err.Error())
		return
	}

	response.Success(c, "quiz retrieved", quiz)
}

// Update godoc
//
//	@Summary		Update quiz
//	@Description	Update quiz title, passing score, and replace its questions and options
//	@Tags			Courses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string			true	"Course ID"
//	@Param			section_id	path	string			true	"Section ID"
//	@Param			lesson_id	path	string			true	"Lesson ID"
//	@Param			quiz_id		path	string			true	"Quiz ID"
//	@Param			request		body	QuizRequest		true	"Quiz payload"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/quizzes/{quiz_id} [put]
func (ctr *QuizController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("quiz_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid quiz id", err.Error())
		return
	}

	var req QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	quiz, err := ctr.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "quiz updated", quiz)
}

// Delete godoc
//
//	@Summary		Delete quiz
//	@Description	Delete a quiz and its questions, options and attempts by ID
//	@Tags			Courses
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Param			quiz_id		path	string	true	"Quiz ID"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/quizzes/{quiz_id} [delete]
func (ctr *QuizController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("quiz_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid quiz id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "quiz deleted", nil)
}

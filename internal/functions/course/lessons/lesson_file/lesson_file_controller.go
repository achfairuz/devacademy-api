package lessonfile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/functions/course/lessons"
	"github.com/iyuz/devacademy-api/internal/utils"
	"github.com/iyuz/devacademy-api/pkg/response"
)

const maxFileSize = 50 << 20

type LessonFileController struct {
	service       LessonFileService
	lessonService lessons.LessonService
}

func NewLessonFileController(service LessonFileService, lessonService lessons.LessonService) *LessonFileController {
	return &LessonFileController{service: service, lessonService: lessonService}
}

// Store godoc
//
//	@Summary		Upload lesson file
//	@Description	Upload a file (e.g. PDF, DOCX) to a lesson, validating the lesson exists
//	@Tags			Courses
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Param			file		formData	file	true	"Lesson file (max 50MB)"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/files [post]
func (ctr *LessonFileController) Store(c *gin.Context) {
	lessonUUID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson_id", err.Error())
		return
	}

	if _, err := ctr.lessonService.GetByID(c.Request.Context(), lessonUUID.String()); err != nil {
		response.Error(c, http.StatusNotFound, "lesson not found", err.Error())
		return
	}

	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required", err.Error())
		return
	}

	dir := "./uploads/lessons/" + lessonUUID.String()
	filename, err := utils.SaveUploadedFile(header, dir, maxFileSize)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file", err.Error())
		return
	}

	file, err := ctr.service.Create(c.Request.Context(), lessonUUID, header.Filename, "/uploads/lessons/"+lessonUUID.String()+"/"+filename, int(header.Size))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "lesson file uploaded", file)
}

// Index godoc
//
//	@Summary		List lesson files
//	@Description	Retrieve list of files for a lesson
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Success		200			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/files [get]
func (ctr *LessonFileController) Index(c *gin.Context) {
	lessonUUID, err := uuid.Parse(c.Param("lesson_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid lesson_id", err.Error())
		return
	}

	files, err := ctr.service.GetByLessonID(c.Request.Context(), lessonUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get lesson files", err.Error())
		return
	}

	response.Success(c, "lesson files retrieved", files)
}

// Show godoc
//
//	@Summary		Get lesson file by ID
//	@Description	Retrieve a single lesson file by UUID
//	@Tags			Courses
//	@Produce		json
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Param			file_id		path	string	true	"File ID"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/files/{file_id} [get]
func (ctr *LessonFileController) Show(c *gin.Context) {
	id, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file id", err.Error())
		return
	}

	file, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "lesson file not found", err.Error())
		return
	}

	response.Success(c, "lesson file retrieved", file)
}

// Delete godoc
//
//	@Summary		Delete lesson file
//	@Description	Delete a lesson file by ID
//	@Tags			Courses
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path	string	true	"Course ID"
//	@Param			section_id	path	string	true	"Section ID"
//	@Param			lesson_id	path	string	true	"Lesson ID"
//	@Param			file_id		path	string	true	"File ID"
//	@Success		200			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Failure		404			{object}	response.Response
//	@Router			/courses/{id}/sections/{section_id}/lessons/{lesson_id}/files/{file_id} [delete]
func (ctr *LessonFileController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid file id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "lesson file deleted", nil)
}

package level

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/pkg/response"
)

type LevelController struct {
	service LevelService
}

func NewLevelController(service LevelService) *LevelController {
	return &LevelController{service: service}
}

// Create godoc
//
//	@Summary		Create a new level
//	@Description	Create a course level
//	@Tags			Levels
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	CreateLevelRequest	true	"Level payload"
//	@Success		201		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Router			/levels [post]
func (ctr *LevelController) Create(c *gin.Context) {
	req := &CreateLevelRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := ctr.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "level created", result)
}

// GetByID godoc
//
//	@Summary		Get level by ID
//	@Description	Retrieve a single level by UUID
//	@Tags			Levels
//	@Produce		json
//	@Param			id	path		string	true	"Level ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/levels/{id} [get]
func (ctr *LevelController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	level, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "level not found", err.Error())
		return
	}

	response.Success(c, "level retrieved", level)
}

// GetAll godoc
//
//	@Summary		List levels
//	@Description	Retrieve paginated list of levels
//	@Tags			Levels
//	@Produce		json
//	@Param			page		query	int	false	"Page number"
//	@Param			page_size	query	int	false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/levels [get]
func (ctr *LevelController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	levels, err := ctr.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get levels", err.Error())
		return
	}

	response.Success(c, "levels retrieved", levels)
}

// Update godoc
//
//	@Summary		Update level
//	@Description	Update a level name
//	@Tags			Levels
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path	string				true	"Level ID"
//	@Param			request	body	UpdateLevelRequest	true	"Level payload"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Router			/levels/{id} [put]
func (ctr *LevelController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var req UpdateLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	level, err := ctr.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "level updated", level)
}

// Delete godoc
//
//	@Summary		Delete level
//	@Description	Delete a level by ID
//	@Tags			Levels
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	string	true	"Level ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/levels/{id} [delete]
func (ctr *LevelController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "level deleted", nil)
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/dto"
	"github.com/iyuz/devacademy-api/internal/services"
	"github.com/iyuz/devacademy-api/pkg/response"
)

type CategoryController struct {
	service services.CategoryService
}

func NewCategoryController(service services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// Create godoc
//
//	@Summary		Create a new category
//	@Description	Create category with icon name (multipart/form-data)
//	@Tags			Categories
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name	formData	string	true	"Category name"
//	@Param			icon	formData	string	false	"Category icon name (Lucide)"
//	@Success		201		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Router			/mentor/categories [post]
func (ctr *CategoryController) Create(c *gin.Context) {
	req := &dto.CreateCategoryRequest{}
	if err := c.ShouldBind(req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	if req.Name == "" {
		response.Error(c, http.StatusBadRequest, "invalid request", "name is required")
		return
	}

	result, err := ctr.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create failed", err.Error())
		return
	}

	response.Created(c, "category created", result)
}

// GetByID godoc
//
//	@Summary		Get category by ID
//	@Description	Retrieve a single category by UUID
//	@Tags			Categories
//	@Produce		json
//	@Param			id	path		string	true	"Category ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/mentor/categories/{id} [get]
func (ctr *CategoryController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	category, err := ctr.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "category not found", err.Error())
		return
	}

	response.Success(c, "category retrieved", category)
}

// GetAll godoc
//
//	@Summary		List categories
//	@Description	Retrieve paginated list of categories
//	@Tags			Categories
//	@Produce		json
//	@Param			page		query	int	false	"Page number"
//	@Param			page_size	query	int	false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Router			/mentor/categories [get]
func (ctr *CategoryController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	categories, err := ctr.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get categories", err.Error())
		return
	}

	response.Success(c, "categories retrieved", categories)
}

// Delete godoc
//
//	@Summary		Delete category
//	@Description	Delete a category by ID
//	@Tags			Categories
//	@Produce		json
//	@Param			id	path	string	true	"Category ID"
//	@Success		200	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/mentor/categories/{id} [delete]
func (ctr *CategoryController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "category deleted", nil)
}

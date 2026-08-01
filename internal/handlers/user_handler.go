package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/services"
	"github.com/iyuz/devacademy-api/pkg/response"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type updateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account and return JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		registerRequest	true	"Register payload"
//	@Success		201		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Router			/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := h.service.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusConflict, "registration failed", err.Error())
		return
	}

	response.Created(c, "user registered successfully", result)
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login payload"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "login failed", err.Error())
		return
	}

	response.Success(c, "login successful", result)
}

// GetByID godoc
//
//	@Summary		Get user by ID
//	@Description	Retrieve a single user by UUID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "user not found", err.Error())
		return
	}

	response.Success(c, "user retrieved", user)
}

// GetAll godoc
//
//	@Summary		List all users
//	@Description	Retrieve paginated list of users
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page		query		int	false	"Page number"
//	@Param			page_size	query		int	false	"Items per page"
//	@Success		200			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Router			/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get users", err.Error())
		return
	}

	response.Success(c, "users retrieved", users)
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Update user name and/or email
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string			true	"User ID"
//	@Param			request	body		updateRequest	true	"Update payload"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, req.Name, req.Email)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "update failed", err.Error())
		return
	}

	response.Success(c, "user updated", user)
}

// Delete godoc
//
//	@Summary		Delete user
//	@Description	Soft delete a user by ID
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "user deleted", nil)
}

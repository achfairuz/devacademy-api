package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/pkg/response"
)

type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{service: service}
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
func (ctr *UserController) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := ctr.service.Register(c.Request.Context(), req.FullName, req.Username, req.Email, req.Password)
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
func (ctr *UserController) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := ctr.service.Login(c.Request.Context(), req.Email, req.Password)
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
func (ctr *UserController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	user, err := ctr.service.GetByID(c.Request.Context(), id)
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
func (ctr *UserController) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := ctr.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get users", err.Error())
		return
	}

	response.Success(c, "users retrieved", users)
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Update user full name and/or email
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
func (ctr *UserController) Update(c *gin.Context) {
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

	user, err := ctr.service.Update(c.Request.Context(), id, req.FullName, req.Email)
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
func (ctr *UserController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	if err := ctr.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusNotFound, "delete failed", err.Error())
		return
	}

	response.Success(c, "user deleted", nil)
}

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/apperror"
	"simple-auth/internal/httpx"
	"simple-auth/internal/middleware"
)

// Handler is the HTTP layer for users (≈ UserController in NestJS).
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /users
func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	u, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, ToResponse(u))
}

// GET /users
func (h *Handler) List(c *gin.Context) {
	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, ToResponseList(users))
}

// GET /users/:id
func (h *Handler) Get(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}

	u, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, ToResponse(u))
}

// PATCH /users/:id
func (h *Handler) Update(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}

	actor, ok := currentActor(c)
	if !ok {
		return
	}

	var req UpdateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	u, err := h.svc.Update(c.Request.Context(), actor, id, req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, ToResponse(u))
}

// PATCH /users/:id/role
func (h *Handler) ChangeRole(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}

	actor, ok := currentActor(c)
	if !ok {
		return
	}

	var req UpdateRoleRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	u, err := h.svc.ChangeRole(c.Request.Context(), actor, id, req.Role)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, ToResponse(u))
}

// DELETE /users/:id
func (h *Handler) Delete(c *gin.Context) {
	id, ok := httpx.ParamID(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		httpx.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// currentActor builds an Actor from the token claims set by RequireAuth.
// On failure it writes a 401 and returns false.
func currentActor(c *gin.Context) (Actor, bool) {
	claims, ok := middleware.CurrentUser(c)
	if !ok {
		httpx.Error(c, apperror.Unauthorized("not authenticated"))
		return Actor{}, false
	}
	return Actor{ID: claims.UserID, Role: Role(claims.Role)}, true
}

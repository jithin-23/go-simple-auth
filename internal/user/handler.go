package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/httpx"
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

	var req UpdateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	u, err := h.svc.Update(c.Request.Context(), id, req)
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

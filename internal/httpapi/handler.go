package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Krchnk/go-micro/internal/users"
)

type Handler struct {
	service *users.Service
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *users.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.POST("/users", h.createUser)
	router.GET("/users", h.listUsers)
	router.PUT("/users/:id", h.updateUser)
	router.DELETE("/users/:id", h.deleteUser)

	return router
}

func (h *Handler) createUser(c *gin.Context) {
	var req createUserRequest
	if err := bindJSONStrict(c, &req); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		writeError(c, http.StatusBadRequest, "name and email are required")
		return
	}

	created := h.service.Create(strings.TrimSpace(req.Name), strings.TrimSpace(req.Email))
	c.JSON(http.StatusCreated, created)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, err := parseUserID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req updateUserRequest
	if err := bindJSONStrict(c, &req); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		writeError(c, http.StatusBadRequest, "name and email are required")
		return
	}

	updated, err := h.service.Update(id, strings.TrimSpace(req.Name), strings.TrimSpace(req.Email))
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			writeError(c, http.StatusNotFound, "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, err := parseUserID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			writeError(c, http.StatusNotFound, "user not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.List())
}

func bindJSONStrict(c *gin.Context, out any) error {
	defer c.Request.Body.Close()

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return errors.New("invalid JSON body")
	}

	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func parseUserID(idParam string) (int64, error) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user id in path")
	}

	return id, nil
}

func writeError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, errorResponse{Error: message})
}

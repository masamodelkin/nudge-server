package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masamodelkin/nudge-server/internal/service"
)

type StatusHandler struct {
	service *service.StatusService
}

func NewStatusHandler(s *service.StatusService) *StatusHandler {
	return &StatusHandler{service: s}
}

type statusRequest struct {
	Name         string  `json:"name" binding:"required"`
	NextStatusID *string `json:"next_status_id"`
	IsDone       bool    `json:"is_done"`
}

func (h *StatusHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.POST("statuses", h.Create)
	protected.GET("statuses", h.List)
	protected.PUT("statuses/:id", h.Update)
	protected.DELETE("statuses/:id", h.Delete)
}

func (h *StatusHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.service.Create(userID, &service.StatusRequest{
		Name:         req.Name,
		NextStatusID: req.NextStatusID,
		IsDone:       req.IsDone,
	})
	if err != nil {
		if errors.Is(err, service.ErrStatusValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusCreated, status)
}

func (h *StatusHandler) List(c *gin.Context) {
	userID := c.GetString("userID")

	statuses, err := h.service.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

func (h *StatusHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.service.Update(id, userID, &service.StatusRequest{
		Name:         req.Name,
		NextStatusID: req.NextStatusID,
		IsDone:       req.IsDone,
	})

	if err != nil {
		if errors.Is(err, service.ErrStatusNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrStatusValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *StatusHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	err := h.service.Delete(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrStatusNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status deleted"})
}

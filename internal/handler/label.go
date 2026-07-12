package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masamodelkin/nudge-server/internal/service"
)

type LabelHandler struct {
	service *service.LabelService
}

func NewLabelHandler(s *service.LabelService) *LabelHandler {
	return &LabelHandler{service: s}
}

type labelRequest struct {
	Name  string  `json:"name" binding:"required"`
	Color *string `json:"color"`
}

func (h *LabelHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.POST("labels", h.Create)
	protected.GET("labels", h.List)
	protected.PUT("labels/:id", h.Update)
	protected.DELETE("labels/:id", h.Delete)
}

func (h *LabelHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req labelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label, err := h.service.Create(userID, &service.LabelRequest{
		Name:  req.Name,
		Color: req.Color,
	})
	if err != nil {
		if errors.Is(err, service.ErrLabelValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusCreated, label)
}

func (h *LabelHandler) List(c *gin.Context) {
	userID := c.GetString("userID")

	labels, err := h.service.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, labels)
}

func (h *LabelHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	var req labelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label, err := h.service.Update(id, userID, &service.LabelRequest{
		Name:  req.Name,
		Color: req.Color,
	})
	if err != nil {
		if errors.Is(err, service.ErrLabelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrLabelValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *LabelHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	err := h.service.Delete(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrLabelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "label deleted"})
}

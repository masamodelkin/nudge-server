package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"github.com/masamodelkin/nudge-server/internal/service"
)

type TriggerHandler struct {
	service *service.TriggerService
}

func NewTriggerHandler(s *service.TriggerService) *TriggerHandler {
	return &TriggerHandler{service: s}
}

type triggerRequest struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Config      types.JSONText `json:"config"`
	IsExclusive bool           `json:"is_exclusive"`
}

func (h *TriggerHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.POST("triggers", h.Create)
	protected.GET("triggers", h.List)
	protected.GET("triggers/:id", h.Get)
	protected.PUT("triggers/:id", h.Update)
	protected.DELETE("triggers/:id", h.Delete)
}

func (h *TriggerHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req triggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trigger, err := h.service.Create(userID, &service.TriggerRequest{
		Name:        req.Name,
		Type:        req.Type,
		Config:      req.Config,
		IsExclusive: req.IsExclusive,
	})

	if err != nil {
		if errors.Is(err, service.ErrTriggerValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}
	c.JSON(http.StatusCreated, trigger)
}

func (h *TriggerHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	task, err := h.service.Get(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrTriggerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TriggerHandler) List(c *gin.Context) {
	userID := c.GetString("userID")

	tasks, err := h.service.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TriggerHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	var req triggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trigger, err := h.service.Update(id, userID, &service.TriggerRequest{
		Name:        req.Name,
		Type:        req.Type,
		Config:      req.Config,
		IsExclusive: req.IsExclusive,
	})

	if err != nil {
		if errors.Is(err, service.ErrTriggerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrTriggerValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, trigger)
}

func (h *TriggerHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	if err := h.service.Delete(id, userID); err != nil {
		if errors.Is(err, service.ErrTriggerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trigger deleted"})
}

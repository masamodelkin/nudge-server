package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masamodelkin/nudge-server/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

type taskRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	IsDraft     bool     `json:"is_draft"`
	DueDate     *int64   `json:"due_date"`
	Priority    *int     `json:"priority"`
	Duration    *int     `json:"duration"`
	StatusID    *string  `json:"status_id"`
	LabelIDs    []string `json:"label_ids"`
}

type timeRequest struct {
	Seconds int `json:"minutes" binding:"required"`
}

func (h *TaskHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.POST("tasks", h.Create)
	protected.GET("tasks", h.List)
	protected.GET("tasks/:id", h.Get)
	protected.PUT("tasks/:id", h.Update)
	protected.DELETE("tasks/:id", h.Delete)
	protected.POST("tasks/:id/time/add", h.AddTime)
	protected.POST("tasks/:id/time/set", h.SetTime)
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Create(userID, &service.TaskRequest{
		Name:        req.Name,
		Description: req.Description,
		IsDraft:     req.IsDraft,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Duration:    req.Duration,
		StatusID:    req.StatusID,
		LabelIDs:    req.LabelIDs,
	})
	if err != nil {
		if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Get(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	task, err := h.service.Get(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	userID := c.GetString("userID")

	tasks, err := h.service.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Update(id, userID, &service.TaskRequest{
		Name:        req.Name,
		Description: req.Description,
		IsDraft:     req.IsDraft,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Duration:    req.Duration,
		StatusID:    req.StatusID,
		LabelIDs:    req.LabelIDs,
	})
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	if err := h.service.Delete(id, userID); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

func (h *TaskHandler) AddTime(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	var req timeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.AddTime(id, userID, req.Seconds)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) SetTime(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")

	var req timeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.SetTime(id, userID, req.Seconds)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

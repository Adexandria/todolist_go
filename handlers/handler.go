package handlers

import (
	"SM/models"
	"SM/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service services.ITaskService
}

func TaskHandler(service *services.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetTaskById(c *gin.Context) {
	id := c.Param("id")
	val, _ := strconv.Atoi(id)
	currentTask := h.Service.GetTaskById(val)
	if currentTask == (models.TaskDTO{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": currentTask})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var createTaskDTO models.CreateTaskDTO
	if err := c.ShouldBind(&createTaskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newId, error := h.Service.CreateTask(&createTaskDTO)
	if error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": newId})

}

func (h *Handler) UpdateTask(c *gin.Context) {
	var updateTaskDTO models.UpdateTaskDTO
	if err := c.ShouldBind(&updateTaskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id := c.Param("id")
	val, _ := strconv.Atoi(id)
	currentTask := h.Service.GetTaskById(val)
	if currentTask == (models.TaskDTO{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	isUpdated := h.Service.UpdateTask(val, &updateTaskDTO)
	if isUpdated {
		c.JSON(http.StatusOK, gin.H{"data": currentTask})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	return
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	val, _ := strconv.Atoi(id)
	currentTask := h.Service.GetTaskById(val)
	if currentTask == (models.TaskDTO{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	isDeleted := h.Service.DeleteTask(val)
	if isDeleted {
		c.JSON(http.StatusOK, gin.H{"data": currentTask})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})

	return
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	page := c.DefaultQuery("page", "1")

	pageSize := c.DefaultQuery("pageSize", "20")

	pageInt, _ := strconv.Atoi(page)

	pageSizeInt, _ := strconv.Atoi(pageSize)

	allTasks := h.Service.GetAllTasks(pageInt, pageSizeInt)

	c.JSON(http.StatusOK, gin.H{"data": allTasks})

	return
}

func (h *Handler) GetTaskByMonthAndYear(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	month, _ := strconv.Atoi(c.Query("month")) // ← c.Query not c.Param
	year, _ := strconv.Atoi(c.Query("year"))

	allTasks := h.Service.GetTaskByMonthAndYear(month, year, page, pageSize)
	c.JSON(http.StatusOK, gin.H{"data": allTasks})

	return

}

func (h *Handler) GetTaskByCreatedDate(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	day, _ := strconv.Atoi(c.Query("day")) // ← c.Query not c.Param
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	allTasks := h.Service.GetTaskByCreatedDate(day, month, year, page, pageSize)
	c.JSON(http.StatusOK, gin.H{"data": allTasks})
}

func (h *Handler) GetTaskByYear(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	year, _ := strconv.Atoi(c.Query("year")) // ← c.Query not c.Param

	allTasks := h.Service.GetTaskByYear(year, page, pageSize)
	c.JSON(http.StatusOK, gin.H{"data": allTasks})
	return
}

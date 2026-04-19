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
		c.JSON(http.StatusOK, gin.H{"message": "task updated"})
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
		c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})

	return
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	allTasks := h.Service.GetAllTasks(page, pageSize)

	c.JSON(http.StatusOK, gin.H{"data": allTasks})

	return
}

func (h *Handler) FilterTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	year := c.Query("year")
	month := c.Query("month")
	day := c.Query("day")

	if year != "" && month != "" && day != "" {
		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		d, _ := strconv.Atoi(day)
		tasks := h.Service.GetTaskByCreatedDate(d, m, y, page, pageSize)
		c.JSON(http.StatusOK, gin.H{"data": tasks})
		return
	}

	if year != "" && month != "" {
		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		tasks := h.Service.GetTaskByMonthAndYear(m, y, page, pageSize)
		c.JSON(http.StatusOK, gin.H{"data": tasks})
		return
	}

	if year != "" {
		y, _ := strconv.Atoi(year)
		tasks := h.Service.GetTaskByYear(y, page, pageSize)
		c.JSON(http.StatusOK, gin.H{"data": tasks})
		return
	}

	tasks := h.Service.GetAllTasks(page, pageSize)
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (h *Handler) SearchByTask(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	name := c.Query("name")

	tasks := h.Service.SearchTaskByName(name, page, pageSize)

	c.JSON(http.StatusOK, gin.H{"data": tasks})

	return
}

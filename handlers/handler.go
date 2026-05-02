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

// Constructor function to instantiate with the TaskService dependency injected
func TaskHandler(service *services.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

// Handler method to retrieve a task by its ID
func (h *Handler) GetTaskById(c *gin.Context) {
	id := c.Param("id")
	val, _ := strconv.Atoi(id)
	currentTask := h.Service.GetTaskById(val)
	if currentTask == (models.TaskDTO{}) {
		c.JSON(http.StatusNotFound, services.ActionResult{
			StatusCode: http.StatusNotFound,
			IsSuccess:  false,
			Error: []string{
				"Task not found",
			},
		})
	}
	c.JSON(http.StatusOK, services.ActionResultModel[models.TaskDTO]{
		Data: currentTask,
		ActionResult: services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		},
	})
}

// Handler method to create a new task
func (h *Handler) CreateTask(c *gin.Context) {
	var createTaskDTO models.CreateTaskDTO
	if err := c.ShouldBind(&createTaskDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.ActionResult{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Error:      []string{err.Error()},
		})
	}

	newId, er := h.Service.CreateTask(&createTaskDTO)
	if er != "" {
		c.JSON(http.StatusBadRequest, services.ActionResult{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Error:      []string{er},
		})
	}

	c.JSON(http.StatusOK, services.ActionResultModel[uint]{
		Data: newId,
		ActionResult: services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		},
	})

}

// Handler method to update an existing task by its ID
func (h *Handler) UpdateTask(c *gin.Context) {
	var updateTaskDTO models.UpdateTaskDTO
	if err := c.ShouldBind(&updateTaskDTO); err != nil {
		c.JSON(http.StatusBadRequest, services.ActionResult{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Error:      []string{err.Error()},
		})
	}
	id := c.Param("id")
	val, _ := strconv.Atoi(id)
	currentTask := h.Service.GetTaskById(val)
	if currentTask == (models.TaskDTO{}) {
		c.JSON(http.StatusNotFound, services.ActionResult{
			StatusCode: http.StatusNotFound,
			IsSuccess:  false,
			Error: []string{
				"Task not found",
			},
		})

	}
	isUpdated := h.Service.UpdateTask(val, &updateTaskDTO)
	if isUpdated {
		c.JSON(http.StatusOK, services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		})
	}
	c.JSON(http.StatusBadRequest, services.ActionResult{
		StatusCode: http.StatusBadRequest,
		IsSuccess:  false,
		Error: []string{
			"Failed to update task",
		},
	})
}

// Handler method to delete a task by its ID
func (h *Handler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	val, _ := strconv.Atoi(id)
	currentTask := h.Service.GetTaskById(val)
	if currentTask == (models.TaskDTO{}) {
		c.JSON(http.StatusNotFound, services.ActionResult{
			StatusCode: http.StatusNotFound,
			IsSuccess:  false,
			Error: []string{
				"Task not found",
			},
		})

	}
	isDeleted := h.Service.DeleteTask(val)
	if isDeleted {
		c.JSON(http.StatusOK, services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		})

	}
	c.JSON(http.StatusBadRequest, services.ActionResult{
		StatusCode: http.StatusBadRequest,
		IsSuccess:  false,
		Error: []string{
			"Failed to delete task",
		},
	})
}

// Handler method to retrieve all tasks with pagination
func (h *Handler) GetAllTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	allTasks := h.Service.GetAllTasks(page, pageSize)

	c.JSON(http.StatusOK, services.ActionResultModel[models.PaginatedDTO]{
		Data: allTasks,
		ActionResult: services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		},
	})
}

// Handler method to filter tasks based on query parameters such as year, month, and day
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
		c.JSON(http.StatusOK, services.ActionResultModel[models.PaginatedDTO]{
			Data: tasks,
			ActionResult: services.ActionResult{
				StatusCode: http.StatusOK,
				IsSuccess:  true,
				Error:      nil,
			},
		})
		return
	}

	if year != "" && month != "" {
		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		tasks := h.Service.GetTaskByMonthAndYear(m, y, page, pageSize)
		c.JSON(http.StatusOK, services.ActionResultModel[models.PaginatedDTO]{
			Data: tasks,
			ActionResult: services.ActionResult{
				StatusCode: http.StatusOK,
				IsSuccess:  true,
				Error:      nil,
			},
		})
		return
	}

	if year != "" {
		y, _ := strconv.Atoi(year)
		tasks := h.Service.GetTaskByYear(y, page, pageSize)
		c.JSON(http.StatusOK, services.ActionResultModel[models.PaginatedDTO]{
			Data: tasks,
			ActionResult: services.ActionResult{
				StatusCode: http.StatusOK,
				IsSuccess:  true,
				Error:      nil,
			},
		})
		return

	}

	tasks := h.Service.GetAllTasks(page, pageSize)
	c.JSON(http.StatusOK, services.ActionResultModel[models.PaginatedDTO]{
		Data: tasks,
		ActionResult: services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		},
	})
}

// Handler method to search for tasks by name with pagination
func (h *Handler) SearchByTask(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	name := c.Query("name")

	tasks := h.Service.SearchTaskByName(name, page, pageSize)

	c.JSON(http.StatusOK, services.ActionResultModel[models.PaginatedDTO]{
		Data: tasks,
		ActionResult: services.ActionResult{
			StatusCode: http.StatusOK,
			IsSuccess:  true,
			Error:      nil,
		},
	})

}

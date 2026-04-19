package models

import "fmt"

func mapTaskToDTO(task Task) TaskPaginationDTO {
	var dueDate string

	if task.DueDate != nil {
		dueDate = task.DueDate.Format("2006-01-02")
	}
	return TaskPaginationDTO{
		ID:          fmt.Sprintf("%d", task.ID),
		Name:        task.Name,
		Description: task.Description,
		DueDate:     dueDate,
		CreatedDate: task.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToPaginatedDTO(tasks []Task, total int, page int, pageSize int) PaginatedDTO {
	var dtos []TaskPaginationDTO
	for _, task := range tasks {
		dtos = append(dtos, mapTaskToDTO(task))
	}
	return PaginatedDTO{
		Tasks:    dtos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

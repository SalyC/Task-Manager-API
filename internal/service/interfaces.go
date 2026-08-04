package service

import "taskmanager/internal/models"

type ProjectRepositoryInterface interface {
	Create(project *models.Project) error
	GetByID(id int) (*models.Project, error)
	GetAll() ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id int) error
}

type TaskRepositoryInterface interface {
	Create(task *models.Task) error
	GetByID(id int) (*models.Task, error)
	GetAll(projectID *int, status *string) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id int) error
}

type CommentRepositoryInterface interface {
	Create(comment *models.Comment) error
	GetByID(id int) (*models.Comment, error)
	GetByTaskID(taskID int) ([]models.Comment, error)
	Delete(id int) error
}

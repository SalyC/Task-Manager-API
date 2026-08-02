package service

import (
	"errors"
	"taskmanager/internal/models"
	"taskmanager/internal/repository"
)

type TaskService struct {
	taskRepo   *repository.TaskRepository
	projectSvc *ProjectService
}

func NewTaskService(taskRepo *repository.TaskRepository, projectSvc *ProjectService) *TaskService {
	return &TaskService{
		taskRepo:   taskRepo,
		projectSvc: projectSvc,
	}
}

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskTitleRequired = errors.New("task title is required")
)

func (t *TaskService) Create(task *models.Task) error {
	if task.Title == "" {
		return ErrTaskTitleRequired
	}
	if task.ProjectID != 0 {
		_, err := t.projectSvc.GetByID(task.ProjectID)
		if err != nil {
			return err
		}
	}
	return t.taskRepo.Create(task)
}

func (t *TaskService) GetByID(id int) (*models.Task, error) {
	task, err := t.taskRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (t *TaskService) GetAll(projectID *int, status *string) ([]models.Task, error) {
	tasks, err := t.taskRepo.GetAll(projectID, status)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskService) Update(task *models.Task) error {
	if _, err := t.GetByID(task.ID); err != nil {
		return err
	}
	if task.Title == "" {
		return ErrTaskTitleRequired
	}
	if task.Status != "new" && task.Status != "in_progress" && task.Status != "done" {
		return errors.New("invalid status")
	}
	return t.taskRepo.Update(task)
}

func (t *TaskService) Delete(id int) error {
	if _, err := t.GetByID(id); err != nil {
		return err
	}
	return t.taskRepo.Delete(id)
}

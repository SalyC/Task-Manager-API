package service

import (
	"errors"
	"taskmanager/internal/models"
	"taskmanager/internal/repository"
)

type ProjectService struct {
	repo ProjectRepositoryInterface
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectNameRequired = errors.New("project name is required")
)

func (p *ProjectService) Create(project *models.Project) error {
	if project.Name == "" {
		return ErrProjectNameRequired
	}
	return p.repo.Create(project)
}
func (p *ProjectService) GetByID(id int) (*models.Project, error) {
	project, err := p.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}
func (p *ProjectService) GetAll() ([]models.Project, error) {
	return p.repo.GetAll()
}
func (p *ProjectService) Update(project *models.Project) error {
	if _, err := p.GetByID(project.ID); err != nil {
		return err
	}
	if project.Name == "" {
		return ErrProjectNameRequired
	}

	return p.repo.Update(project)
}
func (p *ProjectService) Delete(id int) error {
	_, err := p.GetByID(id)
	if err != nil {
		return err
	}
	return p.repo.Delete(id)
}

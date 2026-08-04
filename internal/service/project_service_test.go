package service

import (
	"errors"
	"taskmanager/internal/models"
	"testing"
)

type mockProjectRepository struct {
	createFunc  func(project *models.Project) error
	getByIDFunc func(id int) (*models.Project, error)
	getAllFunc  func() ([]models.Project, error)
	updateFunc  func(project *models.Project) error
	deleteFunc  func(id int) error
}

func (m *mockProjectRepository) Create(project *models.Project) error {
	if m.createFunc != nil {
		return m.createFunc(project)
	}
	return nil
}
func (m *mockProjectRepository) GetByID(id int) (*models.Project, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}
func (m *mockProjectRepository) GetAll() ([]models.Project, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return []models.Project{}, nil
}
func (m *mockProjectRepository) Update(project *models.Project) error {
	if m.updateFunc != nil {
		return m.updateFunc(project)
	}
	return nil
}
func (m *mockProjectRepository) Delete(id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func TestProjectService_Create(t *testing.T) {
	tests := []struct {
		name    string
		project *models.Project
		mockErr error
		wantErr bool
	}{
		{
			name:    "успешное создание",
			project: &models.Project{Name: "Test", Description: "Desc"},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "пустое имя",
			project: &models.Project{Name: "", Description: "Desc"},
			mockErr: nil,
			wantErr: true,
		},
		{
			name:    "ошибка репозитория",
			project: &models.Project{Name: "Test", Description: "Desc"},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				createFunc: func(p *models.Project) error {
					return tt.mockErr
				},
			}
			svc := NewProjectService(mockRepo)
			err := svc.Create(tt.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestProjectService_GetById(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		mockProject *models.Project
		mockErr     error
		wantErr     bool
		wantProject *models.Project
	}{
		{
			name:        "успешное получение",
			id:          1,
			mockProject: &models.Project{ID: 1, Name: "Test", Description: "Desc"},
			mockErr:     nil,
			wantErr:     false,
			wantProject: &models.Project{ID: 1, Name: "Test", Description: "Desc"},
		},
		{
			name:        "проект не найден",
			id:          999,
			mockProject: nil,
			mockErr:     nil,
			wantErr:     true,
			wantProject: nil,
		},
		{
			name:        "ошибка репозитория",
			id:          1,
			mockProject: nil,
			mockErr:     errors.New("db error"),
			wantErr:     true,
			wantProject: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				getByIDFunc: func(id int) (*models.Project, error) {
					return tt.mockProject, tt.mockErr
				},
			}
			svc := NewProjectService(mockRepo)
			got, err := svc.GetByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("GetByID() got nil, want project")
				} else if tt.wantProject != nil {
					if got.ID != tt.wantProject.ID || got.Name != tt.wantProject.Name {
						t.Errorf("GetByID() got = %v, want %v", got, tt.wantProject)
					}
				}
			} else {
				if tt.mockProject == nil && tt.mockErr == nil {
					if err != ErrProjectNotFound {
						t.Errorf("GetByID() error = %v, want %v", err, ErrProjectNotFound)
					}
				}
			}
		})
	}
}

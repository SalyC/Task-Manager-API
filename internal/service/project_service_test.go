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
func TestProjectService_GetAll(t *testing.T) {
	tests := []struct {
		name         string
		mockProjects []models.Project
		mockErr      error
		wantErr      bool
		wantLen      int
	}{
		{
			name: "успешное получение",
			mockProjects: []models.Project{
				{ID: 1, Name: "Project A", Description: "Desc"},
				{ID: 2, Name: "Project B", Description: "Desc"},
			},
			mockErr: nil,
			wantErr: false,
			wantLen: 2,
		},
		{
			name:         "пустой список",
			mockProjects: []models.Project{},
			mockErr:      nil,
			wantErr:      false,
			wantLen:      0,
		},
		{
			name:         "ошибка репозитория",
			mockProjects: nil,
			mockErr:      errors.New("db error"),
			wantErr:      true,
			wantLen:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				getAllFunc: func() ([]models.Project, error) {
					return tt.mockProjects, tt.mockErr
				},
			}
			svc := NewProjectService(mockRepo)
			got, err := svc.GetAll()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != tt.wantLen {
					t.Errorf("GetAll() len = %v, want %v", len(got), tt.wantErr)
				}
			}
		})
	}
}

func TestProjectService_Update(t *testing.T) {
	tests := []struct {
		name               string
		project            *models.Project
		mockGetByIDProject *models.Project
		mockGetByIDError   error
		mockUpdateError    error
		wantErr            bool
	}{
		{
			name:               "успешный апдейт",
			project:            &models.Project{ID: 1, Name: "Test", Description: "Desc"},
			mockGetByIDProject: &models.Project{ID: 1, Name: "Old", Description: "Old"},
			mockGetByIDError:   nil,
			mockUpdateError:    nil,
			wantErr:            false,
		},
		{
			name:               "пустое имя",
			project:            &models.Project{ID: 1, Name: "", Description: "Desc"},
			mockGetByIDProject: &models.Project{ID: 1, Name: "Old", Description: "Old"},
			mockGetByIDError:   nil,
			mockUpdateError:    nil,
			wantErr:            true,
		},
		{
			name:               "проект не найден",
			project:            &models.Project{ID: 999, Name: "Test", Description: "Desc"},
			mockGetByIDProject: nil,
			mockGetByIDError:   nil,
			mockUpdateError:    nil,
			wantErr:            true,
		},
		{
			name:               "ошибка репозитория при Update",
			project:            &models.Project{ID: 1, Name: "Test", Description: "Desc"},
			mockGetByIDProject: &models.Project{ID: 1, Name: "Old", Description: "Old"},
			mockGetByIDError:   nil,
			mockUpdateError:    errors.New("db error"),
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				getByIDFunc: func(id int) (*models.Project, error) {
					return tt.mockGetByIDProject, tt.mockGetByIDError
				},
				updateFunc: func(project *models.Project) error {
					return tt.mockUpdateError
				},
			}
			svc := NewProjectService(mockRepo)
			err := svc.Update(tt.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == "пустое имя" && err != ErrProjectNameRequired {
				t.Errorf("Update() error = %v, want %v", err, ErrProjectNameRequired)
			}
			if tt.name == "проект не найден" && err != ErrProjectNotFound {
				t.Errorf("Update() error = %v, want %v", err, ErrProjectNotFound)
			}
		})
	}
}

func TestProjectService_Delete(t *testing.T) {
	tests := []struct {
		name               string
		id                 int
		mockGetByIDProject *models.Project
		mockGetByIDError   error
		mockDeleteError    error
		wantErr            bool
	}{
		{
			name:               "успешное удаление",
			id:                 1,
			mockGetByIDProject: &models.Project{ID: 1, Name: "Test", Description: "Desc"},
			mockGetByIDError:   nil,
			mockDeleteError:    nil,
			wantErr:            false,
		},
		{
			name:               "Проект не найден",
			id:                 999,
			mockGetByIDProject: nil,
			mockGetByIDError:   nil,
			mockDeleteError:    nil,
			wantErr:            true,
		},
		{
			name:               "Ошибка репозитория при удалении",
			id:                 1,
			mockGetByIDProject: &models.Project{ID: 1, Name: "Test", Description: "Desc"},
			mockGetByIDError:   nil,
			mockDeleteError:    errors.New("db error"),
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockProjectRepository{
				getByIDFunc: func(id int) (*models.Project, error) {
					return tt.mockGetByIDProject, tt.mockGetByIDError
				},
				deleteFunc: func(id int) error {
					return tt.mockDeleteError
				},
			}

			svc := NewProjectService(mockRepo)
			err := svc.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.name == "проект не найден" && err == ErrProjectNotFound {
					t.Errorf("Delete() error = %v, want %v", err, ErrProjectNotFound)
				}
			}
		})
	}
}

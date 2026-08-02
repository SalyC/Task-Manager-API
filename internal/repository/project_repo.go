package repository

import (
	"database/sql"
	"taskmanager/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	query := "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id, created_at"
	return r.db.QueryRow(query, project.Name, project.Description).Scan(&project.ID, &project.CreatedAt)
}
func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	var p models.Project
	query := "SELECT id, name, description, created_at FROM projects WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at FROM projects ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}
func (r *ProjectRepository) Update(project *models.Project) error {
	query := "UPDATE projects SET name = $1, description = $2 WHERE id = $3"
	_, err := r.db.Exec(query, project.Name, project.Description, project.ID)
	return err
}
func (r *ProjectRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM projects WHERE id = $1", id)
	return err
}

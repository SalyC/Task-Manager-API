package repository

import (
	"database/sql"
	"strconv"
	"taskmanager/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (t *TaskRepository) Create(task *models.Task) error {
	query := "INSERT INTO tasks (project_id, title, description, status, priority, deadline) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at"
	return t.db.QueryRow(query, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Deadline).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}
func (t *TaskRepository) GetByID(id int) (*models.Task, error) {
	var tm models.Task
	query := "SELECT id, project_id, title, description, status, priority, deadline, created_at, updated_at FROM tasks WHERE id = $1"
	err := t.db.QueryRow(query, id).Scan(&tm.ID, &tm.ProjectID, &tm.Title, &tm.Description, &tm.Status, &tm.Priority, &tm.Deadline, &tm.CreatedAt, &tm.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tm, nil
}

func (t *TaskRepository) GetAll(projectID *int, status *string) ([]models.Task, error) {
	query := "SELECT id, project_id, title, description, status, priority, deadline, created_at, updated_at FROM tasks WHERE 1=1"
	args := []interface{}{}
	argInd := 1

	if projectID != nil {
		query += " AND project_id = $" + strconv.Itoa(argInd)
		args = append(args, *projectID)
		argInd++
	}
	if status != nil {
		query += " AND status = $" + strconv.Itoa(argInd)
		args = append(args, *status)
		argInd++
	}
	query += " ORDER BY id"

	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Deadline, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (t *TaskRepository) Update(task *models.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, deadline = $5, updated_at = NOW() WHERE id = $6"
	_, err := t.db.Exec(query, task.Title, task.Description, task.Status, task.Priority, task.Deadline, task.ID)
	return err
}
func (t *TaskRepository) Delete(id int) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = $1")
	return err
}

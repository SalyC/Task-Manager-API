package models

import "time"

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Task struct {
	ID          int        `json:"id"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	Author    string    `json:"author"`
	Content   string    `json:"context"`
	CreatedAt time.Time `json:"created_at"`
}

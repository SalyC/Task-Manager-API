package models

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateTaskRequest struct {
	ProjectID   int     `json:"project_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status,omitempty"`
	Priority    string  `json:"priority,omitempty"`
	Deadline    *string `json:"deadline,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	Deadline    *string `json:"deadline,omitempty"`
}

type CreateCommentRequest struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

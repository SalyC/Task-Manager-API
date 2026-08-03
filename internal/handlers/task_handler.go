package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"taskmanager/internal/models"
	"taskmanager/internal/service"
	"time"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{service: svc}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProjectID   int     `json:"project_id"`
		Title       string  `json:"title"`
		Description string  `json:"description,omitempty"`
		Status      string  `json:"status,omitempty"`
		Priority    string  `json:"priority,omitempty"`
		Deadline    *string `json:"deadline,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.ProjectID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Project ID is required and must be > 0")
		return
	}
	if input.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Task title is required")
		return
	}
	if input.Status != "" && input.Status != "new" && input.Status != "in_progress" && input.Status != "done" {
		respondWithError(w, http.StatusBadRequest, "Invalid status")
		return
	}
	if input.Priority != "" && input.Priority != "low" && input.Priority != "medium" && input.Priority != "high" {
		respondWithError(w, http.StatusBadRequest, "Invalid priority")
		return
	}

	var deadline *time.Time
	if input.Deadline != nil && *input.Deadline != "" {
		parsed, err := time.Parse("2006-01-02", *input.Deadline)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid deadline format, use YYYY-MM-DD")
			return
		}
		deadline = &parsed
	}

	task := &models.Task{
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		Deadline:    deadline,
	}

	if err := h.service.Create(task); err != nil {
		if err == service.ErrTaskTitleRequired {
			respondWithError(w, http.StatusBadRequest, "Task title is required")
			return
		}
		if err == service.ErrProjectNotFound {
			respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	task, err := h.service.GetByID(id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get task")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	statusStr := r.URL.Query().Get("status")

	var projectID *int
	if projectIDStr != "" {
		id, err := strconv.Atoi(projectIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid project_id")
			return
		}
		projectID = &id
	}

	var status *string
	if statusStr != "" {
		if statusStr != "new" && statusStr != "in_progress" && statusStr != "done" {
			respondWithError(w, http.StatusBadRequest, "Invalid status")
			return
		}
		status = &statusStr
	}

	tasks, err := h.service.GetAll(projectID, status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
		Priority    *string `json:"priority"`
		Deadline    *string `json:"deadline"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	existing, err := h.service.GetByID(id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get task")
		return
	}

	if input.Title != nil {
		if *input.Title == "" {
			respondWithError(w, http.StatusBadRequest, "Task title cannot be empty")
			return
		}
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Status != nil {
		status := *input.Status
		if status != "new" && status != "in_progress" && status != "done" {
			respondWithError(w, http.StatusBadRequest, "Invalid status")
			return
		}
		existing.Status = status
	}
	if input.Priority != nil {
		priority := *input.Priority
		if priority != "low" && priority != "medium" && priority != "high" {
			respondWithError(w, http.StatusBadRequest, "Invalid priority")
			return
		}
		existing.Priority = priority
	}
	if input.Deadline != nil {
		if *input.Deadline == "" {
			existing.Deadline = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *input.Deadline)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid deadline format, use YYYY-MM-DD")
				return
			}
			existing.Deadline = &parsed
		}
	}

	if err := h.service.Update(existing); err != nil {
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		if err == service.ErrTaskTitleRequired {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	respondWithJSON(w, http.StatusOK, existing)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

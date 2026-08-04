package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"taskmanager/internal/models"
	"taskmanager/internal/service"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{service: svc}
}

//! Если это читать кто то будет, то нумерацию ошибок я добавил для себя как эксперимент, лучше не стоит так делать

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["taskId"]
	taskId, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error 400: Invalid ID")
		return
	}
	if _, err := h.service.GetByTaskID(taskId); err != nil {
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Error 404: Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error 500: Failed to check task")
		return
	}
	var input struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error 400: Invalid JSON")
		return
	}
	log.Printf("Decoded input: %+v", input)

	if _, err := h.service.GetByTaskID(taskId); err != nil {
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to check task")
		return
	}

	if input.Author == "" {
		respondWithError(w, http.StatusBadRequest, "Error 400: Author is required")
		return
	}
	if input.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Error 400: Comment content is required")
		return
	}

	comment := &models.Comment{
		TaskID:  taskId,
		Author:  input.Author,
		Content: input.Content,
	}

	if err := h.service.Create(comment); err != nil {
		if err == service.ErrCommentContentRequired {
			respondWithError(w, http.StatusBadRequest, "Error 400: Content is required")
			return
		}
		if err == service.ErrTaskNotFound {
			respondWithError(w, http.StatusNotFound, "Error 404: Task not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error 500: Failed to get task")
		return
	}
	respondWithJSON(w, http.StatusCreated, comment)
}
func (h *CommentHandler) GetCommentsByTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["taskId"]
	taskId, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error 400: Invalid ID")
		return
	}
	comments, err := h.service.GetByTaskID(taskId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error 500: Failed to get task")
		return
	}
	respondWithJSON(w, http.StatusOK, comments)
}
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error 400: Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		if err == service.ErrCommentNotFound {
			respondWithError(w, http.StatusNotFound, "Error 404: Comment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error 500: Failed to delete comment")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

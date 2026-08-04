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

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: svc}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if input.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Project name is required")
		return
	}

	project := &models.Project{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := h.service.Create(project); err != nil {
		log.Printf("ERROR creating project: %v", err)
		if err == service.ErrProjectNameRequired {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	respondWithJSON(w, http.StatusCreated, project)
}
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.service.GetByID(id)
	if err != nil {
		if err != service.ErrProjectNotFound {
			respondWithError(w, http.StatusBadRequest, "Project not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, "Failed to get project")
		return
	}
	respondWithJSON(w, http.StatusOK, project)
}
func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to get projects")
		return
	}
	respondWithJSON(w, http.StatusOK, projects)
}
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if input.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Project name is required")
		return
	}
	project := &models.Project{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := h.service.Update(project); err != nil {
		if err == service.ErrProjectNotFound {
			respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		if err == service.ErrProjectNameRequired {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	updated, err := h.service.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to update project")
		return
	}
	respondWithJSON(w, http.StatusOK, updated)
}
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == service.ErrProjectNotFound {
			respondWithError(w, http.StatusBadRequest, "Project not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

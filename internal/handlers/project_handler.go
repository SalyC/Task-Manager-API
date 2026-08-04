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

// CreateProject godoc
// @Summary      Создать новый проект
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        request body models.CreateProjectRequest true "Данные проекта"
// @Success      201  {object}  models.Project
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects [post]
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

// GetProject godoc
// @Summary      Получить проект по ID
// @Tags         projects
// @Produce      json
// @Param        id path int true "ID проекта"
// @Success      200  {object}  models.Project
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id} [get]
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

// GetAllProjects godoc
// @Summary      Получить все проекты
// @Tags         projects
// @Produce      json
// @Success      200  {array}  models.Project
// @Failure      500  {object}  map[string]string
// @Router       /projects [get]
func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to get projects")
		return
	}
	respondWithJSON(w, http.StatusOK, projects)
}

// UpdateProject godoc
// @Summary      Обновить проект по ID
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id path int true "ID проекта"
// @Param        request body models.UpdateProjectRequest true "Данные для обновления"
// @Success      200  {object}  models.Project
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id} [put]
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

// DeleteProject godoc
// @Summary      Удалить проект
// @Tags         projects
// @Param        id path int true "ID проекта"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id} [delete]
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

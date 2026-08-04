package main

import (
	"database/sql"
	"log"
	"net/http"

	"taskmanager/internal/config"
	"taskmanager/internal/handlers"
	"taskmanager/internal/repository"
	"taskmanager/internal/service"

	_ "taskmanager/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// @title           Task Manager API
// @version         1.0
// @description      API для управления проектами, задачами и комментариями
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("Database connected")

	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	projectSvc := service.NewProjectService(projectRepo)
	taskSvc := service.NewTaskService(taskRepo, projectSvc)
	commentSvc := service.NewCommentService(commentRepo, taskSvc)

	projectHandler := handlers.NewProjectHandler(projectSvc)
	taskHandler := handlers.NewTaskHandler(taskSvc)
	commentHandler := handlers.NewCommentHandler(commentSvc)

	r := mux.NewRouter()

	//! Сваггер
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//! роутинг по проектам
	r.HandleFunc("/projects", projectHandler.GetAllProjects).Methods("GET")
	r.HandleFunc("/projects", projectHandler.CreateProject).Methods("POST")
	r.HandleFunc("/projects/{id}", projectHandler.GetProject).Methods("GET")
	r.HandleFunc("/projects/{id}", projectHandler.UpdateProject).Methods("PUT")
	r.HandleFunc("/projects/{id}", projectHandler.DeleteProject).Methods("DELETE")

	//! роутинг по задачам
	r.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")

	//! роутинг по комментам
	r.HandleFunc("/tasks/{taskId}/comments", commentHandler.GetCommentsByTask).Methods("GET")
	r.HandleFunc("/tasks/{taskId}/comments", commentHandler.CreateComment).Methods("POST")
	r.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")

	log.Printf("Server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}

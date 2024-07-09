package main

import (
	_ "HL_todo_list/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Todo List API
// @version 1.0
// @description This is a simple Todo List API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath
func main() {
	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/api/todo-list/health", HealthHandler).Methods("GET")
	r.HandleFunc("/api/todo-list/tasks", CreateTask).Methods("POST")
	r.HandleFunc("/api/todo-list/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/api/todo-list/tasks/{id}", DeleteTask).Methods("DELETE")
	r.HandleFunc("/api/todo-list/tasks/{id}/done", MarkTaskDone).Methods("PUT")
	r.HandleFunc("/api/todo-list/tasks", GetTasksByStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

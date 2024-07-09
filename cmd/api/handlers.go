package main

import (
	"HL_todo_list/pkg/models"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var tasks = []models.Task{}
var id = 0
var validate = validator.New()

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the given title and active date
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task to create"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Task already exists"
// @Router /api/todo-list/tasks [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, t := range tasks {
		if t.Title == task.Title && t.ActiveAt == task.ActiveAt {
			http.Error(w, "task already exists", http.StatusNotFound)
			return
		}

	}

	if err := validate.Struct(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.ID = generateID()
	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": task.ID})
}

// UpdateTask godoc
// @Summary Update an existing task
// @Description Update an existing task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.Task true "Updated task data"
// @Success 204
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{id} [put]
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatedTask
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 204
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound)
}

// MarkTaskDone godoc
// @Summary Mark a task as done
// @Description Mark a task as done by ID
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 204
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{id}/done [put]
func MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "task not found", http.StatusNotFound)
}

// GetTasksByStatus godoc
// @Summary Get tasks by status
// @Description Get tasks by their status (active or done)
// @Tags tasks
// @Param status query string false "Task status" Enums(active, done) default(active)
// @Produce json
// @Success 200 {array} models.Task
// @Router /api/todo-list/tasks [get]
func GetTasksByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	currentDate := time.Now().Format("2006-01-02")
	var filteredTasks []models.Task
	for _, task := range tasks {
		isActive := task.ActiveAt <= currentDate
		if (status == "done" && task.Done) || (status == "active" && !task.Done && isActive) {
			date, _ := time.Parse("2006-01-02", task.ActiveAt)
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				task.Title = "ВЫХОДНОЙ - " + task.Title
			}
			filteredTasks = append(filteredTasks, task)
		}
	}

	sort.Slice(filteredTasks, func(i, j int) bool {
		return filteredTasks[i].ID < filteredTasks[j].ID
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredTasks)
}

func generateID() string {
	id++
	return strconv.Itoa(id)
}

// HealthHandler godoc
// @Summary Health check
// @Description Check the health of the service
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

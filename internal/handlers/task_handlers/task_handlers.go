package task_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"main/internal/models"
)

var tasks = make(map[uint]*models.Task)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// StartTask starts a task
// @Summary Start a task
// @Tags tasks
// @Description Starts a task for the user
// @ID start-task
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} ErrorResponse
// @Router /tasks/start [get]
func StartTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, exists := tasks[uint(id)]
	if !exists {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}

	task.StartTask()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// EndTask ends a task
// @Summary End a task
// @Description Ends a task for the user
// @ID end-task
// @Tags tasks
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} ErrorResponse
// @Router /tasks/end [get]
func EndTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, exists := tasks[uint(id)]
	if !exists {
		http.Error(w, "Task not found", http.StatusBadRequest)
		return
	}

	task.EndTask()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KMrsR/task-manager/internal/httputils"
	"github.com/KMrsR/task-manager/internal/models"
	"github.com/KMrsR/task-manager/internal/storage"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	storage storage.TaskStorage
}

func NewTaskHandler(s storage.TaskStorage) *TaskHandler {
	return &TaskHandler{storage: s}
}

// хэндлер создания задачи
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		httputils.ResponseWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.storage.AddTask(task); err != nil {
		httputils.ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httputils.RespondWithJSON(w, http.StatusCreated, task)
}

// хэндлер получить задачу по id
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// вытащим данные из url path
	id := mux.Vars(r)["id"]
	if id == "" {
		httputils.ResponseWithError(w, "id required", http.StatusBadRequest)
		return
	}
	task, err := h.storage.GetTaskByID(id)
	if err != nil {
		httputils.ResponseWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, task)
}

// хэндлер получить все задачи
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputils.ResponseWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tasks, err := h.storage.GetTasks()
	if err != nil {
		httputils.ResponseWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, tasks)
}

// хэндлер обновления задачи
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var updatedTask models.Task
	// Ограничиваем размер тела запроса (защита от DoS)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
	// вытащим данные из url path
	id := mux.Vars(r)["id"]
	if id == "" {
		httputils.ResponseWithError(w, "id required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		httputils.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.storage.UpdateTask(id, updatedTask)
	if err != nil {
		httputils.ResponseWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, updatedTask)
}

// хэндлер удаления задачи
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// вытащим данные из url path
	id := mux.Vars(r)["id"]
	if id == "" {
		httputils.ResponseWithError(w, "id required", http.StatusBadRequest)
		return
	}
	err := h.storage.DeleteTask(id)
	if err != nil {
		httputils.ResponseWithError(w, err.Error(), http.StatusNotFound)
		return
	}
	// приуспешном удалении
	w.WriteHeader(http.StatusNoContent)
}

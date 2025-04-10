package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KMrsR/task-manager/internal/models"
	"github.com/KMrsR/task-manager/internal/storage"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	storage storage.TaskStorage
}

func NewTaskHadler(s storage.TaskStorage) *TaskHandler {
	return &TaskHandler{storage: s}
}

// хэндлер создания задачи
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		responseWithError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.storage.AddTask(task); err != nil {
		responseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// хэндлер получить задачу по id
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	// вытащим данные из url path
	id := mux.Vars(r)["id"]
	if id == "" {
		responseWithError(w, "id required", http.StatusBadRequest)
		return
	}
	task, err := h.storage.GetTaskByID(id)
	if err != nil {
		responseWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// хэндлер получить все задачи
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tasks, err := h.storage.GetTasks()
	if err != nil {
		responseWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

// хэндлер обновления задачи
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var updatedTask models.Task
	// Ограничиваем размер тела запроса (защита от DoS)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
	// вытащим данные из url path
	id := mux.Vars(r)["id"]
	if id == "" {
		responseWithError(w, "id required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		responseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.storage.UpdateTask(id, updatedTask)
	if err != nil {
		responseWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedTask)
}

// хэндлер удаления задачи
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// вытащим данные из url path
	id := mux.Vars(r)["id"]
	if id == "" {
		responseWithError(w, "id required", http.StatusBadRequest)
		return
	}
	err := h.storage.DeleteTask(id)
	if err != nil {
		responseWithError(w, err.Error(), http.StatusNotFound)
		return
	}
	// приуспешном удалении
	w.WriteHeader(http.StatusNoContent)
}

// функция дл стандартных ответов
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// функция для ответа с ошибкой
func responseWithError(w http.ResponseWriter, msg string, code int) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

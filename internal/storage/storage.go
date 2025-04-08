package storage

import (
	"fmt"
	"sync"

	"github.com/KMrsR/task-manager/internal/models"
)

type TaskStorage interface {
	GetTasks() ([]models.Task, error)
	GetTaskByID(id string) (*models.Task, error)
	AddTask(task models.Task) error
	UpdateTask(id string, updatedTask models.Task) error
	DeleteTask(id string) error
	Clear()
}

type MemoryStorage struct {
	tasks map[string]models.Task
	mu    sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{tasks: make(map[string]models.Task)}
}

// Возвращает все задачи
func (s *MemoryStorage) GetTasks() ([]models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l := len(s.tasks)
	if l == 0 {
		return nil, fmt.Errorf("no tasks found")
	}

	tasks := make([]models.Task, 0, l)
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Добавляет новую задачу
// возвращает ошибку если ID пуст
func (s *MemoryStorage) AddTask(task models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	s.tasks[task.ID] = task
	return nil
}

// Возвращает задачу по ID
// возвращает ошибку если ID не существует
func (s *MemoryStorage) GetTaskByID(id string) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ex := s.tasks[id]
	if !ex {
		return nil, fmt.Errorf("task with ID %v not found", id)
	}

	return &task, nil
}

// Обновляем задачу по ID
// возвращает ошибку ID не существует
func (s *MemoryStorage) UpdateTask(id string, updatedTask models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id != updatedTask.ID {
		return fmt.Errorf("task ID mismatch")
	}

	_, ex := s.tasks[id]
	if !ex {
		return fmt.Errorf("task with ID %v not found", id)
	}
	s.tasks[id] = updatedTask
	return nil
}

// Удаляем задачу по ID
// возвращает ошибку ID не существует или мапа пуста
func (s *MemoryStorage) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ex := s.tasks[id]
	if !ex {
		return fmt.Errorf("task with ID %v not found", id)
	}

	delete(s.tasks, id)
	return nil
}

// Метод для очистки хранилища
func (s *MemoryStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = make(map[string]models.Task)
}

package storage

import (
	"strconv"
	"sync"
	"testing"

	"github.com/KMrsR/task-manager/internal/models"
)

func TestAddTask(t *testing.T) {
	type testCase struct {
		name     string
		task     models.Task
		expected interface{}
	}
	store := NewMemoryStorage()
	testCases := []testCase{
		{
			name:     "empty ID",
			task:     models.Task{},
			expected: "task ID cannot be empty",
		},
		{
			name: "normal case",
			task: models.Task{ID: "1",
				Title:    "first task",
				Status:   "in_progress",
				Deadline: "2025-04-31"},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer store.Clear()
			err := store.AddTask(tc.task)
			// проверим добавление пустого ID
			if err != nil {
				if err.Error() != tc.expected {
					t.Errorf("TestAddTask %s failed, got %v, expected %v", tc.name, err, tc.expected)
				}
			}
			// проверим что все данные добавились корректно
			if store.tasks[tc.task.ID].ID != tc.task.ID ||
				store.tasks[tc.task.ID].Title != tc.task.Title ||
				store.tasks[tc.task.ID].Status != tc.task.Status ||
				store.tasks[tc.task.ID].Deadline != tc.task.Deadline {
				t.Errorf("tasks missmatch, expected %v, got %v", tc.task, store.tasks[tc.task.ID])
			}
		})
	}
}

func TestGetTasks(t *testing.T) {

	testTasks := []models.Task{
		{ID: "1",
			Title:    "first task",
			Status:   "in_progress",
			Deadline: "2025-04-31"},
		{ID: "2",
			Title:    "second task",
			Status:   "in_progress",
			Deadline: "2025-04-30"},
	}

	store := NewMemoryStorage()
	defer store.Clear()

	t.Run("empty storage", func(t *testing.T) {
		tasks, err := store.GetTasks()
		// проверка получения ожидаемых ошибок
		if err == nil || err.Error() != "no tasks found" {
			t.Errorf("Expected 'no tasks found', got %v", err)
		}
		// проверим что список задач пуст
		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks, got %v", len(tasks))
		}
	})

	//добавляем тестовые задачи в хранилище
	for _, v := range testTasks {
		err := store.AddTask(v)
		if err != nil {
			t.Errorf("failed to add task^ %v", v)
		}
	}
	// получаем задачи из хранилища
	tasks, err := store.GetTasks()
	if err != nil {
		t.Errorf("failed to get tasks")
	}

	t.Run("storage with tasks", func(t *testing.T) {
		// проверим что длинна полченного слайса задач одинаков с исходным
		if len(tasks) != len(testTasks) {
			t.Errorf("expected %d tasks, got %d", len(testTasks), len(tasks))
		}
		// проверим что содержимое исходного и полученного слайса задач одинаков
		for i, task := range tasks {
			if task.ID != testTasks[i].ID ||
				task.Title != testTasks[i].Title ||
				task.Status != testTasks[i].Status ||
				task.Deadline != testTasks[i].Deadline {
				t.Errorf("tasks %d missmatch: expected %v, got%v", i, testTasks[i], task)
			}
		}
	})

	t.Run("get task by ID", func(t *testing.T) {
		// получим все задачи по ID
		for _, task := range tasks {
			taskByID, err := store.GetTaskByID(task.ID)
			if err != nil {
				t.Errorf("failed to get task bi ID")
			}
			// проверим что полученные данные верны
			if taskByID.ID != task.ID ||
				taskByID.Status != task.Status ||
				taskByID.Deadline != task.Deadline ||
				taskByID.Title != task.Title {
				t.Errorf("task missmatch, got %v, expected %v", taskByID, task)

			}
		}
		// проверим получения по несуществующему ID
		_, err := store.GetTaskByID("badID")
		if err != nil {
			if err.Error() != "task with ID badID not found" {
				t.Errorf("testing badId failed")
			}
		}
	})
}

func TestUpdateAndDeleteTask(t *testing.T) {
	task := models.Task{ID: "1",
		Title:    "first task",
		Status:   "in_progress",
		Deadline: "2025-04-31"}
	stroe := NewMemoryStorage()
	defer stroe.Clear()
	stroe.AddTask(task)
	// проверим что ID и задача не совпадают
	t.Run("task ID mismath", func(t *testing.T) {
		err := stroe.UpdateTask("badID", models.Task{ID: "2", Title: "Updated first task"})
		if err != nil {
			if err.Error() != "task ID mismatch" {
				t.Errorf("expected task ID mismatch, got %v", err)
			}
		}
	})
	// проверим обновление несуществующей задачи
	t.Run("task with ID not found", func(t *testing.T) {
		err := stroe.UpdateTask("2", models.Task{ID: "2", Title: "Updated first task"})
		if err != nil {
			if err.Error() != "task with ID 2 not found" {
				t.Errorf("expected task with ID 2 not found, got %v", err)
			}
		}
	})
	// проверим обновление задачи
	t.Run("updatedTask", func(t *testing.T) {
		err := stroe.UpdateTask("1", models.Task{ID: "1", Title: "Updated first task"})
		if err != nil {
			t.Fatal(err)
		}
		if stroe.tasks[task.ID].Title != "Updated first task" {
			t.Errorf("task updated with wrong data")
		}
	})

	// проверим удаление по ID
	t.Run("delete task by ID", func(t *testing.T) {
		// попытаемся удалить по несуществующему ID
		err := stroe.DeleteTask("2")
		if err != nil {
			if err.Error() != "task with ID 2 not found" {
				t.Errorf("expected task with ID 2 not found, got %v", err)
			}
		}
		// нормальное удаление
		err = stroe.DeleteTask(task.ID)
		if err != nil {
			t.Errorf("task deleting failed")
		}
		//проверим что задача удалена
		_, err = stroe.GetTaskByID(task.ID)
		if err == nil {
			t.Errorf("task deleting failed")
		}
	})
}

// тест на конкурентность
func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	store := NewMemoryStorage()
	cnt := 100
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := store.AddTask(models.Task{ID: strconv.Itoa(id)})
			if err != nil {
				t.Errorf("failed to add task")
			}
		}(i)
	}
	wg.Wait()
	// проверим что в хранилище 100 задач
	if len(store.tasks) != cnt {
		t.Errorf("cuncurrency test failed")
	}
}

/*
подсказки
# 1. Генерируем файл покрытия
go test -coverprofile coverage.out

# 2. Создаем HTML-отчет (вариант для PowerShell)
go tool cover -html="coverage.out" -o "coverage.html"

# 3. Открываем в браузере
Start-Process coverage.html

*/

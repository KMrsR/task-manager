package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KMrsR/task-manager/internal/handlers"
	"github.com/KMrsR/task-manager/internal/middleware"
	"github.com/KMrsR/task-manager/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	// Чтение конфига из env
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Инициализация хранилища
	storage, err := storage.NewPostgresStorage(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer storage.Close()

	// store := storage.NewMemoryStorage()
	handler := handlers.NewTaskHandler(storage)

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Handle("/tasks", middleware.RequireJSON(http.HandlerFunc(handler.CreateTask))).Methods("POST")
	router.HandleFunc("/tasks", handler.GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", handler.GetTaskByID).Methods("GET")
	router.Handle("/task/{id}", middleware.RequireJSON(http.HandlerFunc(handler.UpdateTask))).Methods("PUT")
	router.HandleFunc("/task/{id}", handler.DeleteTask).Methods("DELETE")

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

package main

import (
	"log"
	"net/http"

	"github.com/KMrsR/task-manager/internal/handlers"
	"github.com/KMrsR/task-manager/internal/middleware"
	"github.com/KMrsR/task-manager/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	store := storage.NewMemoryStorage()
	handler := handlers.NewTaskHandler(store)

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

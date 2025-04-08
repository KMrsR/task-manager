package models

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`   //"new","in_progress","done"
	Deadline string `json:"deadline"` //"YYYY-MM-DD"
}

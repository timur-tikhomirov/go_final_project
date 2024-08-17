package main

const (
	DefaultPort = "7540"
	DateFormat  = "20060102"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

var ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

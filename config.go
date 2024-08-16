package main

const (
	DefaultPort = "7540"
	DateFormat  = "20060102"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

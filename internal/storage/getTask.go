package storage

import (
	"fmt"

	"github.com/timur-tikhomirov/go_final_project/configs"
)

// Получение задачи по id
func (s *Store) GetTask(id string) (configs.Task, error) {
	var t configs.Task
	if id == "" {
		return configs.Task{}, fmt.Errorf(`{"error":"Не указан id"}`)
	}
	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return configs.Task{}, fmt.Errorf(`{"error":"Задача не найдена"}`)
	}
	return t, nil
}

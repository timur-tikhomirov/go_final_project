package storage

import (
	"fmt"
	"time"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/tasks"
)

// Выполнение задачи
func (s *Store) TaskDone(id string) error {
	var task configs.Task
	// Проверяем id
	if id == "" {
		return fmt.Errorf(`{"error":"Не указан id"}`)
	}

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return fmt.Errorf(`{"error":"Задача не найдена"}`)
	}
	if task.Repeat == "" {
		_, err := s.db.Exec("DELETE FROM scheduler WHERE id=?", task.ID)
		if err != nil {
			return fmt.Errorf(`{"error":"Не удалось удалить задачу"}`)
		}
	} else {
		next, err := tasks.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
		}
		task.Date = next
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = s.db.Exec(query, task.Date, task.ID)
	if err != nil {
		return fmt.Errorf(`{"error":"Ошибка обновления даты"}`)
	}
	return nil
}

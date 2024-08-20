package storage

import (
	"fmt"
	"time"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/tasks"
)

// Добавление задачи
func (s *Store) PostTask(t configs.Task) (string, error) {
	var err error

	if t.Title == "" {
		return "", fmt.Errorf(`{"error":"Не указан заголовок задачи"}`)
	}

	// Проверяем наличие даты
	if t.Date == "" {
		t.Date = time.Now().Format(configs.DateFormat)
	}

	_, err = time.Parse(configs.DateFormat, t.Date)
	if err != nil {
		return "", fmt.Errorf(`{"error":"Некорректный формат даты"}`)
	}
	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if t.Date < time.Now().Format(configs.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := tasks.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return "", fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(configs.DateFormat)
		}
	}

	// Добавляем задачу в базу
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", fmt.Errorf(`{"error":"Не удалось добавить задачу"}`)
	}

	// Возвращаем идентификатор добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(`{"error":"Не удалось вернуть id новой задачи"}`)
	}
	return fmt.Sprintf("%d", id), nil
}

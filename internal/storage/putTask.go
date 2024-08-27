package storage

import (
	"fmt"
	"time"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/tasks"
)

// Редактирование задачи
func (s *Store) PutTask(t configs.Task) error {
	// Проверяем, что id в теле не пустой
	if t.ID == "" {
		return fmt.Errorf(`{"error":"Не указан id"}`)
	}
	// Проверяем обязательное поле Title
	if t.Title == "" {
		return fmt.Errorf(`{"error":"Не указан заголовок задачи"}`)
	}
	// Проверяем дату
	if t.Date == "" {
		t.Date = time.Now().Format(configs.DateFormat)
	}

	_, err := time.Parse(configs.DateFormat, t.Date)
	if err != nil {
		return fmt.Errorf(`{"error":"Некорректный формат даты"}`)
	}

	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if t.Date < time.Now().Format(configs.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := tasks.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {

				return fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(configs.DateFormat)
		}
	}

	// Обновляем задачу в базе
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {

		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`) // вот эта вот история не работает, посчитаем измененные ряды
	}

	// Считаем измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`{"error":"Не удалось посчитать измененные строки"}`)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`)
	}
	// Если ошибок не возникло, возвращаем nil
	return nil
}

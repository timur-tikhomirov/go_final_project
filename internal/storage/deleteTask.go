package storage

import "fmt"

// Удаление задачи из БД
func (s *Store) DeleteTask(id string) error {
	// Проверяем, что id не пустой
	if id == "" {
		return fmt.Errorf(`{"error":"Не указан id"}`)
	}
	query := "DELETE FROM scheduler WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(`{"error":"Не удалось удалить задачу"}`)
	}
	// Считаем измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf(`{"error":"Не удалось посчитать измененные строки"}`)
	}

	if rowsAffected == 0 {

		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`)
	}
	// Если ошибок нет, возвращаем nil
	return nil
}

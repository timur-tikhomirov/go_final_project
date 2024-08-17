package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Обработчик для task/done
func TaskDoneHandler(res http.ResponseWriter, req *http.Request) {

	var task Task

	// Получаем id из запроса
	id := req.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if id == "" {
		err := errors.New("id не указан")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}
	var err error

	// Получаем задачу из базы данных для дальнейших операций

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		err := errors.New("Задача с таким id не найдена")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return

	}

	// Если указано правило, вычисляем следующую дату
	if task.Repeat != "" {

		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			err := errors.New("Правило повторения указано в неправильном формате")
			ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(ErrorResponse)
			return

		}

		// Обновляем дату задачи
		query := `UPDATE scheduler SET date = ? WHERE id = ?`
		if _, err = DB.Exec(query, nextDate, id); err != nil {
			err := errors.New("Не удалось обновить дату задачи")
			ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(ErrorResponse)
			return
		}
	} else {
		// Если правило не указано, удаляем задачу
		query := `DELETE FROM scheduler WHERE id = ?`
		if _, err = DB.Exec(query, id); err != nil {
			err := errors.New("Не удалось удалить задачу")
			ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(ErrorResponse)
			return
		}
	}

	// Возвращаем ответ
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(nil)
}

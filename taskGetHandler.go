package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Обработчик GET для /task
func TaskGetHandler(res http.ResponseWriter, req *http.Request) {
	var task Task
	// Получаем  id из запроса
	id := req.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if id == "" {
		err := errors.New("id не указан")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Запрос в БД
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		err := errors.New("Задача с таким id не найдена")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return

	}

	// Возвращаем ответ
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(task)
}

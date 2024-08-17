package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Обработчик DELETE для /task
func TaskDeleteHandler(res http.ResponseWriter, req *http.Request) {
	// Получаем значение параметра id из запроса
	id := req.URL.Query().Get("id")

	// Проверяем, что id не пустой
	if id == "" {
		err := errors.New("id не указан")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Удаляем задачу из базы данных
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(query, id)
	if err != nil {
		err := errors.New("Задача с таким id не найдена")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Считаем измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err := errors.New("Не получилось посчитать измененные строки")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	if rowsAffected == 0 {
		err := errors.New("Задача с таким id не найдена")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Возвращаем ответ
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(nil)
}

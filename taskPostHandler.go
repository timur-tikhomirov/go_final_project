package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Обработчик POST для /task
func TaskPostHandler(res http.ResponseWriter, req *http.Request) {
	var task Task

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Получаем тело запроса
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		err := errors.New("Ошибка десериализации JSON")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	defer req.Body.Close()

	// Проверяем обязательное поле Title
	if task.Title == "" {
		err := errors.New("Не указан заголовок задачи")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Проверяем наличие даты
	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	}

	_, err = time.Parse(DateFormat, task.Date)
	if err != nil {
		err := errors.New("Некорректный формат даты")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if task.Date < time.Now().Format(DateFormat) {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				err := errors.New("Правило повторения указано в неправильном формате")
				ErrorResponse.Error = err.Error()
				json.NewEncoder(res).Encode(ErrorResponse)
				return
			}
			task.Date = nextDate
		} else {
			task.Date = time.Now().Format(DateFormat)
		}
	}

	// Добавляем задачу в базу
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(res, "Не удалось добавить задачу", http.StatusInternalServerError)
		return
	}

	// Получаем идентификатор добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(res, "Не удалось вернуть id новой задачи", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ
	response := map[string]int{"id": int(id)}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

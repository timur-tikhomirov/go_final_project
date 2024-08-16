package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Обработчик следующей даты
func NextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(DateFormat, now)
	if err != nil {
		http.Error(res, "неверный формат текущей даты", http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Возвращаем ответ
	_, err = res.Write([]byte(nextDate))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

}

// Обработчик эндпоинта /task
func TaskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		taskPostHandler(res, req)
	default:
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Обработчик POST для /task
func taskPostHandler(res http.ResponseWriter, req *http.Request) {
	var task Task

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	// Проверяем обязательное поле Title
	if task.Title == "" {
		http.Error(res, "пустое поле Title", http.StatusBadRequest)
		return
	}

	// Проверяем наличие даты
	date := task.Date
	if date == "" {
		date = time.Now().Format(DateFormat)
	}

	startDate, err := time.Parse(DateFormat, date)
	if err != nil {
		http.Error(res, "неверный формат даты", http.StatusBadRequest)
		return
	}

	// Если дата меньше текущей, устанавливаем следующую дату по правилу
	if startDate.Before(time.Now()) {
		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), date, task.Repeat)
			if err != nil {
				http.Error(res, "ошибка определения даты", http.StatusBadRequest)
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
		http.Error(res, "не удалось добавить задачу", http.StatusInternalServerError)
		return
	}

	// Получаем идентификатор добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(res, "не получить id добавленной задачи", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ
	response := map[string]int{"id": int(id)}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

// Обработчик GET для /task

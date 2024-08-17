package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Обработчик GET для /tasks
func TasksGetHandler(res http.ResponseWriter, req *http.Request) {

	// Максимальное число возвращаемых задач

	// Текущая дата чтобы не отображать старые задачи
	now := time.Now().Format(DateFormat)

	// Запрос к базе данных
	var t Task
	var tasks []Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date ASC LIMIT ?`
	rows, err := DB.Query(query, now, MaxTasks)
	if err != nil {
		err := errors.New("Ошибка запроса к базе данных")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			err := errors.New("Ошибка распознавания данных")
			ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(ErrorResponse)
			return
		}
		tasks = append(tasks, t)
	}
	// Если нет задач, присвоим пустой массив
	if len(tasks) == 0 {
		tasks = []Task{}
	}

	response := TasksResponse{
		Tasks: tasks,
	}

	// Возвращаем ответ
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(response)
}

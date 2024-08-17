package main

import (
	"encoding/json"
	"errors"
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
		http.Error(res, "Некорректный формат даты", http.StatusBadRequest)
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
		http.Error(res, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик POST для /task
func taskPostHandler(res http.ResponseWriter, req *http.Request) {
	var task Task

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

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

	startDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		err := errors.New("Некорректный формат даты")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if startDate.Before(time.Now()) {
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

// Обработчик GET для /task
// Обработчик PUT для /task
// Обработчик DELETE для /task

// Обработчик GET для /tasks
func GetTasksHandler(res http.ResponseWriter, req *http.Request) {

	// Максимальное число возвращаемых задач
	const maxTasks = 50

	// Текущая дата чтобы не отображать старые задачи
	now := time.Now().Format(DateFormat)

	// Запрос к базе данных
	var t Task
	var tasks []Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date >= ? ORDER BY date ASC LIMIT ?`
	rows, err := DB.Query(query, now, maxTasks)
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

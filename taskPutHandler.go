package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Обработчик PUT для /task
func TaskPutHandler(res http.ResponseWriter, req *http.Request) {
	var task Task

	// Получаем id из запроса
	//id := req.URL.Query().Get("id")

	// Получаем тело запроса
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		err := errors.New("Ошибка десериализации JSON")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	defer req.Body.Close()

	// Проверяем, что id в теле не пустой
	if task.ID == "" {
		err := errors.New("id в теле не указан")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Проверяем, что id в query не пустой
	//if id == "" {
	//	err := errors.New("id в query не указан")
	//	ErrorResponse.Error = err.Error()
	//	json.NewEncoder(res).Encode(ErrorResponse)
	//	return
	//}

	// Проверяем обязательное поле Title
	if task.Title == "" {
		err := errors.New("Не указан заголовок задачи")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Проверяем дату
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

	// Обновляем задачу в базе
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		err := errors.New("Задача с таким id не найдена") // вот эта вот история не работает, посчитаем измененные ряды
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
		err := errors.New("Задача не изменена")
		ErrorResponse.Error = err.Error()
		json.NewEncoder(res).Encode(ErrorResponse)
		return
	}

	// Возвращаем ответ
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(nil)
}

package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/tasks"
)

// Обработчик следующей даты
func NextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(configs.DateFormat, now)
	if err != nil {
		http.Error(res, "Некорректный формат даты", http.StatusBadRequest)
		return
	}
	nextDate, err := tasks.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Возвращаем ответ
	_, err = res.Write([]byte(nextDate))
	if err != nil {
		log.Fatal(err)
		return
	}

}

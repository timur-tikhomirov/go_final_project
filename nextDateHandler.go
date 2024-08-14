package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var dateFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	//проверка на пустой repeat
	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	startDate, err := time.Parse(dateFormat, date)
	//проверка на неверный формат даты
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}
	//разделяем правило на тип правила и его аргумент
	ruleSplited := strings.Split(repeat, " ")
	//тип правила
	ruleType := ruleSplited[0]

	switch ruleType {
	case "d":
		if len(ruleSplited) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		//количество дней для переноса задачи
		daysToMove, err := strconv.Atoi(ruleSplited[1])

		if err != nil {
			return "", err
		}
		if daysToMove > 400 {

			return "", fmt.Errorf("количество дней не должно превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysToMove)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем на указанное количество дней
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysToMove)
		}
		return newDate.Format(dateFormat), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем еще на год
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(dateFormat), nil

	default:
		return "", fmt.Errorf("некорректный ввод")
	}
}

func nextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(dateFormat, now)
	if err != nil {
		return
	}
	response, err := NextDate(nowTime, date, repeat)
	io.WriteString(res, response)
}

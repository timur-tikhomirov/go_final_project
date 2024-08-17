package main

import "net/http"

// Обработчик эндпоинта /task
func TaskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		TaskPostHandler(res, req)
	case http.MethodGet:
		TaskGetHandler(res, req)
	case http.MethodPut:
		TaskPutHandler(res, req)
	case http.MethodDelete:
		TaskDeleteHandler(res, req)
	default:
		http.Error(res, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

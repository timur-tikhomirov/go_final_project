package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/storage"
)

// Обработчик DELETE для /task
func TaskDeleteHandler(store storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		err := store.DeleteTask(id)
		if err != nil {
			err := errors.New("Задача с таким id не найдена")
			configs.ErrorResponse.Error = err.Error()
			json.NewEncoder(res).Encode(configs.ErrorResponse)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

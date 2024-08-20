package handler

import (
	"encoding/json"
	"net/http"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/storage"
)

// Обработчик PUT для /task
func TaskPutHandler(store storage.Store) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var t configs.Task
		err := json.NewDecoder(req.Body).Decode(&t)
		if err != nil {
			http.Error(res, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}
		err = store.PutTask(t)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{}); err != nil {
			http.Error(res, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/timur-tikhomirov/go_final_project/configs"
)

// Получение листинга задач по фильтрам
func (s *Store) GetTasks(search string) ([]configs.Task, error) {
	var t configs.Task
	var tasks []configs.Task
	var rows *sql.Rows
	var err error
	if search == "" {
		rows, err = s.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", configs.MaxTasks)
	} else if date, error := time.Parse("02.01.2006", search); error == nil {
		query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, date.Format(configs.DateFormat), configs.MaxTasks)
	} else {
		search = "%%%" + search + "%%%"
		query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, search, search, configs.MaxTasks)
	}
	if err != nil {
		return []configs.Task{}, fmt.Errorf(`{"error":"ошибка запроса"}`)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []configs.Task{}, fmt.Errorf(`{"error":"Ошибка распознавания данных"}`)
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		tasks = []configs.Task{}
	}

	return tasks, nil
}

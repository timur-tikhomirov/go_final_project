package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func main() {
	// Определяем директорию приложения и проверяем наличие базы данных
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	log.Println(dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	DB, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	// создаем таблицу и индекс
	if install {

		_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128) NOT NULL
			);`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("База данных создана")
	} else {
		log.Println("База данных была создана ранее")
	}

	// Определяем порт из окружения, если переменная окружения отсутствует - устанавливаем порт по умолчанию
	port := DefaultPort
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) != 0 {
		port = envPort
	}
	port = ":" + port

	// Обрабатываем запрос
	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", GetTasksHandler)

	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}

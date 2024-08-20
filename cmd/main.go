package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/timur-tikhomirov/go_final_project/configs"
	"github.com/timur-tikhomirov/go_final_project/internal/handler"
	"github.com/timur-tikhomirov/go_final_project/internal/storage"

	_ "modernc.org/sqlite"
)

func main() {
	// Открываем/создаем базу данных
	dataBase := storage.OpenDataBase()
	defer dataBase.Close()
	store := storage.NewStore(dataBase)

	// Определяем порт из окружения, если переменная окружения отсутствует - устанавливаем порт по умолчанию
	port := configs.DefaultPort
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) != 0 {
		port = envPort
	}
	port = ":" + port
	// Определяем директорию приложения
	filePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	//Определяем корневую папку проекта
	projectPath := strings.TrimRight(filePath, "/cmd")
	log.Println("Директория проекта", projectPath)

	// Создаем хендлер для файлов фронта
	fileServer := http.FileServer(http.Dir(projectPath + configs.WebDir))
	// Обрабатываем запросы
	http.Handle("/", fileServer)
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)
	http.HandleFunc("GET /api/task", handler.TaskGetHandler(store))
	http.HandleFunc("POST /api/task", handler.TaskPostHandler(store))
	http.HandleFunc("PUT /api/task", handler.TaskPutHandler(store))
	http.HandleFunc("DELETE /api/task", handler.TaskDeleteHandler(store))
	http.HandleFunc("/api/tasks", handler.TasksGetHandler(store))
	http.HandleFunc("/api/task/done", handler.TaskDoneHandler(store))

	log.Println("Приложение запущено на порту", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}

}

package main

import (
	"net/http"
)

func main() {
	// Определяем порт из окружения, если переменная окружения отсутствует - устанавливаем порт по умолчанию
	//port := "7540"
	//envPort := os.Getenv("TODO_PORT")
	//if len(envPort) != 0 {
	//	port = envPort
	//}
	//port = ":" + port

	// Обрабатываем запрос
	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}

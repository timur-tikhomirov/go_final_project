package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

type Store struct {
	db *sql.DB
}

// Открываем/создаем  БД
func OpenDataBase() *sql.DB {
	// Определяем директорию приложения и проверяем наличие базы данных
	filePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(filePath), "scheduler.db")

	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}
	log.Println("Путь к базе данных", dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	// создаем таблицу и индекс
	if install {

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128) NOT NULL
			);`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("База данных создана")
	} else {
		log.Println("База данных была создана ранее")
	}
	return db

}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

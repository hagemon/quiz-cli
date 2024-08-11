package db

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Metadata struct {
	ID           int    `db:"id"`
	QuizName     string `db:"quiz_name"`
	NumOfQuizzes int    `db:"num_of_quizzes"`
	Loc          int    `db:"loc"`
	ErrNum       int    `db:"err_num"`
}

type Quiz struct {
	ID         int    `db:"id"`
	QuizID     int    `db:"quiz_id"`
	Question   string `db:"question"`
	Answer     string `db:"answer"`
	Analysis   string `db:"analysis"`
	Difficulty string `db:"difficulty"`
	Options    string `db:"options"`
	OptionsNum int    `db:"options_num"`
	QuizType   string `db:"quiz_type"`
	IsError    bool   `db:"is_error"`
	IsVisited  bool   `db:"is_visited"`
}

func InitDB() error {

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS metadata (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		quiz_name TEXT NOT NULL,
		num_of_quizzes INT NOT NULL,
		loc INT NOT NULL DEFAULT 0,
		err_num INT NOT NULL DEFAULT 0
	);`

	createQuizTableSQL := `CREATE TABLE IF NOT EXISTS quiz (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		quiz_id INTEGER NOT NULL,
		question TEXT NOT NULL,
		answer TEXT NOT NULL,
		analysis TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		options TEXT NOT NULL,
		options_num INT NOT NULL,
		quiz_type TEXT NOT NULL,
		is_error BOOLEAN NOT NULL DEFAULT FALSE,
		is_visited BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	_, err = db.Exec(createQuizTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database and table initialized successfully.")
	return nil
}

func GetDBPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dbpath := filepath.Join(usr.HomeDir, ".quiz", "quiz.db")
	dbDir := filepath.Dir(dbpath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			return "", err
		} else {
			log.Println("Database created at:", dbDir)
		}
	}
	return dbpath, nil

}

func OpenDB() (*sql.DB, error) {
	dbpath, err := GetDBPath()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

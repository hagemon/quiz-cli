package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
)

func QueryMetadata() ([]Metadata, error) {
	var metadata []Metadata
	results, err := query("SELECT * FROM metadata", &metadata)
	if err != nil {
		return nil, err
	}

	metadata = make([]Metadata, len(results))
	for i, result := range results {
		metadata[i] = result.(Metadata)
	}

	return metadata, nil
}

func query(query string, target interface{}) ([]interface{}, error) {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice {
		return nil, fmt.Errorf("target must be a pointer to a slice")
	}

	// 获取切片元素的类型
	elemType := targetVal.Elem().Type().Elem()
	elemVal := reflect.New(elemType).Elem()

	// 准备查询
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 创建一个扫描目标
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = elemVal.Field(i).Addr().Interface()
	}

	// 扫描结果并填充到结构体中
	var results []interface{}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		results = append(results, elemVal.Interface())
	}

	return results, nil
}

func CheckQuizID(quizID string) (bool, error) {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()
	var exists bool
	err = db.QueryRow(`SELECT COUNT(*) FROM metadata WHERE id = ?`, quizID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetQuizInfo(quizID string) (Metadata, error) {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		log.Fatal(err)
		return Metadata{}, err
	}
	defer db.Close()
	var metadata Metadata
	err = db.QueryRow(`SELECT * FROM metadata WHERE id = ?`, quizID).Scan(&metadata.ID, &metadata.QuizName, &metadata.NumOfQuizzes, &metadata.Loc, &metadata.ErrNum)
	if err == sql.ErrNoRows {
		return Metadata{}, fmt.Errorf("no rows found for quiz ID %s", quizID)
	}
	if err != nil {
		return Metadata{}, err
	}
	return metadata, nil
}

func GetAllQuiz(quizID int) ([]Quiz, error) {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var quizzes []Quiz
	rows, err := db.Query(`SELECT * FROM quiz WHERE quiz_id = ?`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var quiz Quiz
		err = rows.Scan(&quiz.ID, &quiz.QuizID, &quiz.Question, &quiz.Answer, &quiz.Analysis, &quiz.Difficulty, &quiz.Options, &quiz.OptionsNum, &quiz.QuizType, &quiz.IsError, &quiz.IsVisited)
		if err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func GetRandomQuiz(quizID int, num int) ([]Quiz, error) {
	quiz, err := GetAllQuiz(quizID)
	if err != nil {
		return nil, err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(quiz), func(i, j int) {
		quiz[i], quiz[j] = quiz[j], quiz[i]
	})
	return quiz[:num], nil
}

func GetErrorQuiz(quizID int) ([]Quiz, error) {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var quizzes []Quiz
	rows, err := db.Query(`SELECT * FROM quiz WHERE quiz_id = ? AND is_error = TRUE`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var quiz Quiz
		err = rows.Scan(&quiz.ID, &quiz.QuizID, &quiz.Question, &quiz.Answer, &quiz.Analysis, &quiz.Difficulty, &quiz.Options, &quiz.OptionsNum, &quiz.QuizType, &quiz.IsError, &quiz.IsVisited)
		if err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

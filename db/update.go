package db

import "database/sql"

func UpdateQuizError(quiz Quiz) error {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE quiz SET is_error = TRUE WHERE id = :id", quiz.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateQuizVisit(quiz Quiz) error {
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE quiz SET is_visited = TRUE WHERE id = :id", quiz.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE metadata SET loc = GREATEST(loc, :id) WHERE id = :quiz_id", quiz.ID, quiz.QuizID)
	if err != nil {
		return err
	}
	return nil
}

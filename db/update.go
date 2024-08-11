package db

import (
	"database/sql"
	"log"
)

func UpdateQuizError(quiz Quiz) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE quiz SET is_error = TRUE WHERE id = :id", quiz.ID)

	if err != nil {
		log.Printf("Update Quiz Error Status: %v", err)
		return err
	}
	_, err = db.Exec(`UPDATE metadata
	SET err_num = (
		SELECT COUNT(*)
		FROM quiz 
		WHERE quiz_id = :id and is_error = TRUE
	)
	WHERE id = :id`,
		sql.Named("id", quiz.QuizID))
	if err != nil {
		log.Printf("Update Quiz Error Metadata Status: %v", err)
		return err
	}
	return nil
}

func UpdateQuizVisit(quiz Quiz) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE quiz SET is_visited = TRUE WHERE id = :id", quiz.ID)
	if err != nil {
		log.Printf("Update Quiz Visit Status: %v", err)
		return err
	}
	_, err = db.Exec("UPDATE metadata SET loc = CASE WHEN loc > :id THEN loc ELSE :id END WHERE id = :quiz_id",
		sql.Named("id", quiz.ID),
		sql.Named("quiz_id", quiz.QuizID))
	if err != nil {
		log.Printf("Update Quiz Visit Metadata Status: %v", err)
		return err
	}
	return nil
}

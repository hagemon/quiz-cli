package db

import "log"

func UpdateQuizError(quiz Quiz) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE quiz SET is_error = TRUE WHERE id = :id", quiz.ID)

	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE metadata
	SET err_num = (
		SELECT COUNT(*)
		FROM quiz 
		WHERE quiz_id = :quiz_id and is_error = TRUE
	)
	WHERE id = :id`, quiz.ID, quiz.QuizID)
	if err != nil {
		log.Println(err)
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
		return err
	}
	_, err = db.Exec("UPDATE metadata SET loc = CASE WHEN loc > :id THEN loc ELSE :id END WHERE id = :quiz_id", quiz.ID, quiz.ID, quiz.QuizID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

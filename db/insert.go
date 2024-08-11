package db

import (
	"log"
)

func InsertQuizzes(quizzes []Quiz) error {

	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()
	for _, quiz := range quizzes {
		_, err = db.Exec(`INSERT INTO quiz (question, answer, analysis, difficulty, options, options_num, quiz_type, quiz_id) VALUES (:question, :answer, :analysis, :difficulty, :options, :options_num, :quiz_type, :quiz_id)`, quiz.Question, quiz.Answer, quiz.Analysis, quiz.Difficulty, quiz.Options, quiz.OptionsNum, quiz.QuizType, quiz.QuizID)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func InsertMetadata(metadata Metadata) (int, error) {

	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer db.Close()
	result, err := db.Exec(`INSERT INTO metadata (quiz_name, num_of_quizzes) VALUES (:quiz_name, :num_of_quizzes)`,
		metadata.QuizName, metadata.NumOfQuizzes)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

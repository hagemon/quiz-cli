package main

import (
	"log"
	"os"
	"quiz-cli/db"
	"quiz-cli/pages"

	"github.com/rivo/tview"
)

func main() {
	logFile, err := os.OpenFile("quiz.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	err = db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		os.Exit(1)
	}

	app := tview.NewApplication()
	layout := pages.NewHomePage(app)
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}

}

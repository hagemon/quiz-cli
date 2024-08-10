package pages

import (
	"fmt"
	"quiz-cli/db"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func handleStartCmd(pages *tview.Pages, quizID string) error {
	// Check if the quizID exists in the database
	metadata, err := db.GetQuizInfo(quizID)
	if err != nil {
		return err
	}

	// Create a new quiz page
	quizModePage, err := newQuizModePage(pages, metadata)
	if err != nil {
		return err
	}
	pages.AddAndSwitchToPage("quizMode", quizModePage, true)
	return nil
}

func newQuizModePage(pages *tview.Pages, metadata db.Metadata) (tview.Primitive, error) {
	quizModeLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	quizModeTitle := tview.NewTextView().SetText("Select Quiz Mode").SetTextAlign(tview.AlignCenter)
	quizMode := tview.NewList()
	quizMode.AddItem("1. All Quizs", "", 0, nil).
		AddItem("2. Wrong Quizs", "", 0, nil).
		AddItem("3. Random 100 Quizs", "", 0, nil)
	prompt := tview.NewTextView().SetText("Press ESC to exit").SetTextAlign(tview.AlignCenter)
	quizMode.SetSelectedFunc(func(index int, mainText, secondaryText string, key rune) {
		var quiz []db.Quiz
		var err error
		startIndex := 0
		prompt.SetText("Getting quiz data...")
		switch index {
		case 0:
			quiz, err = db.GetAllQuiz(metadata.ID)
			startIndex = 1500
		case 1:
			quiz, err = db.GetErrorQuiz(metadata.ID)
		case 2:
			quiz, err = db.GetRandomQuiz(metadata.ID, 100)
		}
		if err != nil {
			prompt.SetText(err.Error() + " Press ESC to exit")
			return
		}
		if len(quiz) == 0 {
			prompt.SetText("No quiz found" + " Press ESC to exit")
			return
		}
		quizPage := newQuizPage(pages, quiz, startIndex)
		pages.AddAndSwitchToPage("quiz", quizPage, true)
	})
	quizModeLayout.AddItem(quizModeTitle, 2, 0, false).
		AddItem(quizMode, 0, 1, true).
		AddItem(prompt, 1, 0, false)
	pageLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(quizModeLayout, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)

	pageLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage("home")
		}
		return event
	})
	return pageLayout, nil
}

func newQuizPage(pages *tview.Pages, quiz []db.Quiz, startIndex int) tview.Primitive {
	quizString, optionsString, _ := formatQuiz(quiz[startIndex])
	num_options := len(strings.Split(optionsString, "\n"))
	quizContentView := tview.NewTextView().SetText(quizString).SetWrap(true)
	optionsContentView := tview.NewTextView().SetText(optionsString).SetWrap(true).SetRegions(true).SetDynamicColors(true)
	answerContentView := tview.NewTextView().SetText("").SetWrap(true)
	contentview := tview.NewFlex().SetDirection(tview.FlexRow)
	contentview.AddItem(tview.NewBox(), 1, 0, false).
		AddItem(quizContentView, 0, 1, true).
		AddItem(optionsContentView, 0, 1, true).
		AddItem(answerContentView, 0, 1, true).
		AddItem(tview.NewBox(), 0, 0, false)

	quizLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	quizLayout.AddItem(tview.NewBox(), 10, 1, false).
		AddItem(contentview, 0, 1, true).
		AddItem(tview.NewBox(), 10, 1, false)

	selectedOption := []string{}
	showAnswer := false

	quizLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEsc:
			pages.SwitchToPage("home")
		case event.Rune() == ' ':
			if !showAnswer {
				if strings.Join(selectedOption, "") == quiz[startIndex].Answer {
					answerContentView.SetText("Correct")
				} else {
					answerContentView.SetText("Answer: " + quiz[startIndex].Answer)
				}
				showAnswer = true
			} else {
				optionsContentView.Highlight("")
				answerContentView.SetText("")
				selectedOption = []string{}
				startIndex++
				if startIndex < len(quiz) {
					quizString, optionsString, _ := formatQuiz(quiz[startIndex])
					quizContentView.SetText(quizString)
					optionsContentView.SetText(optionsString)
				}
				showAnswer = false
			}

		case event.Rune() == 'E' || event.Rune() == 'e':
			showAnswer = true
		case event.Rune() >= '1' && event.Rune() <= '9':
			if event.Rune() > rune(num_options+'0') {
				return event
			}
			option := string('A' + event.Rune() - '1')
			if quiz[startIndex].QuizType == "多选题" {
				if contains(selectedOption, option) {
					selectedOption = remove(selectedOption, option)

				} else {
					selectedOption = append(selectedOption, option)
				}
			} else {
				if len(selectedOption) == 0 {
					selectedOption = append(selectedOption, option)
				} else {
					selectedOption[0] = option
				}
			}
			optionsContentView.Highlight(selectedOption...)

		}

		return event
	})
	return quizLayout
}

func formatQuiz(quiz db.Quiz) (string, string, string) {
	options := strings.Split(quiz.Options, "$$")
	for i, option := range options {
		if option != "" {
			options[i] = fmt.Sprintf("[\"%c\"]%c. %s[\"green\"]", 'A'+i, 'A'+i, option)
		}
	}
	optionsString := strings.Join(options, "\n")
	quizString := fmt.Sprintf("%s\n%s", quiz.QuizType, quiz.Question)
	answerString := fmt.Sprintf("Answer: %s", quiz.Answer)
	return quizString, optionsString, answerString
}

func contains(selectedOption []string, option string) bool {
	for _, opt := range selectedOption {
		if opt == option {
			return true
		}
	}
	return false
}
func remove(selectedOption []string, option string) []string {
	for i, opt := range selectedOption {
		if opt == option {
			return append(selectedOption[:i], selectedOption[i+1:]...)
		}
	}
	return selectedOption
}

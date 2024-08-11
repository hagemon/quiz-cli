package pages

import (
	"fmt"
	"quiz-cli/db"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tealeg/xlsx"
)

func handleImport(app *tview.Application, pages *tview.Pages, filename string, homeLayout *tview.Flex, messageField *tview.TextView) error {
	// 解析Excel文件
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	// todo: 页面上展示一个表格，每一行展示xlFile的一个列信息，包括序号、列名和第一个案例数据
	// 创建一个新的表格
	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1) // 固定表头

	// 设置表头
	table.SetCell(0, 0, tview.NewTableCell("No.").SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("Column").SetSelectable(false).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("Example").SetSelectable(false).SetAlign(tview.AlignCenter))

	// 获取第一个Sheet
	if len(xlFile.Sheets) == 0 {
		return fmt.Errorf("the Excel file does not contain any sheets")
	}
	sheet := xlFile.Sheets[0]

	// 遍历列并设置表格内容
	for colIndex, cell := range sheet.Rows[0].Cells {
		columnName := cell.String()
		var exampleData string
		if len(sheet.Rows) > 1 && len(sheet.Rows[1].Cells) > colIndex {
			exampleData = sheet.Rows[1].Cells[colIndex].String()
		} else {
			exampleData = "N/A"
		}
		if utf8.RuneCountInString(exampleData) > 10 {
			runes := []rune(exampleData)
			exampleData = string(runes[:10]) + "..."
		}
		table.SetCell(colIndex+1, 0, tview.NewTableCell(fmt.Sprintf("%d", colIndex+1)).SetAlign(tview.AlignCenter))
		table.SetCell(colIndex+1, 1, tview.NewTableCell(columnName).SetAlign(tview.AlignCenter))
		table.SetCell(colIndex+1, 2, tview.NewTableCell(exampleData).SetAlign(tview.AlignCenter))
	}

	questions := []string{
		"1. Question: ",
		"2. Quiz Type: ",
		"3. Options (e.g. 3-6): ",
		"4. Answer: ",
		"5. Analysis (Optional): ",
		"6. Difficulty (Optional): ",
	}

	// 交互式输入
	labelView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetText(fmt.Sprintf("Configures:\n\nPlease input the corresponding number on the left\n%s", questions[0]))

	inputField := tview.NewInputField().
		SetFieldWidth(20)

	// 添加提示信息
	prompt := tview.NewTextView().
		SetText("Press 'ESC' to cancel").
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignCenter)

	tableAndInputLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(table, 0, 3, true).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(labelView, 0, 1, false).
				AddItem(inputField, 0, 1, false),
			0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)

	// 创建一个新的布局，将表格和提示信息放在一起
	tableLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 1, 0, false). // 上边距
		AddItem(tableAndInputLayout, 0, 1, true).
		AddItem(prompt, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false) // 下边距

	// 将表格居中显示的容器
	importLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 1, 1, false). // 左边距
		AddItem(tableLayout, 0, 1, true).
		AddItem(tview.NewBox(), 1, 1, false) // 右边距

	// 设置新的布局为应用的根布局，并监听键盘事件
	pages.AddAndSwitchToPage("import", importLayout, true)
	app.SetFocus(inputField)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.SwitchToPage("home")
			app.SetFocus(homeLayout)
		}
		return event
	})

	answers := make([]string, len(questions))
	currentQuestion := 0

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			currentAnswer := inputField.GetText()
			answers[currentQuestion] = currentAnswer
			currentQuestion++
			nextQuestion := ""
			if currentQuestion < len(questions) {
				nextQuestion = questions[currentQuestion]
			}
			label := fmt.Sprintf("%s%s\n%s", labelView.GetText(false), currentAnswer, nextQuestion)
			labelView.SetText(label)
			inputField.SetText("")

			if currentQuestion >= len(questions) {
				app.SetFocus(importLayout)
				inputField.SetDisabled(true)
				prompt.SetText("Press 'y' to confirm import or 'ESC' to cancel")
			}

		}
	})

	importLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.SwitchToPage("home")
			app.SetFocus(homeLayout)
		case tcell.KeyRune:
			if event.Rune() == 'y' && currentQuestion >= len(questions) {
				prompt.SetText("Importing Quizs...")
				err := handleImportConfirm(strings.TrimSuffix(filename, ".xlsx"), sheet, answers)
				if err != nil {
					messageField.SetText(err.Error())
				} else {
					messageField.SetText(fmt.Sprintf("Import %s successful", filename))
				}
				pages.SwitchToPage("home")
				app.SetFocus(homeLayout)
			}
		}
		return event
	})

	return nil
}

func handleImportConfirm(quiz_name string, sheet *xlsx.Sheet, answers []string) error {
	metadata_id, err := insertMetadata(quiz_name, len(sheet.Rows)-1)
	if err != nil {
		return err
	}
	err = insertQuizzes(metadata_id, sheet, answers)
	if err != nil {
		return err
	}
	return nil
}

func insertMetadata(quiz_name string, num_of_quizzes int) (int, error) {
	metadata := db.Metadata{
		QuizName:     quiz_name,
		Loc:          0,
		NumOfQuizzes: num_of_quizzes,
		ErrNum:       0,
	}
	id, err := db.InsertMetadata(metadata)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertQuizzes(metadata_id int, sheet *xlsx.Sheet, answers []string) error {
	quizzes := []db.Quiz{}

	questionIndex, err := strconv.Atoi(answers[0])
	if err != nil {
		return err
	}
	quizTypeIndex, err := strconv.Atoi(answers[1])
	if err != nil {
		return err
	}
	answerIndex, err := strconv.Atoi(answers[3])
	if err != nil {
		return err
	}
	var analysisIndex, difficultyIndex int
	if answers[4] == "" {
		analysisIndex = -1
	} else {
		analysisIndex, err = strconv.Atoi(answers[4])
		if err != nil {
			return err
		}
	}
	if answers[5] == "" {
		difficultyIndex = -1
	} else {
		difficultyIndex, err = strconv.Atoi(answers[5])
		if err != nil {
			return err
		}
	}

	start, end := 0, 0
	if strings.Contains(answers[2], "-") {
		parts := strings.Split(answers[2], "-")
		start, err = strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		end, err = strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
	} else {
		start, err = strconv.Atoi(answers[2])
		if err != nil {
			return err
		}
		end = start
	}

	for i, row := range sheet.Rows {
		if i == 0 { // Skip the first row as it's the header
			continue
		}
		cells := row.Cells
		var options []string
		for i := start; i < end; i++ {
			optionString := cells[i-1].String()
			if optionString != "" {
				options = append(options, optionString)
			}
		}
		analysis := ""
		if analysisIndex != -1 {
			analysis = cells[analysisIndex-1].String()
		}
		difficulty := ""
		if difficultyIndex != -1 {
			difficulty = cells[difficultyIndex-1].String()
		}
		quiz := db.Quiz{
			QuizID:     metadata_id,
			Question:   cells[questionIndex-1].String(),
			Answer:     cells[answerIndex-1].String(),
			Analysis:   analysis,
			Difficulty: difficulty,
			Options:    strings.Join(options, "$$"),
			OptionsNum: len(options),
			QuizType:   cells[quizTypeIndex-1].String(),
		}
		quizzes = append(quizzes, quiz)
	}
	err = db.InsertQuizzes(quizzes)
	if err != nil {
		return err
	}
	return nil
}

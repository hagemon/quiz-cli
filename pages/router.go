package pages

import (
	"fmt"
	"strings"

	"quiz-cli/db"

	"github.com/rivo/tview"
)

func handleCommand(app *tview.Application, pages *tview.Pages, command string, homeLayout *tview.Flex, messageField *tview.TextView, metaTable *tview.Table) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	homeLayout.RemoveItem(messageField)
	homeLayout.RemoveItem(metaTable)

	switch parts[0] {
	case "show":
		metadatas, err := db.QueryMetadata()
		if err != nil {
			messageField.SetText(fmt.Sprintf("%v", err))
			homeLayout.AddItem(messageField, 0, 1, true)
		} else {
			metaTable.SetCell(0, 0, tview.NewTableCell("ID").SetExpansion(1).SetAlign(tview.AlignCenter))
			metaTable.SetCell(0, 1, tview.NewTableCell("Quiz Name").SetExpansion(1).SetAlign(tview.AlignCenter))
			metaTable.SetCell(0, 2, tview.NewTableCell("Number of Quizzes").SetExpansion(1).SetAlign(tview.AlignCenter))
			metaTable.SetCell(0, 3, tview.NewTableCell("Location").SetExpansion(1).SetAlign(tview.AlignCenter))
			metaTable.SetCell(0, 4, tview.NewTableCell("Error Number").SetExpansion(1).SetAlign(tview.AlignCenter))

			for i, metadata := range metadatas {
				metaTable.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", metadata.ID)).SetAlign(tview.AlignCenter))
				metaTable.SetCell(i+1, 1, tview.NewTableCell(metadata.QuizName).SetAlign(tview.AlignCenter))
				metaTable.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("%d", metadata.NumOfQuizzes)).SetAlign(tview.AlignCenter))
				metaTable.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprintf("%d", metadata.Loc)).SetAlign(tview.AlignCenter))
				metaTable.SetCell(i+1, 4, tview.NewTableCell(fmt.Sprintf("%d", metadata.ErrNum)).SetAlign(tview.AlignCenter))
			}
			homeLayout.AddItem(metaTable, 0, 1, true)
		}
		// 跳转到
	case "import":
		if len(parts) < 2 {
			return "Error: no file provided."
		}
		filename := parts[1]
		err := handleImport(app, pages, filename, homeLayout, messageField)
		if err != nil {
			messageField.SetText(fmt.Sprintf("%v", err))
			homeLayout.AddItem(messageField, 0, 1, true)
		}
	case "q":
		app.Stop()
	case "start":
		if len(parts) < 2 {
			return "Error: no quiz ID provided."
		}
		quizID := parts[1]
		err := handleStartCmd(pages, quizID)
		if err != nil {
			messageField.SetText(fmt.Sprintf("%v", err))
			homeLayout.AddItem(messageField, 0, 1, true)
		}

	default:
		messageField.SetText("Unknown command: " + parts[0])
		homeLayout.AddItem(messageField, 0, 1, true)
	}
	return ""
}

package pages

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewHomePage(app *tview.Application) tview.Primitive {

	textView, linesCount := newWelcomeView()
	inputField := newInputField()

	messageField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true)

	metaTable := tview.NewTable().SetBorders(true)

	// composite layout
	homeLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, linesCount, 1, true)

	// homeLayout := tview.NewFlex().SetDirection(tview.FlexRow).
	// 	AddItem(innerLayout, 0, 1, true)

	pages := tview.NewPages().AddPage("home", homeLayout, true, true)

	// add action

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := inputField.GetText()
			inputField.SetText("")
			homeLayout.RemoveItem(inputField)
			app.SetFocus(textView)
			handleCommand(app, pages, command, homeLayout, messageField, metaTable)
		} else if key == tcell.KeyEscape {
			inputField.SetText("")
			homeLayout.RemoveItem(inputField)
			app.SetFocus(textView)
		}
	})

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == ':' {
			homeLayout.AddItem(inputField, 0, 1, true)
			app.SetFocus(inputField)
		}
		return event
	})

	return pages

}

func newWelcomeView() (*tview.TextView, int) {
	welcome_info := `Welcome to QUIZ system.
	
	Press ':' to enter command mode.

		show               show quiz in database
		import filename    import a xlsx file to database
        start  id          start quizing with id
		q                  exit

	Press 'CTRL+C' to exit program.
	`

	// views
	textView := tview.NewTextView().
		SetText(welcome_info).
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true)

	linesCount := strings.Count(welcome_info, "\n") + 1

	return textView, linesCount
}

func newInputField() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(":").
		SetFieldWidth(20)

	return inputField
}

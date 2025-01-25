package uiprogress

import (
	"github.com/rivo/tview"
)

func ProjectLogDisplay(log string) *tview.Flex {
	box := tview.NewFlex()
	box.Box = tview.NewBox()
	box.SetDirection(tview.FlexRow)
	box.SetBorder(true)
	box.SetTitle("Log")
	box.SetTitleColor(tview.Styles.SecondaryTextColor)

	box.AddItem(tview.NewTextView().SetText("⮕ "+log), 0, 1, false)
	return box
}

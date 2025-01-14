package uiprogress

import (
	"github.com/rivo/tview"
	"goFlaky/core/progress"
)

func ProjectDisplay(progress []progress.ProjectProgress) *tview.Flex {
	box := tview.NewFlex()
	box.Box = tview.NewBox()
	box.SetDirection(tview.FlexRow)
	box.SetBorder(true)
	box.SetTitle("Projects")
	box.SetTitleColor(tview.Styles.SecondaryTextColor)
	for _, v := range progress {
		box.AddItem(tview.NewTextView().SetText(v.Identifier), 1, 0, false)
	}

	return box
}

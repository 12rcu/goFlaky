package uiprogress

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"goFlaky/core/progress"
)

func ProjectProgressDisplay(progress []progress.ProjectProgress) *tview.Flex {
	box := tview.NewFlex()
	box.Box = tview.NewBox()
	box.SetDirection(tview.FlexRow)
	box.SetBorder(true)
	box.SetTitle("Progress")
	box.SetTitleColor(tview.Styles.SecondaryTextColor)
	for _, v := range progress {
		box.AddItem(createProgressBar(v.Index, v.Runs), 1, 0, false)
	}

	return box
}

func createProgressBar(index int, runs int) *tview.Flex {
	box := tview.NewFlex()
	box.SetDirection(tview.FlexColumn)
	progressDone := tview.
		NewTextView().
		SetText("").
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.NewHexColor(0xb294bb))
	progressRemain := tview.
		NewTextView().
		SetText("").
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.NewHexColor(0x85678f))

	box.AddItem(progressDone, 0, index, false)
	box.AddItem(progressRemain, 0, runs-index, false)

	return box
}

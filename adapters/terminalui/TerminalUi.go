package terminalui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"goFlaky/adapters/terminalui/uiprogress"
	"goFlaky/core/progress"
	"sync"
	"time"
)

func TerminalUi(
	initProgress []progress.ProjectProgress,
	progressChannel chan []progress.ProjectProgress,
	logChannel chan string,
	waitGroup *sync.WaitGroup,
) {
	app := tview.NewApplication()
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewHexColor(0x1d1f21),
		ContrastBackgroundColor:     tcell.NewHexColor(0x282a2e),
		MoreContrastBackgroundColor: tcell.NewHexColor(0x373b41),
		BorderColor:                 tcell.NewHexColor(0x707880),
		TitleColor:                  tcell.NewHexColor(0xf0c674),
		GraphicsColor:               tcell.NewHexColor(0xf0c674),
		PrimaryTextColor:            tcell.NewHexColor(0xc5c8c6),
		SecondaryTextColor:          tcell.NewHexColor(0x81a2be),
		TertiaryTextColor:           tcell.NewHexColor(0xb5bd68),
		InverseTextColor:            tcell.NewHexColor(0x5f819d),
		ContrastSecondaryTextColor:  tcell.NewHexColor(0x85678f),
	}

	box := tview.NewFlex()
	box.SetDirection(tview.FlexColumn)
	box.SetBorder(true)
	box.SetTitle("[::b]GoFlaky")

	wrapper := progressWrapper{
		projectProgress: initProgress,
		projectLogs:     "",
		done:            false,
	}

	updateUi(box, wrapper.projectProgress, wrapper.projectLogs)

	go drawApp(app, box)

	go receiveProgress(progressChannel, &wrapper)
	go receiveLogs(logChannel, &wrapper)

	for !wrapper.done {
		time.Sleep(1 * time.Second)
		app.QueueUpdateDraw(func() {
			updateUi(box, wrapper.projectProgress, wrapper.projectLogs)
		})
	}
	time.Sleep(10 * time.Second)
	app.Stop()
	waitGroup.Done()
}

func drawApp(app *tview.Application, box *tview.Flex) {
	if err := app.SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

func receiveProgress(progressChannel chan []progress.ProjectProgress, wrapper *progressWrapper) {
	for prg := range progressChannel {
		wrapper.projectProgress = prg
	}
	wrapper.done = true
}

func receiveLogs(logChannel chan string, wrapper *progressWrapper) {
	for prg := range logChannel {
		log := prg + "\n" + wrapper.projectLogs

		wrapper.projectLogs = log
	}
}

func updateUi(box *tview.Flex, currProgress []progress.ProjectProgress, currLogs string) {
	box.Clear()
	box.AddItem(uiprogress.ProjectDisplay(currProgress), 0, 3, false)
	box.AddItem(uiprogress.ProjectProgressDisplay(currProgress), 0, 3, false)
	box.AddItem(uiprogress.ProjectLogDisplay(currLogs), 0, 10, false)
}

type progressWrapper struct {
	projectProgress []progress.ProjectProgress
	projectLogs     string
	done            bool
}

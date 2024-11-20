package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})))
}
func main() {
	a := app.New()

	w := a.NewWindow("vsplit bug")

	largeContainer := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("test1"), widget.NewLabel("just a test 1"),
		widget.NewLabel("test2"), widget.NewLabel("just a test 2"),
		widget.NewLabel("test3"), widget.NewLabel("just a test 3"),
		widget.NewLabel("test4"), widget.NewLabel("just a test 4"),
		widget.NewLabel("test5"), widget.NewLabel("just a test 5"),
		widget.NewLabel("test6"), widget.NewLabel("just a test 6"),
		widget.NewLabel("test7"), widget.NewLabel("just a test 7"),
		widget.NewLabel("test8"), widget.NewLabel("just a test 8"),
		widget.NewLabel("test9"), widget.NewLabel("just a test 9"),
		widget.NewLabel("test10"), widget.NewLabel("just a test 10"),
		widget.NewLabel("test11"), widget.NewLabel("just a test 11"),
		widget.NewLabel("test12"), widget.NewLabel("just a test 12"),
		widget.NewLabel("test13"), widget.NewLabel("just a test 13"),
		widget.NewLabel("test14"), widget.NewLabel("just a test 14"),
		widget.NewLabel("test15"), widget.NewLabel("just a test 15"),
		widget.NewLabel("test16"), widget.NewLabel("just a test 16"),
		widget.NewLabel("test17"), widget.NewLabel("just a test 17"),
		widget.NewLabel("test18"), widget.NewLabel("just a test 18"),
		widget.NewLabel("test19"), widget.NewLabel("just a test 19"),
		widget.NewLabel("test20"), widget.NewLabel("just a test 20"),
	)

	l := widget.NewLabel("not clicked")
	w.Resize(fyne.NewSize(800, 600))
	b := widget.NewButton("test just a button", func() {
		l.SetText("Clicked")
	})

	b.Icon = theme.IconForWidget(theme.IconNameRadioButton, b)
	b.Importance = widget.MediumImportance
	b.Alignment = widget.ButtonAlignLeading
	b.IconPlacement = widget.ButtonIconLeadingText
	b.Refresh()
	w.SetContent(container.NewVSplit(largeContainer, container.NewVBox(b, l)))

	w.ShowAndRun()
}

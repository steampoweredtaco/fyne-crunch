package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/steampoweredtaco/fyne-crunch/widgets/pinboard"
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

	w := a.NewWindow("pinboard demo")

	// The pinboard return is used to interact with the pinboard programmatically and the
	// scroll is the canvas object that should be used for views. This will change in the
	// future to only return a single item that does both which would be more idiomatic
	pb, scroll := pinboard.NewPinBoard(
		pinboard.NewPinBoardItem("test1", widget.NewLabel("just a test 1")),
		pinboard.NewPinBoardItem("test2", widget.NewLabel("just a test 2")),
		pinboard.NewPinBoardItem("test3", widget.NewLabel("just a test 3")),
		pinboard.NewPinBoardItem("test4", widget.NewLabel("just a test 4")),
		pinboard.NewPinBoardItem("test5", widget.NewLabel("just a test 5")),
		pinboard.NewPinBoardItem("test6", widget.NewLabel("just a test 6")),
		pinboard.NewPinBoardItem("test7", widget.NewLabel("just a test 7")),
		pinboard.NewPinBoardItem("test8", widget.NewLabel("just a test 8")),
		pinboard.NewPinBoardItem("test9", widget.NewLabel("just a test 9")),
		pinboard.NewPinBoardItem("test10", widget.NewLabel("just a test 10")),
		pinboard.NewPinBoardItem("test11", widget.NewLabel("just a test 11")),
		pinboard.NewPinBoardItem("test12", widget.NewLabel("just a test 12")),
		pinboard.NewPinBoardItem("test13", widget.NewLabel("just a test 13")),
		pinboard.NewPinBoardItem("test14", widget.NewLabel("just a test 14")),
		pinboard.NewPinBoardItem("test15", widget.NewLabel("just a test 15")),
		pinboard.NewPinBoardItem("test16", widget.NewLabel("just a test 16")),
		pinboard.NewPinBoardItem("test17", widget.NewLabel("just a test 17")),
		pinboard.NewPinBoardItem("test18", widget.NewLabel("just a test 18")),
		pinboard.NewPinBoardItem("test19", widget.NewLabel("just a test 19")),
		pinboard.NewPinBoardItem("test20", widget.NewLabel("just a test 20")),
	)

	pb.AddItem(pinboard.NewPinBoardItem("test20", widget.NewLabel("just a test 21")))
	pb.AddItem(pinboard.NewPinBoardItem("test21 - nil detail", nil))
	pb.Show()

	w.Resize(fyne.NewSize(800, 600))

	w.SetContent(scroll)

	w.ShowAndRun()
}

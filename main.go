package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Baby Tracker")

	label := widget.NewLabel("Welcome to Baby Tracker!")
	myWindow.SetContent(container.NewVBox(
		label,
		widget.NewButton("Add Feed", func() {
			label.SetText("Feed logged!")
		}),
	))

	myWindow.ShowAndRun()
}

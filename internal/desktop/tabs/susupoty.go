package tabs

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"babytracker/internal/models"
	"babytracker/internal/storage"
)

// CreateSusuPotyTab creates the diaper tracking interface.
func CreateSusuPotyTab() *fyne.Container {
	dateBinding := binding.NewString()
	timeBinding := binding.NewString()
	notesBinding := binding.NewString()

	diaperTypeSelect := widget.NewSelect(
		[]string{"Wet", "Dirty", "Mixed"},
		func(selected string) {},
	)
	diaperTypeSelect.PlaceHolder = "Select type..."

	dateEntry := widget.NewEntryWithData(dateBinding)
	dateEntry.SetPlaceHolder(dateFormat)
	dateBinding.Set(time.Now().Format(dateFormat))

	timeEntry := widget.NewEntryWithData(timeBinding)
	timeEntry.SetPlaceHolder(timeFormat + " (24hr format)")
	timeBinding.Set(time.Now().Format(timeFormat))

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.Bind(notesBinding)
	notesEntry.SetPlaceHolder("Any observations...")
	notesEntry.Resize(fyne.NewSize(400, 80))

	diaperForm := widget.NewForm(
		&widget.FormItem{Text: "Diaper Type", Widget: diaperTypeSelect},
		&widget.FormItem{Text: "Date", Widget: dateEntry},
		&widget.FormItem{Text: "Time", Widget: timeEntry},
		&widget.FormItem{Text: "Notes", Widget: notesEntry},
	)

	logButton := widget.NewButton("Log Change", func() {
		dateStr, _ := dateBinding.Get()
		if dateStr == "" {
			dateStr = time.Now().Format(dateFormat)
		}

		changeTime := time.Now()
		if timeStr, _ := timeBinding.Get(); timeStr != "" {
			if parsedTime, err := time.Parse(timeFormat, timeStr); err == nil {
				if parsedDate, err := time.Parse(dateFormat, dateStr); err == nil {
					changeTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
						parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, parsedDate.Location())
				}
			}
		}

		notes, _ := notesBinding.Get()

		entry := models.DiaperEntry{
			Date:  dateStr,
			Time:  changeTime,
			Type:  diaperTypeSelect.Selected,
			Notes: notes,
		}

		if err := storage.SaveDiaper(&entry); err != nil {
			fmt.Printf("Error saving diaper change: %v\n", err)
			return
		}

		fmt.Printf("Diaper change logged: %s on %s\n", entry.Type, dateStr)

		diaperTypeSelect.ClearSelected()
		dateBinding.Set(time.Now().Format(dateFormat))
		timeBinding.Set(time.Now().Format(timeFormat))
		notesBinding.Set("")
	})

	quickWetBtn := widget.NewButton("Quick Wet", func() {
		diaperTypeSelect.SetSelected("Wet")
		dateBinding.Set(time.Now().Format(dateFormat))
		timeBinding.Set(time.Now().Format(timeFormat))
	})

	quickDirtyBtn := widget.NewButton("Quick Dirty", func() {
		diaperTypeSelect.SetSelected("Dirty")
		dateBinding.Set(time.Now().Format(dateFormat))
		timeBinding.Set(time.Now().Format(timeFormat))
	})

	quickActions := container.NewHBox(quickWetBtn, quickDirtyBtn)

	recentLabel := widget.NewLabel("Recent Changes")
	recentLabel.TextStyle.Bold = true
	recentPlaceholder := widget.NewLabel("Recent diaper changes will appear here")

	return container.NewVBox(
		widget.NewCard("The Susu-Poty Chronicles", "Log diaper changes",
			container.NewVBox(diaperForm, quickActions, logButton)),
		widget.NewSeparator(),
		widget.NewCard("Recent Activity", "Your recent diaper logs",
			container.NewVBox(recentLabel, recentPlaceholder)),
	)
}

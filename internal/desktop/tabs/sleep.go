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

// CreateSleepTab creates the sleep tracking interface.
func CreateSleepTab() *fyne.Container {
	dateBinding := binding.NewString()
	startTimeBinding := binding.NewString()
	endTimeBinding := binding.NewString()
	notesBinding := binding.NewString()

	sleepTypeSelect := widget.NewSelect(
		[]string{"Nap", "Night"},
		func(selected string) {},
	)
	sleepTypeSelect.PlaceHolder = "Select sleep type..."

	qualitySelect := widget.NewSelect(
		[]string{"Good", "Fair", "Poor"},
		func(selected string) {},
	)
	qualitySelect.PlaceHolder = "Select quality..."

	dateEntry := widget.NewEntryWithData(dateBinding)
	dateEntry.SetPlaceHolder(dateFormat)
	dateBinding.Set(time.Now().Format(dateFormat))

	startTimeEntry := widget.NewEntryWithData(startTimeBinding)
	startTimeEntry.SetPlaceHolder("HH:MM:SS (24hr)")
	startTimeBinding.Set(time.Now().Format(timeFormat))

	endTimeEntry := widget.NewEntryWithData(endTimeBinding)
	endTimeEntry.SetPlaceHolder("HH:MM:SS (24hr, optional)")

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.Bind(notesBinding)
	notesEntry.SetPlaceHolder("Sleep observations...")
	notesEntry.Resize(fyne.NewSize(400, 80))

	sleepForm := widget.NewForm(
		&widget.FormItem{Text: "Sleep Type", Widget: sleepTypeSelect},
		&widget.FormItem{Text: "Date", Widget: dateEntry},
		&widget.FormItem{Text: "Start Time", Widget: startTimeEntry},
		&widget.FormItem{Text: "End Time", Widget: endTimeEntry},
		&widget.FormItem{Text: "Quality", Widget: qualitySelect},
		&widget.FormItem{Text: "Notes", Widget: notesEntry},
	)

	logButton := widget.NewButton("Log Sleep", func() {
		dateStr, _ := dateBinding.Get()
		if dateStr == "" {
			dateStr = time.Now().Format(dateFormat)
		}

		startTime := time.Now()
		if startStr, _ := startTimeBinding.Get(); startStr != "" {
			if parsed, err := time.Parse(timeFormat, startStr); err == nil {
				if parsedDate, err := time.Parse(dateFormat, dateStr); err == nil {
					startTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
						parsed.Hour(), parsed.Minute(), parsed.Second(), 0, parsedDate.Location())
				}
			}
		}

		var endTime time.Time
		var duration int
		if endStr, _ := endTimeBinding.Get(); endStr != "" {
			if parsed, err := time.Parse(timeFormat, endStr); err == nil {
				if parsedDate, err := time.Parse(dateFormat, dateStr); err == nil {
					endTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
						parsed.Hour(), parsed.Minute(), parsed.Second(), 0, parsedDate.Location())
					duration = int(endTime.Sub(startTime).Minutes())
				}
			}
		}

		notes, _ := notesBinding.Get()

		entry := models.SleepEntry{
			Date:      dateStr,
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  duration,
			Type:      sleepTypeSelect.Selected,
			Quality:   qualitySelect.Selected,
			Notes:     notes,
		}

		if err := storage.SaveSleep(&entry); err != nil {
			fmt.Printf("Error saving sleep: %v\n", err)
			return
		}

		fmt.Printf("Sleep logged: %s on %s\n", entry.Type, dateStr)

		sleepTypeSelect.ClearSelected()
		qualitySelect.ClearSelected()
		dateBinding.Set(time.Now().Format(dateFormat))
		startTimeBinding.Set(time.Now().Format(timeFormat))
		endTimeBinding.Set("")
		notesBinding.Set("")
	})

	quickNapBtn := widget.NewButton("Quick Nap", func() {
		sleepTypeSelect.SetSelected("Nap")
		dateBinding.Set(time.Now().Format(dateFormat))
		startTimeBinding.Set(time.Now().Format(timeFormat))
	})

	quickNightBtn := widget.NewButton("Quick Night", func() {
		sleepTypeSelect.SetSelected("Night")
		dateBinding.Set(time.Now().Format(dateFormat))
		startTimeBinding.Set(time.Now().Format(timeFormat))
	})

	quickActions := container.NewHBox(quickNapBtn, quickNightBtn)

	recentLabel := widget.NewLabel("Recent Sleep")
	recentLabel.TextStyle.Bold = true
	recentPlaceholder := widget.NewLabel("Recent sleep history will appear here")

	return container.NewVBox(
		widget.NewCard("Log Sleep", "Track naps and night sleep",
			container.NewVBox(sleepForm, quickActions, logButton)),
		widget.NewSeparator(),
		widget.NewCard("Recent Activity", "Your recent sleep logs",
			container.NewVBox(recentLabel, recentPlaceholder)),
	)
}

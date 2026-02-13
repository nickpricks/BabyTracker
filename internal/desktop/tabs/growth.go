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

// CreateGrowthTab creates the growth tracking interface.
func CreateGrowthTab() *fyne.Container {
	dateBinding := binding.NewString()
	weightBinding := binding.NewFloat()
	heightBinding := binding.NewFloat()
	headCircBinding := binding.NewFloat()
	notesBinding := binding.NewString()

	dateEntry := widget.NewEntryWithData(dateBinding)
	dateEntry.SetPlaceHolder(dateFormat)
	dateBinding.Set(time.Now().Format(dateFormat))

	weightEntry := widget.NewEntryWithData(binding.FloatToString(weightBinding))
	weightEntry.SetPlaceHolder("Weight in kg")

	heightEntry := widget.NewEntryWithData(binding.FloatToString(heightBinding))
	heightEntry.SetPlaceHolder("Height in cm")

	headCircEntry := widget.NewEntryWithData(binding.FloatToString(headCircBinding))
	headCircEntry.SetPlaceHolder("Head circumference in cm")

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.Bind(notesBinding)
	notesEntry.SetPlaceHolder("Growth observations, doctor notes...")
	notesEntry.Resize(fyne.NewSize(400, 80))

	growthForm := widget.NewForm(
		&widget.FormItem{Text: "Date", Widget: dateEntry},
		&widget.FormItem{Text: "Weight (kg)", Widget: weightEntry},
		&widget.FormItem{Text: "Height (cm)", Widget: heightEntry},
		&widget.FormItem{Text: "Head Circ. (cm)", Widget: headCircEntry},
		&widget.FormItem{Text: "Notes", Widget: notesEntry},
	)

	logButton := widget.NewButton("Log Growth", func() {
		dateStr, _ := dateBinding.Get()
		if dateStr == "" {
			dateStr = time.Now().Format(dateFormat)
		}

		weight, _ := weightBinding.Get()
		height, _ := heightBinding.Get()
		headCirc, _ := headCircBinding.Get()
		notes, _ := notesBinding.Get()

		entry := models.GrowthEntry{
			Date:              dateStr,
			Weight:            weight,
			Height:            height,
			HeadCircumference: headCirc,
			Notes:             notes,
		}

		if err := storage.SaveGrowth(&entry); err != nil {
			fmt.Printf("Error saving growth: %v\n", err)
			return
		}

		fmt.Printf("Growth logged for %s\n", dateStr)

		dateBinding.Set(time.Now().Format(dateFormat))
		weightBinding.Set(0)
		heightBinding.Set(0)
		headCircBinding.Set(0)
		notesBinding.Set("")
	})

	recentLabel := widget.NewLabel("Recent Measurements")
	recentLabel.TextStyle.Bold = true
	recentPlaceholder := widget.NewLabel("Recent growth measurements will appear here")

	return container.NewVBox(
		widget.NewCard("Log Growth", "Track weight, height, and head circumference",
			container.NewVBox(growthForm, logButton)),
		widget.NewSeparator(),
		widget.NewCard("Recent Activity", "Your recent growth logs",
			container.NewVBox(recentLabel, recentPlaceholder)),
	)
}

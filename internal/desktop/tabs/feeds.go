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

const (
	dateFormat = time.DateOnly
	timeFormat = time.TimeOnly
)

// CreateFeedsTab creates the feeding tracker interface.
func CreateFeedsTab() *fyne.Container {
	dateBinding := binding.NewString()
	timeBinding := binding.NewString()
	quantityBinding := binding.NewFloat()
	notesBinding := binding.NewString()

	feedTypeSelect := widget.NewSelect(
		[]string{"Bottle", "Breast (Left)", "Breast (Right)", "Breast (Both)", "Solid Food"},
		func(selected string) {
			fmt.Printf("Feed type selected: %s\n", selected)
		},
	)
	feedTypeSelect.PlaceHolder = "Select feed type..."

	dateEntry := widget.NewEntryWithData(dateBinding)
	dateEntry.SetPlaceHolder(dateFormat)
	dateBinding.Set(time.Now().Format(dateFormat))

	timeEntry := widget.NewEntryWithData(timeBinding)
	timeEntry.SetPlaceHolder(timeFormat + " (24hr format)")
	timeBinding.Set(time.Now().Format(timeFormat))

	quantityEntry := widget.NewEntryWithData(binding.FloatToString(quantityBinding))
	quantityEntry.SetPlaceHolder("Amount in ml or oz")

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.Bind(notesBinding)
	notesEntry.SetPlaceHolder("Notes: How did baby respond? Any concerns?")
	notesEntry.Resize(fyne.NewSize(400, 80))

	feedForm := widget.NewForm(
		&widget.FormItem{Text: "Feed Type", Widget: feedTypeSelect},
		&widget.FormItem{Text: "Date", Widget: dateEntry},
		&widget.FormItem{Text: "Time", Widget: timeEntry},
		&widget.FormItem{Text: "Quantity (optional)", Widget: quantityEntry},
		&widget.FormItem{Text: "Notes", Widget: notesEntry},
	)

	logButton := widget.NewButton("Log Feed", func() {
		dateStr, _ := dateBinding.Get()
		if dateStr == "" {
			dateStr = time.Now().Format(dateFormat)
		}

		feedTime := time.Now()
		if timeStr, _ := timeBinding.Get(); timeStr != "" {
			if parsedTime, err := time.Parse(timeFormat, timeStr); err == nil {
				parsedDate, err := time.Parse(dateFormat, dateStr)
				if err == nil {
					feedTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
						parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, parsedDate.Location())
				} else {
					now := time.Now()
					feedTime = time.Date(now.Year(), now.Month(), now.Day(),
						parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, now.Location())
				}
			}
		}

		quantity, _ := quantityBinding.Get()
		notes, _ := notesBinding.Get()

		feed := models.FeedEntry{
			Date:     dateStr,
			Time:     feedTime,
			Type:     feedTypeSelect.Selected,
			Quantity: quantity,
			Notes:    notes,
		}

		err := storage.SaveFeed(&feed)
		if err != nil {
			fmt.Printf("Error saving feed: %v\n", err)
			return
		}

		fmt.Printf("Feed logged successfully at %s %s\n", dateStr, feedTime.Format(timeFormat))

		feedTypeSelect.ClearSelected()
		dateBinding.Set(time.Now().Format(dateFormat))
		timeBinding.Set(time.Now().Format(timeFormat))
		quantityBinding.Set(0)
		notesBinding.Set("")
	})

	quickBottleBtn := widget.NewButton("Quick Bottle", func() {
		feedTypeSelect.SetSelected("Bottle")
		dateBinding.Set(time.Now().Format(dateFormat))
		timeBinding.Set(time.Now().Format(timeFormat))
		quantityEntry.FocusGained()
	})

	quickBreastBtn := widget.NewButton("Quick Breast", func() {
		feedTypeSelect.SetSelected("Breast (Both)")
		dateBinding.Set(time.Now().Format(dateFormat))
		timeBinding.Set(time.Now().Format(timeFormat))
		notesEntry.FocusGained()
	})

	quickActions := container.NewHBox(quickBottleBtn, quickBreastBtn)

	recentFeedsLabel := widget.NewLabel("Recent Feeds")
	recentFeedsLabel.TextStyle.Bold = true
	recentFeedsPlaceholder := widget.NewLabel("Recent feeding history will appear here")

	return container.NewVBox(
		widget.NewCard("Log New Feed", "Track feeding times, amounts, and notes",
			container.NewVBox(feedForm, quickActions, logButton)),
		widget.NewSeparator(),
		widget.NewCard("Recent Activity", "Your recent feeding logs",
			container.NewVBox(recentFeedsLabel, recentFeedsPlaceholder)),
	)
}

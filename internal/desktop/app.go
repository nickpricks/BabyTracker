package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"babytracker/internal/desktop/tabs"
)

// App represents the main application structure.
type App struct {
	fyneApp fyne.App
	window  fyne.Window
}

// NewApp creates and initializes a new Baby Tracker application.
func NewApp() *App {
	myApp := app.New()
	myApp.SetIcon(theme.AccountIcon())

	myWindow := myApp.NewWindow("Baby Tracker")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.CenterOnScreen()

	return &App{
		fyneApp: myApp,
		window:  myWindow,
	}
}

// CreateMainContent creates and returns the main tabbed interface.
func (a *App) CreateMainContent() fyne.CanvasObject {
	feedsTab := tabs.CreateFeedsTab()
	sleepTab := tabs.CreateSleepTab()
	growthTab := tabs.CreateGrowthTab()
	diaperTab := tabs.CreateSusuPotyTab()

	tabsList := container.NewAppTabs(
		container.NewTabItem("Feeds", feedsTab),
		container.NewTabItem("Sleep", sleepTab),
		container.NewTabItem("Growth", growthTab),
		container.NewTabItem("Susu-Poty", diaperTab),
	)
	tabsList.SetTabLocation(container.TabLocationTop)

	return tabsList
}

// SetupWindow configures the main window with content and properties.
func (a *App) SetupWindow() {
	content := a.CreateMainContent()
	a.window.SetContent(content)
	a.window.SetMaster()
	a.window.SetCloseIntercept(func() {
		a.window.Close()
	})
}

// Run starts the Baby Tracker application. Blocks until closed.
func (a *App) Run() {
	a.SetupWindow()
	a.window.ShowAndRun()
}

// GetWindow returns the main application window.
func (a *App) GetWindow() fyne.Window {
	return a.window
}

// GetApp returns the Fyne application instance.
func (a *App) GetApp() fyne.App {
	return a.fyneApp
}

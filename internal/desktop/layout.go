package desktop

import (
	"fyne.io/fyne/v2/container"

	"babytracker/internal/desktop/tabs"
)

// CreateMainLayout constructs the primary tabbed interface.
func CreateMainLayout() *container.AppTabs {
	mainTabs := container.NewAppTabs()

	mainTabs.Append(container.NewTabItem("Feeds", tabs.CreateFeedsTab()))
	mainTabs.Append(container.NewTabItem("Sleep", tabs.CreateSleepTab()))
	mainTabs.Append(container.NewTabItem("Growth", tabs.CreateGrowthTab()))
	mainTabs.Append(container.NewTabItem("Susu-Poty", tabs.CreateSusuPotyTab()))

	return mainTabs
}

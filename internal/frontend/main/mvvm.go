package yaac_frontend_main

import (
	"fmt"

	"fyne.io/fyne/v2"
	pages "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
	settings "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages/settings"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type FrontendMain struct {
	MVVM shared.MVVM
}

func New(mvvm shared.MVVM) *FrontendMain {
	pages.New(mvvm)
	settings.New(mvvm)
	return &FrontendMain{
		MVVM: mvvm,
	}
}

func (*FrontendMain) ReceiveNewTable(table shared.AttendanceList) {
	gv.Window.Show()
	yaac_shared.App.SendNotification(fyne.NewNotification(fmt.Sprintf("%s %o %s", "Es ist ein neue Liste für den Kurs", table.CourseID, "eingetroffen!"), ""))
}

func (*FrontendMain) ReceiveError(source string, err error) {
	yaac_shared.App.SendNotification(fyne.NewNotification(fmt.Sprintf("%s %s %s", "Es ist ein Fehler im", source, "aufgetreten!"), err.Error()))
}

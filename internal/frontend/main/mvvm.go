package yaac_frontend_main

import (
	"fmt"

	"fyne.io/fyne/v2"
	yaac_frontend_pages "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
	yaac_frontend_settings "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages/settings"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type FrontendMain struct {
	MVVM yaac_shared.MVVM
}

func New(mvvm yaac_shared.MVVM) *FrontendMain {
	yaac_frontend_pages.New(mvvm)
	yaac_frontend_settings.New(mvvm)
	return &FrontendMain{
		MVVM: mvvm,
	}
}

func (*FrontendMain) ReceiveNewTable(table yaac_shared.AttendanceList) {
	gv.Window.Show()
	yaac_shared.App.SendNotification(fyne.NewNotification(fmt.Sprintf("%s %o %s", "Es ist ein neue Liste für den Kurs", table.CourseID, "eingetroffen!"), ""))
}

func (*FrontendMain) ReceiveError(source string, err error) {
	yaac_shared.App.SendNotification(fyne.NewNotification(fmt.Sprintf("%s %s %s", "Es ist ein Fehler im", source, "aufgetreten!"), err.Error()))
}

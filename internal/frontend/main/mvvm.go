package yaac_frontend_main

import (
	"fyne.io/fyne/v2/dialog"
	yaac_frontend_pages "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
	yaac_frontend_settings "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages/settings"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type FrontendMain struct {
	MVVM shared.MVVM
}

func New(mvvm shared.MVVM) *FrontendMain {
	yaac_frontend_pages.New(mvvm)
	yaac_frontend_settings.New(mvvm)
	return &FrontendMain{
		MVVM: mvvm,
	}
}

func (*FrontendMain) ReceiveNewTable(table shared.AttendanceList) {
	gv.Window.Show()
	dialog.ShowInformation("Es ist ein neue Liste eingetroffen!", "Nun auf ihrer Übersichtsseite einzusehen!", gv.Window)
}

func (*FrontendMain) ReceiveError(source string, err error) {
	dialog.ShowError(err, gv.Window)
}

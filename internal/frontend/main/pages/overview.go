package yaac_frontend_pages

//TODO: Verification needed when not all attendances = true
import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

var overviewGrid *fyne.Container

func OverviewScreen(w fyne.Window) fyne.CanvasObject {
	title := ReturnHeader("Anwesenheitsliste der Kurse - Heute")

	buttonImageContainer := ReturnVerifyImageContainer()
	insertList := widget.NewButton("", func() {
		OpenImageUpload(w)
	})

	overviewGrid = container.NewGridWrap(fyne.NewSize(250, 250))
	LoadOverviewWidgets(w, overviewGrid)

	header := container.NewVBox(container.NewBorder(nil, nil, nil, container.NewPadded(container.NewPadded(container.NewPadded(container.NewPadded(buttonImageContainer, insertList)))), container.NewGridWrap(fyne.NewSize(400, 200), title)), widget.NewSeparator())
	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(overviewGrid))
}

/*
ReturnVerifyImageContainer returns the buttonImageContaier containing the image for insertList Button.
*/
func ReturnVerifyImageContainer() *fyne.Container {
	image := canvas.NewImageFromResource(yaac_shared.ResourceImageUpload2Png)
	image.FillMode = canvas.ImageFillOriginal
	buttonImageContainer := container.NewCenter(image)
	return buttonImageContainer
}

/*
LoadOverviewWidgets loads all OverviewWidgets for each course and adds them to the overviewGrid
*/
func LoadOverviewWidgets(w fyne.Window, grid *fyne.Container) {
	grid.RemoveAll()
	courses, _ := myMVVM.Courses()
	for _, element := range courses {
		var students []string
		lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: element.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
		if len(lists[0].Attendancies) > 0 {
			students = ReturnNonAttending(lists[0].Attendancies)
		} else {
			students = append(students, "Kein Listeingang")
		}
		grid.Add(NewOverviewWidget(w, element.Name, int(element.ID), students, len(lists[0].Attendancies)))
	}
}

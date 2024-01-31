package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type students struct {
	name   *widget.Label
	course *widget.Label
	year   *widget.Label
}

func studentScreen(_ fyne.Window) fyne.CanvasObject {
	student := &students{
		name:   widget.NewLabel("Max, Alberti"),
		course: widget.NewLabel(""),
		year:   widget.NewLabel(""),
	}
	selection := widget.NewLabel("")
	courseDropdown := widget.NewSelect([]string{
		"TIK",
		"TIT",
	}, func(s string) {
		student.course.SetText(s)
		selection.SetText(updateSelection(student))
	})
	yearDropdown := widget.NewSelect([]string{
		"2021",
		"2022",
		"2023",
	}, func(s string) {
		student.year.SetText(s)
		selection.SetText(updateSelection(student))
	})

	dropdownArea := container.NewHBox(courseDropdown, yearDropdown)
	selectionArea := container.NewVBox(selection, widget.NewSeparator())

	var attendanceData = [][]string{
		[]string{"Datum", "Anwesenheit"},
		[]string{"10.10.2023", "Anwesend"},
		[]string{"11.10.2023", "Anwesend"}}
	attendanceList := widget.NewTable(
		func() (int, int) {
			return len(attendanceData), len(attendanceData[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.TableCellID, cell fyne.CanvasObject) {
			cell.(*widget.Label).SetText(attendanceData[i.Row][i.Col])
		})
	attendanceList.SetColumnWidth(0, 140)
	attendanceList.SetRowHeight(2, 50)
	studentView := container.NewBorder(dropdownArea, nil, nil, nil, container.NewBorder(selectionArea, nil, nil, nil, attendanceList))
	return studentView
}

func updateSelection(student *students) string {
	return fmt.Sprintf("%s - %s %s", student.name.Text, student.course.Text, student.year.Text)
}

package yaac_shared

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
)

const APP_NAME = "YAAC"

var App fyne.App

func GetApp() *fyne.App {
	if App == nil {
		App = app.NewWithID(APP_NAME)
	}

	return &App
}

type Table = imgproc.Table
type TableRow = imgproc.TableRow

type EmailData struct {
	MailServer string
	Email      string
	Password   string
}

type Student struct {
	LName string
	FName string
}

type Attendance struct {
	Student         Student
	Attending       bool
	DayOfAttendance string
}

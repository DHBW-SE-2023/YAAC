package yaac_shared

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
)

const APP_NAME = "YAAC"

var App fyne.App

func GetApp() *fyne.App {
	if App == nil {
		App = app.NewWithID(APP_NAME)
	}

	return &App
}

type EmailData struct {
	MailServer string
	Email      string
	Password   string
}

type Attendance = database.Attendance
type AttendanceList = database.AttendanceList
type Student = database.Student
type Course = database.Course
type Setting = database.Setting

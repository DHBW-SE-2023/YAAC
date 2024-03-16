package yaac_shared

import (
	"time"
)

type MVVM interface {
	// Frontend
	NotifyError(source string, err error)
	NotifyNewList(list AttendanceList)

	// Mail
	UpdateMailCredentials(credentials EmailData) error
	GetMailsToday() ([]MailData, error)
	CheckMailConnection() bool

	// Imgproc
	ValidateTable(img []byte) (Table, error)

	// Database
	InsertList(list AttendanceList) (AttendanceList, error)
	UpdateList(list AttendanceList) (AttendanceList, error)
	LatestList(course Course, date time.Time) (AttendanceList, error)
	InsertCourse(course Course) (Course, error)
	CourseStudents(course Course) ([]Student, error)
	Courses() ([]Course, error)
	AllAttendanceListInRangeByCourse(course Course, start time.Time, end time.Time) ([]AttendanceList, error)
	AllAttendanceListInRange(start time.Time, end time.Time) ([]AttendanceList, error)
	Settings() ([]Setting, error)
	SettingsUpdate(settings []Setting) ([]Setting, error)
	SettingsReset() ([]Setting, error)
	CourseByName(name string) (Course, error)
	Students(student Student) ([]Student, error)
	InsertStudent(student Student) (Student, error)
}

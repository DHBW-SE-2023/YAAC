package yaac_mvvm

import (
	"time"

	database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
)

var backend *database.BackendDatabase

func (m *MVVM) ConnectDatabase(dbPath string) error {
	backend = database.NewBackend(m, dbPath)
	return backend.ConnectDatabase()
}

// This function allows only one list per day. If a list already exists, it overwrites it.
// AttendanceList must be filled out fully, except the ID, CreatedAt, UpdatedAt, and DeletedAt.
func (m *MVVM) InsertList(list database.AttendanceList) (database.AttendanceList, error) {
	return backend.InsertList(list)
}

// `list` needs the field `Id` to be not null.
func (m *MVVM) UpdateList(list database.AttendanceList) (database.AttendanceList, error) {
	return backend.UpdateList(list)
}

// [..., end)
func (m *MVVM) LatestList(course database.Course, end time.Time) (database.AttendanceList, error) {
	return backend.LatestList(course, end)
}

func (m *MVVM) InsertCourse(course database.Course) (database.Course, error) {
	return backend.InsertCourse(course)
}

func (m *MVVM) Courses() ([]database.Course, error) {
	return backend.Courses()
}

// [start, end)
// Returns all lists, even outdated ones
func (m *MVVM) AllAttendanceListInRangeByCourse(course database.Course, start time.Time, end time.Time) ([]database.AttendanceList, error) {
	return backend.AllAttendanceListInRangeByCourse(course, start, end)
}

// [start, end)
// Returns all lists, even outdated ones
func (m *MVVM) AllAttendanceListInRange(start time.Time, end time.Time) ([]database.AttendanceList, error) {
	return backend.AllAttendanceListInRange(start, end)
}

func (m *MVVM) Settings() ([]database.Setting, error) {
	return backend.Settings()
}

func (m *MVVM) SettingsUpdate(settings []database.Setting) ([]database.Setting, error) {
	return backend.SettingsUpdate(settings)
}

func (m *MVVM) SettingsReset() ([]database.Setting, error) {
	return backend.SettingsReset()
}

func (m *MVVM) InsertStudent(student database.Student) (database.Student, error) {
	return backend.InsertStudent(student)
}

func (m *MVVM) CourseStudents(course database.Course) ([]database.Student, error) {
	return backend.CourseStudents(course)
}

func (m *MVVM) CourseByName(name string) (database.Course, error) {
	return backend.CourseByName(name)
}

func (m *MVVM) Students(student database.Student) ([]database.Student, error) {
	return backend.Students(student)
}

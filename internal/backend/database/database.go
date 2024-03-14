package yaac_backend_database

import (
	"log"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Course struct {
	gorm.Model
	Name     string
	Students []Student
}

type Student struct {
	gorm.Model
	FirstName        string `gorm:"check:FirstName!=''"`
	LastName         string `gorm:"check:LastName!=''"`
	CourseID         uint
	IsImmatriculated bool
}

type Attendance struct {
	StudentID        uint `gorm:"primaryKey"`
	AttendanceListID uint `gorm:"primaryKey"` //`gorm:"primaryKey;foreignKey:Id;references:AttendanceList"`
	IsAttending      bool
	NameROI          Rectangle
	SignatureROI     Rectangle
	TotalROI         Rectangle
}

type AttendanceList struct {
	ID uint `gorm:"primaryKey"`
	// CreatedAt is a primary key to ensure, that only one list per day can exist
	CreatedAt    time.Time `gorm:"primaryKey"`
	UpdatedAt    time.Time
	CourseID     uint //`gorm:"foreignKey:Id;references:Course"`
	ReceivedAt   time.Time
	Attendancies []Attendance
	Image        []byte // image as png
}

type Setting struct {
	gorm.Model
	Setting string `gorm:"index"`
	Value   string
}

func (item *BackendDatabase) ConnectDatabase() error {
	// We save private data, so noone but us may read it
	dbPath := item.Path

	err := os.MkdirAll(path.Dir(dbPath), 0700)
	if err != nil {
		log.Fatalf("Could not create the database: %v", err)
		return err
	}

	// Ensure that the file exists
	fd, err := os.OpenFile(dbPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err == nil {
		fd.Close()
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
		return err
	}

	item.DB = db

	return item.DB.AutoMigrate(&Course{}, &Student{}, &AttendanceList{}, &Attendance{}, &Setting{})
}

// This function allows only one list per day. If a list already exists, it overwrites it.
// AttendanceList must be filled out fully, except the ID.
func (item *BackendDatabase) InsertList(list AttendanceList) (AttendanceList, error) {
	year, month, day := list.ReceivedAt.Date()
	list.CreatedAt = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	_ = item.DB.Model(&AttendanceList{}).Delete(&AttendanceList{}, item.DB.Where(&AttendanceList{CreatedAt: list.CreatedAt, CourseID: list.CourseID}))
	err := item.DB.Model(&AttendanceList{}).Create(&list).Error
	return list, err
}

// `list` needs the field `Id` to be not null.
func (item *BackendDatabase) UpdateList(list AttendanceList) (AttendanceList, error) {
	err := item.DB.Save(&list).Error
	return list, err
}

// [..., end)
func (item *BackendDatabase) LatestList(course Course, end time.Time) (AttendanceList, error) {
	list := AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Preload("Attendancies").Joins("JOIN Course c ON c.ID = CourseID").Where("ReceivedAt < ?", end).Order("ReceivedAt DESC").Take(&list).Error
	return list, err
}

// [start, end)
// Returns all lists, even outdated ones
func (item *BackendDatabase) AllAttendanceListInRangeByCourse(course Course, start time.Time, end time.Time) ([]AttendanceList, error) {
	list := []AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Preload("Attendancies").Where("CourseID = ?", course.ID).Where("ReceivedAt BETWEEN ? AND ?", start, end).Order("ReceivedAt DESC").Find(&list).Error
	return list, err
}

// [start, end)
// Returns all lists, even outdated ones
func (item *BackendDatabase) AllAttendanceListInRange(start time.Time, end time.Time) ([]AttendanceList, error) {
	list := []AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Preload("Attendancies").Where("ReceivedAt BETWEEN ? AND ?", start, end).Order("ReceivedAt DESC").Find(&list).Error
	return list, err
}

func (item *BackendDatabase) InsertCourse(course Course) (Course, error) {
	err := item.DB.Model(&Course{}).Save(&course).Error
	return course, err
}

func (item *BackendDatabase) Courses() ([]Course, error) {
	courses := []Course{}
	err := item.DB.Model(&Course{}).Find(&courses).Error
	return courses, err
}

func (item *BackendDatabase) CourseByName(name string) (Course, error) {
	course := Course{}
	err := item.DB.Model(&Course{}).Where(&Course{Name: name}).Take(&course).Error
	return course, err
}

func (item *BackendDatabase) CourseStudents(course Course) ([]Student, error) {
	students := []Student{}
	err := item.DB.Model(&Course{}).Joins("JOIN Student ON Course.ID = Student.CourseID").Where(&course).Select("Student.*").Find(&students).Error
	return students, err
}

func (item *BackendDatabase) InsertStudent(student Student) (Student, error) {
	err := item.DB.Model(&Student{}).Save(&student).Error
	return student, err
}

func (item *BackendDatabase) Settings() ([]Setting, error) {
	settings := []Setting{}
	err := item.DB.Model(&Setting{}).Find(&settings).Error
	return settings, err
}

func (item *BackendDatabase) SettingsUpdate(settings []Setting) ([]Setting, error) {
	err := item.DB.Model(&Setting{}).Save(&settings).Error
	return settings, err
}

func (item *BackendDatabase) SettingsReset() ([]Setting, error) {
	// FIXME: Add default settings
	settings := []Setting{}
	err := item.DB.Model(&Setting{}).Delete(&Setting{}).Create(&settings).Error
	return settings, err
}

func (item *BackendDatabase) Students(student Student) ([]Student, error) {
	students := []Student{}
	err := item.DB.Model(&Student{}).Where(student).Find(&students).Error
	return students, err
}

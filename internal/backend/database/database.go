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
	Name     string `gorm:"unique"`
	Students []Student
}

type Student struct {
	gorm.Model
	FirstName        string `gorm:"check:FirstName!='';type:varchar(64)"`
	LastName         string `gorm:"check:LastName!='';type:varchar(64)"`
	FullName         string `gorm:"check:FullName!='';type:varchar(128);uniqueIndex:namecourseunique;index;not null"`
	CourseID         uint   `gorm:"uniqueIndex:namecourseunique;not null"`
	IsImmatriculated bool   `gorm:"default:true"`
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
	Setting string `gorm:"primaryKey"`
	Value   string
}

const DEFAULT_ORDER = "ReceivedAt DESC"

// Ensure that the database is present and all tables being there
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
		FullSaveAssociations: true,
		AllowGlobalUpdate:    true,
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

// `list` needs the field `Id` to be	 not null.
func (item *BackendDatabase) UpdateList(list AttendanceList) (AttendanceList, error) {
	err := item.DB.Save(&list).Error
	return list, err
}

// Get the latest list before a certain date for a course
// [..., end)
func (item *BackendDatabase) LatestList(course Course, end time.Time) (AttendanceList, error) {
	list := AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Preload("Attendancies").Joins("JOIN Course c ON c.ID = CourseID").Where("ReceivedAt < ?", end).Order(DEFAULT_ORDER).Take(&list).Error
	return list, err
}

// Get all attendance lists for a course in a time range.
//
// [start, end)
func (item *BackendDatabase) AllAttendanceListInRangeByCourse(course Course, start time.Time, end time.Time) ([]AttendanceList, error) {
	list := []AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Preload("Attendancies").Where("CourseID = ?", course.ID).Where("ReceivedAt BETWEEN ? AND ?", start, end).Order(DEFAULT_ORDER).Find(&list).Error
	return list, err
}

// Get all attendance lists for aall courses in a time range.
//
// [start, end)
func (item *BackendDatabase) AllAttendanceListInRange(start time.Time, end time.Time) ([]AttendanceList, error) {
	list := []AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Preload("Attendancies").Where("ReceivedAt BETWEEN ? AND ?", start, end).Order(DEFAULT_ORDER).Find(&list).Error
	return list, err
}

// Insert a new course. It is not checked whether the course already exists.
func (item *BackendDatabase) InsertCourse(course Course) (Course, error) {
	err := item.DB.Model(&Course{}).Save(&course).Error
	return course, err
}

// Get all courses saved in the database.
func (item *BackendDatabase) Courses() ([]Course, error) {
	courses := []Course{}
	err := item.DB.Model(&Course{}).Find(&courses).Error
	return courses, err
}

// Get a course by name.
func (item *BackendDatabase) CourseByName(name string) (Course, error) {
	course := Course{}
	err := item.DB.Model(&Course{}).Where(&Course{Name: name}).Take(&course).Error
	return course, err
}

// Get all students in a `course`.
func (item *BackendDatabase) CourseStudents(course Course) ([]Student, error) {
	students := []Student{}
	err := item.DB.Model(&Course{}).Joins("JOIN Student ON Course.ID = Student.CourseID").Where(&course).Select("Student.*").Find(&students).Error
	return students, err
}

// Add a new student to the database.
// student.FullName will be set automatically.
func (item *BackendDatabase) InsertStudent(student Student) (Student, error) {
	if student.FullName == "" {
		student.FullName = student.LastName + ", " + student.FirstName
	}

	err := item.DB.Model(&Student{}).Save(&student).Error
	return student, err
}

// Get all settings saved in the database
func (item *BackendDatabase) Settings() ([]Setting, error) {
	settings := []Setting{}
	err := item.DB.Model(&Setting{}).Find(&settings).Error
	return settings, err
}

// Update settings. `settings` need not contain all key-value pairs.
func (item *BackendDatabase) SettingsUpdate(settings []Setting) ([]Setting, error) {
	err := item.DB.Model(&Setting{}).Save(&settings).Error
	return settings, err
}

// Reset all settings. This clears the `Setting` table
func (item *BackendDatabase) SettingsReset() ([]Setting, error) {
	settings := []Setting{
		{Setting: "MailServer", Value: "imap.mail.de:993"},
		{Setting: "UserEmail", Value: "anwesenheits_listen@mail.de"},
		{Setting: "UserEmailPassword", Value: "DHBW-YAAC-2024!"},
	}
	err := item.DB.Model(&Setting{}).Delete(&Setting{}).Create(&settings).Error
	return settings, err
}

// Query all students based on `student`.

func (item *BackendDatabase) Students(student Student) ([]Student, error) {
	students := []Student{}
	err := item.DB.Model(&Student{}).Where(student).Find(&students).Error
	return students, err
}

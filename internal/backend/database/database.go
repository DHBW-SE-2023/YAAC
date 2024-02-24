package yaac_backend_database

import (
	"log"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name     string
	Students []Student `gorm:"foreignKey:Course"`
}

type Student struct {
	gorm.Model
	FirstName        string
	LastName         string
	Course           uint
	IsImmatriculated bool
}

type Attendance struct {
	gorm.Model
	Student     uint `gorm:"primaryKey;foreignKey:Id;references:Student"`
	List        uint `gorm:"primaryKey"`
	IsAttending bool
}

type AttendanceList struct {
	gorm.Model
	Course       uint `gorm:"foreignKey:Id;references:Course"`
	ReceivedAt   time.Time
	Attendencies []Attendance `gorm:"foreignKey:List"`
}

type Setting struct {
	gorm.Model
	Setting string
	Value   any
}

func (item *BackendDatabase) ConnectDatabase(dbPath string) error {
	// We save private data, so noone but us may read it
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

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
		return err
	}

	item.DB = db

	db.AutoMigrate(&Course{}, &Student{}, &AttendanceList{}, &Attendance{}, &Setting{})

	return nil
}

func (item *BackendDatabase) InsertList(list AttendanceList) (AttendanceList, error) {
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
	err := item.DB.Model(&AttendanceList{}).Where(&AttendanceList{Course: course.ID}).Where("ReceivedAt < ?", end).Order("ReceivedAt DESC").Take(&list).Error
	return list, err
}

// [start, end)
// Returns all lists, even outdated ones
func (item *BackendDatabase) AllAttendanceListInRangeByCourse(course Course, start time.Time, end time.Time) ([]AttendanceList, error) {
	list := []AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Where(&AttendanceList{Course: course.ID}).Where("ReceivedAt BETWEEN ? AND ?", start, end).Order("ReceivedAt DESC").Find(&list).Error
	return list, err
}

// [start, end)
// Returns all lists, even outdated ones
func (item *BackendDatabase) AllAttendanceListInRange(start time.Time, end time.Time) ([]AttendanceList, error) {
	list := []AttendanceList{}
	err := item.DB.Model(&AttendanceList{}).Where("ReceivedAt BETWEEN ? AND ?", start, end).Order("ReceivedAt DESC").Find(&list).Error
	return list, err
}

func (item *BackendDatabase) Courses() ([]Course, error) {
	courses := []Course{}
	err := item.DB.Model(&Course{}).Find(&courses).Error
	return courses, err
}

func (item *BackendDatabase) CourseStudents(course Course) ([]Student, error) {
	students := []Student{}
	err := item.DB.Model(&Course{}).Where(course).Select("Students").Find(&students).Error
	return students, err
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
	settings := []Setting{}
	err := item.DB.Model(&Setting{}).Delete(&Setting{}).Create(&settings).Error
	return settings, err
}

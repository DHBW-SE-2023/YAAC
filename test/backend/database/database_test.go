package test

import (
	"log"
	"os"
	"testing"
	"time"

	backend "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	mvvm "github.com/DHBW-SE-2023/YAAC/internal/mvvm"
	"gorm.io/gorm"
)

var testTime = time.Date(2000, 01, 01, 1, 1, 1, 1, time.UTC)
var testByteArray = []byte("This is a test string simulating the byte array from a file")

func setupDatabase() (*backend.BackendDatabase, error) {
	dbPath := "./testdb.db"
	mvvm := mvvm.New()
	conn := backend.NewBackend(mvvm, dbPath)
	var err error = nil
	where := ""

	where = "ConnectDatabase"
	if err = conn.ConnectDatabase(); err != nil {
		goto Cleanup
	}

	{

		where = "InsertCourse"
		courseNames := []string{"TIK22", "TIT22", "TIS22", "TIM22"}
		for _, course := range courseNames {
			c := backend.Course{Name: course}
			if _, err = conn.InsertCourse(c); err != nil {
				goto Cleanup
			}
		}

		students := []backend.Student{
			{
				FirstName:        "Max",
				LastName:         "Mustermann",
				IsImmatriculated: true,
			},
			{
				FirstName:        "Maximilian",
				LastName:         "Mustermann",
				IsImmatriculated: true,
			},
		}

		where = "InsertStudents"
		for _, student := range students {
			if _, err = conn.InsertStudent(student); err != nil {
				goto Cleanup
			}
		}

		attendanceLists := []backend.AttendanceList{
			{
				ReceivedAt: testTime,
				CourseID:   0, // TIK22
				Image:      testByteArray,
				Attendancies: []backend.Attendance{
					{
						StudentID:   0, // Max Mustermann
						IsAttending: false,
					},
					{
						StudentID:   1, // Maximilian Mustermann
						IsAttending: true,
					},
				},
			},
			{
				ReceivedAt: testTime,
				CourseID:   0, // TIK22
				Image:      testByteArray,
				Attendancies: []backend.Attendance{
					{
						StudentID:   0, // Max Mustermann
						IsAttending: true,
					},
					{
						StudentID:   1, // Maximilian Mustermann
						IsAttending: false,
					},
				},
			},
			{
				ReceivedAt: testTime,
				CourseID:   1, // TIT22
				Image:      testByteArray,
				Attendancies: []backend.Attendance{
					{
						StudentID:   0, // Max Mustermann
						IsAttending: true,
					},
					{
						StudentID:   1, // Maximilian Mustermann
						IsAttending: false,
					},
				},
			},
		}

		where = "InsertAttendanceList"
		for _, list := range attendanceLists {
			if _, err = conn.InsertList(list); err != nil {
				goto Cleanup
			}
		}
	}

	return conn, nil

Cleanup:
	log.Fatalf("Could not setup database, where %v, error: %v", where, err)
	os.Remove(dbPath)
	return &backend.BackendDatabase{}, err
}

func clearDatabase(t *testing.T, conn *backend.BackendDatabase) {
	if err := os.Remove(conn.Path); err != nil {
		t.Fatalf("Remove database: %v", err)
	}

}

func TestConnectDatabase(t *testing.T) {
	conn, err := setupDatabase()
	if err != nil {
		t.Fatalf("Couldn't create the database: %v", err)
		return
	}

	defer clearDatabase(t, conn)

	if _, err = os.Stat(conn.Path); err != nil {
		t.Fatalf("Could not access database: %v", err)
		return
	}
}

func TestInsertAttendanceList(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	list, _ := conn.LatestList(backend.Course{Name: "TIK22"}, time.Now())

	correctAttendancies := []bool{true, false}
	if len(list.Attendancies) != len(correctAttendancies) {
		t.Fatalf("Length of attendancies is incorrect: %v, correct: %v", len(list.Attendancies), len(correctAttendancies))
	}

	// Check if the state of the attendancies has updated
	for i, attendancy := range list.Attendancies {
		if attendancy.IsAttending != correctAttendancies[i] {
			s, _ := conn.Students(backend.Student{Model: gorm.Model{ID: attendancy.StudentID}})
			t.Fatalf("StudentID %v is incorrectly marked as %v, correct: %v", s[0], attendancy.IsAttending, correctAttendancies[i])
		}
	}
}

func TestStudents(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	students := []backend.Student{
		{
			FirstName:        "Max",
			LastName:         "Mustermann",
			IsImmatriculated: true,
		},
		{
			FirstName:        "Maximilian",
			LastName:         "Mustermann",
			IsImmatriculated: true,
		},
	}
	s, _ := conn.Students(backend.Student{LastName: "Mustermann"})

	if len(s) != len(students) {
		t.Fatalf("Length of students is incorrect: %v, correct: %v", len(s), len(students))
	}

	for i, student := range students {
		if s[i].CourseID != student.CourseID || s[i].IsImmatriculated != student.IsImmatriculated || s[i].LastName != student.LastName || s[i].FirstName != student.FirstName {
			t.Fatalf("Student mismatched in core properties: %v, correct: %v", s[i], student)
		}
	}
}

func TestAllAttendanceListInRangeByCourse(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	correctLists := []backend.AttendanceList{
		{
			ReceivedAt: testTime.UTC(),
			CourseID:   0, // TIK22
			Image:      testByteArray,
			Attendancies: []backend.Attendance{
				{
					StudentID:   0, // Max Mustermann
					IsAttending: true,
				},
				{
					StudentID:   1, // Maximilian Mustermann
					IsAttending: false,
				},
			},
		},
	}

	lists, err := conn.AllAttendanceListInRangeByCourse(backend.Course{Model: gorm.Model{ID: 0}}, testTime, testTime.Add(1000))
	if err != nil {
		t.Fatalf("AllAttendanceListInRangeByCourse: %v", err)
	}

	if len(correctLists) != len(lists) {
		t.Fatalf("Incorrect length: %v, correct: %v", len(lists), len(correctLists))
	}

	for i, l := range lists {
		c := correctLists[i]
		if c.ReceivedAt != l.ReceivedAt || c.CourseID != l.CourseID || len(c.Attendancies) != len(l.Attendancies) {
			t.Fatalf("Mismatched properties")
		}

		for j, a := range l.Attendancies {
			b := c.Attendancies[j]
			if a.IsAttending != b.IsAttending || a.StudentID != b.StudentID {
				t.Fatalf("Mismatched attendancies properties")
			}
		}
	}
}

func TestLatestList(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	c, _ := conn.CourseByName("TIK22")

	correctList := backend.AttendanceList{
		ReceivedAt: testTime.UTC(),
		CourseID:   c.ID, // TIK22
		Image:      testByteArray,
		Attendancies: []backend.Attendance{
			{
				StudentID:   0, // Max Mustermann
				IsAttending: true,
			},
			{
				StudentID:   1, // Maximilian Mustermann
				IsAttending: false,
			},
		},
	}

	latestList, _ := conn.LatestList(c, time.Now())
	if len(correctList.Attendancies) != len(latestList.Attendancies) || correctList.ReceivedAt != latestList.ReceivedAt || correctList.CourseID != latestList.CourseID {
		t.Fatalf("Mismatched properties")
	}
}

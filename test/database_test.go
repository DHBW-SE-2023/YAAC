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

const AACERROR = "AllAttendanceListInRangeByCourse: %v"
const MPERROR = "Mismatched properties"

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

		tikC, _ := conn.CourseByName("TIK22")
		titC, _ := conn.CourseByName("TIT22")

		students := []backend.Student{
			{
				FirstName:        "Max",
				LastName:         "Mustermann",
				IsImmatriculated: true,
				CourseID:         tikC.ID,
			},
			{
				FirstName:        "Maximilian",
				LastName:         "Mustermann",
				IsImmatriculated: true,
				CourseID:         titC.ID,
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
				CourseID:   tikC.ID, // TIK22
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
			{ // This list should automatically update the previous one.
				ReceivedAt: testTime,
				CourseID:   tikC.ID, // TIK22
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
				CourseID:   titC.ID, // TIT22
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

func TestInsertTwoStudentsWithSameName(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	c, _ := conn.CourseByName("TIK22")

	// This student already exists in the database
	s := backend.Student{
		FirstName:        "Max",
		LastName:         "Mustermann",
		IsImmatriculated: true,
		CourseID:         c.ID,
	}

	_, err := conn.InsertStudent(s)
	if err == nil {
		t.Fatalf("Expected to fail. Two students in the same course may not have the same name")
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

	tikC, _ := conn.CourseByName("TIK22")
	titC, _ := conn.CourseByName("TIT22")

	students := []backend.Student{
		{
			FirstName:        "Max",
			LastName:         "Mustermann",
			IsImmatriculated: true,
			CourseID:         tikC.ID,
		},
		{
			FirstName:        "Maximilian",
			LastName:         "Mustermann",
			IsImmatriculated: true,
			CourseID:         titC.ID,
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

func listsEqual(a backend.AttendanceList, b backend.AttendanceList) bool {
	return a.ReceivedAt == b.ReceivedAt && a.CourseID == b.CourseID && len(a.Attendancies) == len(b.Attendancies)
}

func attendanciesEqual(a backend.Attendance, b backend.Attendance) bool {
	return a.IsAttending == b.IsAttending && a.StudentID == b.StudentID
}

func TestAllAttendanceListInRangeByCourse(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	tikC, _ := conn.CourseByName("TIK22")

	correctLists := []backend.AttendanceList{
		{
			ReceivedAt: testTime.UTC(),
			CourseID:   tikC.ID,
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

	lists, err := conn.AllAttendanceListInRangeByCourse(tikC, testTime, testTime.Add(1000))
	if err != nil {
		t.Fatalf(AACERROR, err)
	}

	if len(correctLists) != len(lists) {
		t.Fatalf("Incorrect length: %v, correct: %v", len(lists), len(correctLists))
	}

	for i, l := range lists {
		c := correctLists[i]
		if !listsEqual(l, c) {
			t.Fatalf(MPERROR)
		}

		for j, a := range l.Attendancies {
			b := c.Attendancies[j]
			if !attendanciesEqual(a, b) {
				t.Fatalf("Mismatched attendancies properties")
			}
		}
	}
}

func TestAllAttendanceListInRange(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	tikC, _ := conn.CourseByName("TIK22")
	titC, _ := conn.CourseByName("TIT22")

	correctLists := []backend.AttendanceList{
		{
			ReceivedAt: testTime.UTC(),
			CourseID:   tikC.ID,
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
			CourseID:   titC.ID,
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

	lists, err := conn.AllAttendanceListInRange(testTime, testTime.Add(1000))
	if err != nil {
		t.Fatalf("AllAttendanceListInRange: %v", err)
	}

	if len(correctLists) != len(lists) {
		t.Fatalf("Incorrect length: %v, correct: %v", len(lists), len(correctLists))
	}

	for i, l := range lists {
		c := correctLists[i]
		if !listsEqual(l, c) {
			t.Fatalf(MPERROR)
		}

		for j, a := range l.Attendancies {
			b := c.Attendancies[j]
			if !attendanciesEqual(a, b) {
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
		t.Fatalf(MPERROR)
	}
}

func TestUpdateList(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	c, _ := conn.CourseByName("TIK22")
	lists, err := conn.AllAttendanceListInRangeByCourse(c, testTime, testTime.Add(1000))

	if err != nil {
		t.Fatalf(AACERROR, err)
	}

	if len(lists) != 1 {
		t.Fatalf("Wrong number of lists: %v, correct: 1", len(lists))
	}

	originalList := lists[0]
	newCourse, _ := conn.CourseByName("TIM22")

	newList := originalList
	newList.CourseID = newCourse.ID

	conn.UpdateList(newList)

	c, _ = conn.CourseByName("TIM22")
	lists, err = conn.AllAttendanceListInRangeByCourse(c, testTime, testTime.Add(1000))

	if err != nil {
		t.Fatalf(AACERROR, err)
	}

	if len(lists) != 1 {
		t.Fatalf("Wrong number of lists: %v, correct: 1", len(lists))
	}

	l := lists[0]

	if l.CourseID != newCourse.ID {
		t.Fatalf("CourseID property not updated: %v, correct: %v", l.CourseID, newCourse.ID)
	}
}

func TestCourseStudents(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	c, _ := conn.CourseByName("TIK22")

	correctStudents := []backend.Student{
		{
			FirstName:        "Max",
			LastName:         "Mustermann",
			IsImmatriculated: true,
			CourseID:         c.ID,
		},
	}

	students, err := conn.CourseStudents(c)

	if err != nil {
		t.Fatalf("CourseStudents: %v", err)
	}

	if len(correctStudents) != len(students) {
		t.Fatalf("Wrong number of students: %v, correct: %v", len(students), len(correctStudents))
	}

	for i, s := range students {
		c := correctStudents[i]

		if !studentsEqual(s, c) {
			t.Logf("%+v\n %+v\n", s, c)
			t.Fatalf("Mismatched properties: index %v", i)
		}
	}
}

func studentsEqual(a, b backend.Student) bool {
	return a.FirstName == b.FirstName && a.LastName == b.LastName && a.IsImmatriculated == b.IsImmatriculated && a.CourseID == b.CourseID
}

func settingsEqual(a, b backend.Setting) bool {
	return a.Setting == b.Setting && a.Value == b.Value
}

func TestSettings(t *testing.T) {
	conn, _ := setupDatabase()
	defer clearDatabase(t, conn)

	correctSettings := []backend.Setting{
		{
			Setting: "test1",
			Value:   "123",
		},
		{
			Setting: "test2",
			Value:   "456",
		},
	}

	_, err := conn.SettingsUpdate(correctSettings)
	if err != nil {
		t.Fatalf("SettingsUpdate: %v", err)
	}

	settings, err := conn.Settings()
	if err != nil {
		t.Fatalf("Settings: %v", err)
	}

	if len(settings) != len(correctSettings) {
		t.Fatalf("Wrong number of settings: %v, correct: %v", len(settings), len(correctSettings))
	}

	for i, s := range settings {
		c := correctSettings[i]
		if !settingsEqual(s, c) {
			t.Fatalf("Mismatched properties: index %v", i)
		}
	}
}

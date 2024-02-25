package database

import (
	yaac_backend_database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	yaac_mvvm "github.com/DHBW-SE-2023/YAAC/internal/mvvm"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

/*	The focus of this test are the statements actually executed on the database.
	Functionality of the actual database and the libraries included are not to be tested here
	thus the test coverage for error handling statements is not great.
*/

// test variables
var location, _ = time.LoadLocation("CET")
var testTime = time.Date(2000, 01, 01, 1, 1, 1, 1, location)
var testByteArray = []byte("This is a test string simulating the byte array from a file")

// setupDatabaseHelper sets up the test database
func setupDatabaseHelper() (*yaac_backend_database.BackendDatabase, error) {
	mvvm := yaac_mvvm.New()

	databaseinst := yaac_backend_database.New(mvvm, "../../testdata/", "data.db")

	databaseinst.CreateDatabase()
	databaseinst.ConnectDatabase()

	// setup test data
	err := databaseinst.InsertCourse("TIK22")
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertCourse("TIT22")
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertCourse("TIS22")
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertStudent("Mustermann", "Max", true, 1)
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertStudent("Muster", "Maximilian", false, 2)
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertAttendance(1, testTime, true)
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertAttendance(2, testTime, false)
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertAttendanceList(testTime, 1, testByteArray)
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}
	err = databaseinst.InsertAttendanceList(testTime, 2, testByteArray)
	if err != nil {
		err := os.Remove("../../testdata/data.db")
		return nil, err
	}

	return databaseinst, nil
}

// clearDatabaseHelper deletes the test database to achieve a persistent environment throughout tests
func clearDatabaseHelper(databaseinst *yaac_backend_database.BackendDatabase) {
	databaseinst.DisconnectDatabase()
	err := os.Remove("../../testdata/data.db")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateDatabase(t *testing.T) {
	databaseinst, err := setupDatabaseHelper()
	if err != nil {
		t.Error("Test Create Database failed")
	}
	defer clearDatabaseHelper(databaseinst)

	if _, err := os.Stat("../../testdata/" + "data.db"); err == nil {
		t.Log("CreateDatabase() tested successfully")
	}
}

func TestConnectDatabase(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	t.Log("ConnectDatabase() tested successfully")
}

func TestDisconnectDatabase(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()

	databaseinst.DisconnectDatabase()

	if _, err := databaseinst.GetAllAttendanceWithStudentName(); err == nil {
		t.Errorf("Test Disconnect Database failed: %q", err)
	} else {
		t.Log("Disconnect Database tested successfully")
	}

	err := os.Remove("../../testdata/data.db")
	if err != nil {
		log.Fatal(err)
	}
}

func TestInsertAttendance(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	// the insert statements already run in setupDatabaseHelper thus we only need to check if the entries are present
	attendingStudents, _ := databaseinst.GetAllAttendanceWithStudentName()

	students := []yaac_shared.Attendance{
		{Student: yaac_shared.Student{
			LName: "Mustermann",
			FName: "Max",
		},
			Attending:       true,
			DayOfAttendance: testTime.Format("2006-01-02"),
		},
		{Student: yaac_shared.Student{
			LName: "Muster",
			FName: "Maximilian",
		},
			Attending:       false,
			DayOfAttendance: testTime.Format("2006-01-02"),
		},
	}

	if !reflect.DeepEqual(attendingStudents, students) {
		t.Error("Failed to test InsertAttendance")
	} else {
		t.Log("Insert Attendance tested successfully")
	}
}

func TestInsertCurrentAttendance(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	// the function uses time.Now() as date, unfortunately this can not be mocked in go
	// for this reason the test assumes that it runs on the same day (as the method uses only the current day)
	timeNow := time.Now()

	err := databaseinst.InsertCurrentAttendance(1, true)
	if err != nil {
		t.Errorf("Failed to test InsertCurrentAttendance with error: %q", err)
	}

	attendingStudents, _ := databaseinst.GetAllAttendanceWithStudentName()

	students := []yaac_shared.Attendance{
		{Student: yaac_shared.Student{
			LName: "Mustermann",
			FName: "Max",
		},
			Attending:       true,
			DayOfAttendance: testTime.Format("2006-01-02"),
		},
		{Student: yaac_shared.Student{
			LName: "Muster",
			FName: "Maximilian",
		},
			Attending:       false,
			DayOfAttendance: testTime.Format("2006-01-02"),
		},
		{Student: yaac_shared.Student{
			LName: "Mustermann",
			FName: "Max",
		},
			Attending:       true,
			DayOfAttendance: timeNow.Format("2006-01-02"),
		},
	}

	if !reflect.DeepEqual(attendingStudents, students) {
		t.Error("Failed to test InsertCurrentAttendance")
	} else {
		t.Log("Insert Attendance tested successfully")
	}
}

func TestUpdateAttendance(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	err := databaseinst.UpdateAttendance(1, testTime, false)
	if err != nil {
		t.Errorf("Failed to test UpdateAttendance with error: %q", err)
	}

	attendingStudents, _ := databaseinst.GetAllAttendanceWithStudentName()

	students := []yaac_shared.Attendance{
		{Student: yaac_shared.Student{
			LName: "Mustermann",
			FName: "Max",
		},
			Attending:       false,
			DayOfAttendance: testTime.Format("2006-01-02"),
		},
		{Student: yaac_shared.Student{
			LName: "Muster",
			FName: "Maximilian",
		},
			Attending:       false,
			DayOfAttendance: testTime.Format("2006-01-02"),
		},
	}

	if !reflect.DeepEqual(attendingStudents, students) {
		t.Error("Failed to test InsertAttendance")
	} else {
		t.Log("Insert Attendance tested successfully")
	}
}

func TestInsertStudent(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	err := databaseinst.InsertStudent("Maximus", "Musterus", true, 3)
	if err != nil {
		t.Errorf("Failed to test InsertStudent with error: %q", err)
	}

	studentsOfTis22, _ := databaseinst.GetAllStudentsPerCourse(3)

	students := []yaac_shared.Student{
		{
			LName: "Maximus",
			FName: "Musterus",
		},
	}

	if !reflect.DeepEqual(studentsOfTis22, students) {
		t.Error("Failed to test InsertAttendance")
	} else {
		t.Log("Insert Attendance tested successfully")
	}
}

func TestInsertStudentConstraints(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	// Testing constraint violation
	longName := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."

	err := databaseinst.InsertStudent(longName, "Musterus", true, 3)
	if err == nil {
		t.Error("Failed to test lName constraint of InsertStudent")
	}

	err = databaseinst.InsertStudent("Maximus", longName, true, 3)
	if err == nil {
		t.Error("Failed to test fName constraint of InsertStudent")
	}
}

func TestInsertAttendanceList(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	// InsertAttendanceList already run in setupDatabaseHelper thus we only need to check if the entries are present
	attendanceLists, _ := databaseinst.GetAllAttendanceLists()

	lists := []yaac_shared.AttendanceList{
		{
			Id:           1,
			TimeReceived: testTime.Format("2006-01-02 15:04:05"),
			CourseId:     1,
			List:         testByteArray,
		},
		{
			Id:           2,
			TimeReceived: testTime.Format("2006-01-02 15:04:05"),
			CourseId:     2,
			List:         testByteArray,
		},
	}
	if !reflect.DeepEqual(attendanceLists, lists) {
		t.Error("Failed to test InsertAttendanceList")
	} else {
		t.Log("InsertAttendanceList tested successfully")
	}
}

/*func TestInsertAttendanceListWithAlreadyExistingEarlierList(t *testing.T) {
	databaseinst, _ := setupDatabaseHelper()
	defer clearDatabaseHelper(databaseinst)

	// timestamp later than testTime but on the same day
	var newTestTime = time.Date(2000, 01, 01, 2, 5, 1, 1, location)

	// inserting a attendance list where an earlier one already exists
	err := databaseinst.InsertAttendanceList(newTestTime, 1, testByteArray)
	if err != nil {
		t.Errorf("Failed to test InsertAttendanceListWithAlreadyExistingEarlierList with error: %q", err)
	}

	attendanceLists, _ := databaseinst.GetAllAttendanceLists()

	lists := []yaac_shared.AttendanceList{
		{
			Id:           1,
			TimeReceived: newTestTime.Format("2006-01-02 15:04:05"),
			CourseId:     1,
			List:         testByteArray,
		},
		{
			Id:           2,
			TimeReceived: testTime.Format("2006-01-02 15:04:05"),
			CourseId:     2,
			List:         testByteArray,
		},
	}
	if !reflect.DeepEqual(attendanceLists, lists) {
		t.Error("Failed to test InsertAttendanceList")
	} else {
		t.Log("InsertAttendanceList tested successfully")
	}
}*/

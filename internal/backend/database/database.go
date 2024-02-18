package yaac_backend_database

import (
	"database/sql"
	"errors"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

// CreateDatabase creates a database in ./data
func (item *BackendDatabase) CreateDatabase() {
	err := os.MkdirAll(item.path, 0755)
	if err != nil {
		log.Fatal("Could not create database folder: ", err)
	}

	_, err = os.Create(item.path + item.dbName)
	if err != nil {
		log.Fatal("Could not create database file: ", err)
	}

	db, err := sql.Open("sqlite3", item.path+item.dbName)
	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	// close database connection when finished
	defer log.Println("Database created")
	defer db.Close()

	// create Student table
	_, err = db.Exec("CREATE TABLE Student (StudentId INTEGER PRIMARY KEY, LName VARCHAR(50) NOT NULL, FName VARCHAR(50) NOT NULL, Course VARCHAR(12) NOT NULL, StatusOfMatriculation BOOLEAN NOT NULL Default True);")
	if err != nil {
		log.Fatal("Could not create student table on database: ", err)
	}

	// create Attendance table
	_, err = db.Exec("CREATE TABLE Attendance (StudentId INTEGER NOT NULL, DayOfAttendance TEXT NOT NULL, Attending BOOLEAN DEFAULT False NOT NULL, PRIMARY KEY (DayOfAttendance, StudentId) , FOREIGN KEY (StudentId) REFERENCES Student (StudentId) ON DELETE CASCADE);")
	if err != nil {
		log.Fatal("Could not create attendance table on database: ", err)
	}

	// create AttendanceList table
	_, err = db.Exec("CREATE TABLE AttendanceList (ListId INTEGER PRIMARY KEY, TimeRecieved TEXT,Course VARCHAR(12) NOT NULL, List BLOB NOT NULL);")
	if err != nil {
		log.Fatal("Could not create attendance list table on database: ", err)
	}
}

// ConnectDatabase creates a connection to the database at given path
func (item *BackendDatabase) ConnectDatabase() {
	db, err := sql.Open("sqlite3", item.path+item.dbName)
	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	item.database = db
}

// DisconnectDatabase closes the connection to the database
func (item *BackendDatabase) DisconnectDatabase() {
	err := item.database.Close()
	if err != nil {
		log.Println("Failed to close connection to database")
	}
}

// InsertAttendance inserts attendance
func (item *BackendDatabase) InsertAttendance(studentId int, timeOfAttendance time.Time, attending bool) error {
	// convert Time to YYYY-MM-DD
	date := timeOfAttendance.Format("2006-01-02")

	// use prepared statement for faster execution and to prevent sql injection attacks
	stmt, err := item.database.Prepare("INSERT INTO Attendance (StudentId, DayOfAttendance, Attending) VALUES (?, ?, ?);")
	if err != nil {
		log.Println("Could not create database prepared statement ", err)
		return err
	}

	defer stmt.Close()

	// execute prepared statement
	_, err = stmt.Exec(studentId, date, attending)
	if err != nil {
		log.Println("Could not add attendance: ", err)
		return err
	}

	return nil
}

// InsertCurrentAttendance inserts attendance for the current day
func (item *BackendDatabase) InsertCurrentAttendance(studentId int, attending bool) error {
	if err := item.InsertAttendance(studentId, time.Now(), attending); err != nil {
		log.Println("Could not add attendance: ", err)
		return err
	}

	return nil
}

func (item *BackendDatabase) UpdateAttendance(studentId int, timeOfAttendance time.Time, attending bool) error {
	// convert Time to YYYY-MM-DD
	date := timeOfAttendance.Format("2006-01-02")

	// prepared statement
	stmt, err := item.database.Prepare("UPDATE Attendance SET Attending = ? WHERE StudentId = ? AND DayOfAttendance = ?;")
	if err != nil {
		log.Fatal("Could not create database prepared statement ", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(attending, studentId, date)
	if err != nil {
		log.Println("Could not update attendance: ", err)
		return err
	}

	return nil
}

// InsertStudent inserts student to the database
func (item *BackendDatabase) InsertStudent(lName string, fName string, statusOfMatriculation bool, course string) error {
	// check constraints matching
	if len(lName) > 50 {
		log.Println("Maximum length for last name is 50")
		return errors.New("input does not match constraints")
	}

	if len(fName) > 50 {
		log.Println("Maximum length for first name is 50")
		return errors.New("input does not match constraints")
	}

	if len(course) > 12 {
		log.Println("Maximum length for course is 12")
		return errors.New("input does not match constraints")
	}

	// use prepared statement for faster execution and to prevent sql injection attacks
	stmt, err := item.database.Prepare("INSERT INTO Student (LName, FName, Course, StatusOfMatriculation) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.Fatal("Could not create database prepared statement ", err)
	}

	defer stmt.Close()

	// execute prepared statement
	_, err = stmt.Exec(lName, fName, course, statusOfMatriculation)
	if err != nil {
		log.Println("Could not add student: ", err)
		return err
	}

	return nil
}

// InsertAttendanceList inserts list with time=now by its relative link to the content root
func (item *BackendDatabase) InsertAttendanceList(dateTime time.Time, course string, list []byte) error {
	// check constraint matching
	if len(course) > 12 {
		log.Println("Maximum length for course is 12")
		return errors.New("input does not match constraints")
	}

	// check for already existing attendance lists for the current day and override if they are older
	// prepared statement to check if an attendance list fot the same day already exists
	searchList, err := item.database.Prepare("SELECT TimeRecieved FROM AttendanceList WHERE Course = ? AND DATE(TimeRecieved) = DATE(?)")
	if err != nil {
		log.Fatal("Could not create database prepared statement", err)
	}
	defer searchList.Close()

	result := searchList.QueryRow(course, dateTime)

	var timeRecievedStr string

	err = result.Scan(&timeRecievedStr)
	if err == nil {
		// found existing attendance lists
		timeExisting, err := time.Parse("2006-01-02 15:04:05", timeRecievedStr)
		if err != nil {
			log.Println("Could not pare time from sqlite to go")
			return err
		}

		// check if existing list needs to be overwritten and overwrite if true
		if timeExisting.Before(dateTime) {
			deleteList, err := item.database.Prepare("DELETE FROM AttendanceList WHERE TimeRecieved = ? AND Course = ?")
			if err != nil {
				log.Fatal("Could not create database prepared statement", err)
			}
			defer deleteList.Close()

			// execute delete statement
			if _, err = deleteList.Exec(timeExisting, course); err != nil {
				log.Println("Could not delete existing attendance list to override")
				return err
			}

			err = item.insertAttendanceListHelper(dateTime, course, list)
			if err != nil {
				return err
			}
		}
	} else {
		// no existing attendance lists found
		err := item.insertAttendanceListHelper(dateTime, course, list)
		if err != nil {
			return err
		}
	}
	return nil
}

func (item *BackendDatabase) insertAttendanceListHelper(dateTime time.Time, course string, list []byte) error {
	// prepared statement to insert new attendance list
	insertList, err := item.database.Prepare("INSERT INTO AttendanceList (TimeRecieved, Course, List) VALUES (?, ?, ?);")
	if err != nil {
		log.Fatal("Could not create database prepared statement", err)
	}
	defer insertList.Close()

	// execute prepared statement
	_, err = insertList.Exec(dateTime, course, list)
	if err != nil {
		log.Println("Could not add attendance list ", err)
		return err
	}

	return nil
}

func (item *BackendDatabase) InsertCurrentAttendanceList(course string, list []byte) error {
	if err := item.InsertAttendanceList(time.Now(), course, list); err != nil {
		log.Println("Could not add attendance list: ", err)
		return err
	}

	return nil
}

// GetStudentFullNameById takes a StudentId and returnes the students full name in one column
func (item *BackendDatabase) GetStudentFullNameById(studentId int) (string, error) {
	if studentId < 0 {
		log.Println("StudentId must be positive integer")
		return "", errors.New("input does not match constraints")
	}

	// prepare statement
	stmt, err := item.database.Prepare("SELECT FName || ' ' || LName FROM Student WHERE StudentId = ?")
	if err != nil {
		log.Println("Could not create database prepared statement ", err)
		return "", err
	}

	result := stmt.QueryRow(studentId)

	var studentName string

	if err := result.Scan(&studentName); err != nil {
		log.Println("Could not get student name by id ", err)
		return "", errors.New("could not get student name by id")
	}

	return studentName, nil
}

func (item *BackendDatabase) GetLatestListDatePerCourse(course string) (string, error) {
	if len(course) > 12 {
		log.Println("Maximum length for course is 12")
		return "", errors.New("input does not match constraints")
	}

	// prepare statement
	stmt, err := item.database.Prepare("SELECT DATE(TimeRecieved) FROM AttendanceList WHERE Course = ? ORDER BY TimeRecieved DESC LIMIT 1")
	if err != nil {
		log.Println("Could not create database prepared statement ", err)
		return "", err
	}

	result := stmt.QueryRow(course)

	var date string

	if err := result.Scan(&date); err != nil {
		log.Println("Could not get date of latest attendance list for course: "+course+" ", err)
		return "", errors.New("could not get date of latest attendance list for course:" + course)
	}

	return date, nil
}

func (item *BackendDatabase) GetAllCourses() ([]string, error) {
	// prepare statement
	stmt, err := item.database.Prepare("SELECT DISTINCT Course FROM AttendanceList")
	if err != nil {
		log.Println("Could not create database prepared statement ", err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			log.Println("Could not get all courses")
			return nil, errors.New("could not get all courses")
		}

		result = append(result, row)
	}

	return result, nil
}

func (item *BackendDatabase) GetAllStudentsPerCourse(course string) ([]yaac_shared.Student, error) {
	if len(course) > 12 {
		log.Println("Maximum length for course is 12")
		return nil, errors.New("input does not match constraints")
	}

	// prepare statement
	stmt, err := item.database.Prepare("SELECT DISTINCT FName, LName FROM Student WHERE Course = ?")
	if err != nil {
		log.Println("Could not create database prepared statement ", err)
		return nil, err
	}

	rows, err := stmt.Query(course)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []yaac_shared.Student
	for rows.Next() {
		var row yaac_shared.Student
		if err := rows.Scan(&row.FName, &row.LName); err != nil {
			log.Println("Could not get all students per course")
			return nil, errors.New("could not get all students per course")
		}

		result = append(result, row)
	}

	return result, nil
}

func (item *BackendDatabase) GetAllAttendanceWithStudentName() ([]yaac_shared.Attendance, error) {
	// prepare statement
	stmt, err := item.database.Prepare("SELECT DISTINCT s.FName, s.LName, a.Attending FROM Attendance AS a LEFT OUTER JOIN Student AS s ON a.StudentId = s.StudentId")
	if err != nil {
		log.Println("Could not create database prepared statement ", err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []yaac_shared.Attendance
	for rows.Next() {
		var row yaac_shared.Attendance
		if err := rows.Scan(&row.Student.FName, &row.Student.LName, &row.Attending); err != nil {
			log.Println("Could not get all students per course")
			return nil, errors.New("could not get all students per course")
		}

		result = append(result, row)
	}

	return result, nil
}

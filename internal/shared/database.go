package yaac_shared

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

// CreateDatabase creates a database in ./data
func CreateDatabase() {
	err := os.MkdirAll("./data", 0755)
	if err != nil {
		log.Fatal("Could not create database folder: ", err)
	}

	_, err = os.Create("./data/data.db")
	if err != nil {
		log.Fatal("Could not create database file: ", err)
	}

	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	// close database connection when finished
	defer log.Println("Database created")
	defer db.Close()

	// create Student table
	_, err = db.Exec("CREATE TABLE Student (StudentId INTEGER PRIMARY KEY, LName VARCHAR(50) NOT NULL, FName VARCHAR(50) NOT NULL, Course VARCHAR(12) NOT NULL, StatusOfMatriculation BOOLEAN NOT NULL Default=True);")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}

	// create Attendance table
	_, err = db.Exec("CREATE TABLE Attendance (StudentId INTEGER NOT NULL, DayOfAttendance TEXT NOT NULL DEFAULT=DATE(), Attending BOOLEAN DEFAULT=False NOT NULL, PRIMARY KEY (DayOfAttendance, StudentId) , FOREIGN KEY (StudentId) REFERENCES Student (StudentId) ON DELETE CASCADE);")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}

	// create AttendanceList table
	_, err = db.Exec("CREATE TABLE AttendanceList (ListId INTEGER PRIMARY KEY, TimeRecieved TEXT DEFAULT=DATETIME(),Course VARCHAR(12) NOT NULL, ListPath TEXT NOT NULL;")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}
}

// ConnectDatabase creates a connection to the database at ./data/data.db
func ConnectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatal("Connecting to database failed: ", err)
	}

	return db
}

// DisconnectDatabase closes the connection to the database
func DisconnectDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Failed to close connection to database")
	}
}

// inserts attendance for today
func insertAttendance(db *sql.DB, attending bool) error {
	// use prepared statement for faster execution and to prevent sql injection attacks
	stmt, err := db.Prepare("INSERT INTO Attendance (DayOfAttendance, Attending) VALUES (DATE(), ?);")
	if err != nil {
		log.Fatal("Could not create database prepared statement ", err)
	}

	defer stmt.Close()

	// execute prepared statement
	_, err = stmt.Exec(attending)
	if err != nil {
		log.Println("Could not add attendance: ", err)
		return err
	}

	return nil
}

// InsertStudent inserts student to the database
func InsertStudent(db *sql.DB, lName string, fName string, statusOfMatriculation bool, course string) error {
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
	stmt, err := db.Prepare("INSERT INTO Student (LName, FName, Course, StatusOfMatriculation) VALUES (?, ?, ?, ?);")
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
func InsertAttendanceList(db *sql.DB, course string, listPath string) error {
	// check constraint matching
	if len(course) > 12 {
		log.Println("Maximum length for course is 12")
		return errors.New("input does not match constraints")
	}

	// prepared statement
	stmt, err := db.Prepare("INSERT INTO AttendanceList (TimeRecieved, Course, ListPath) VALUES (DATETIME(), ?, ?);")
	if err != nil {
		log.Fatal("Could not create database prepared statement", err)
	}

	defer stmt.Close()

	// execute prepared statement
	_, err = stmt.Exec(course, listPath)
	if err != nil {
		log.Println("Could not add attendance list ", err)
		return err
	}

	return nil
}

// GetStudentFullNameById takes a StudentId and returnes the students full name in one column
func GetStudentFullNameById(db *sql.DB, studentId int) (string, error) {
	if studentId < 0 {
		log.Println("StudentId must be positive integer")
		return "", errors.New("input does not match constraints")
	}

	// prepare statement
	stmt, err := db.Prepare("SELECT FName || ' ' || LName FROM Student WHERE StudentId = ?")
	if err != nil {
		log.Fatal("Could not create database prepared statement ", err)
	}

	result := stmt.QueryRow(studentId)

	studentName := ""

	if err := result.Scan(&studentName); err != nil {
		log.Println("Could not get student name by id ", err)
		return "", errors.New("could not get student name by id")
	}

	return studentName, nil
}

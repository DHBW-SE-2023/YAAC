package yaac_shared

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func createDatabase() {
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

	// close database in the end
	defer db.Close()

	// create Student table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Student (StudentId INTEGER PRIMARY KEY, SName VARCHAR(30) NOT NULL, SFirstName VARCHAR(30) NOT NULL, Course CHAR(3) NOT NULL, StatusOfMatriculation BOOLEAN NOT NULL Default=True);")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}

	// create Lecture table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Lecture (LectureId INTEGER PRIMARY KEY, LName VARCHAR(30) NOT NULL, Room CHAR(4), Online BOOLEAN NOT NULL Default=False, LecturerName VARCHAR(50));")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}

	// create Attendance table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Attendance (LectureId INTEGER NOT NULL, StudentId INTEGER NOT NULL, DATE TEXT NOT NULL DEFAULT=DATE(),PRIMARY KEY (LectureId, StudentId) , FOREIGN KEY (StudentId) REFERENCES Student (StudentId) ON DELETE CASCADE), FOREIGN KEY (LectureId) REFERENCES Lecture (LectureId) ON DELETE CASCADE);")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}

	// create AttendanceList table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS AttendanceList (ListId INTEGER PRIMARY KEY, Date TEXT DEFAULT=DATE(),Course CHAR(3), ListImage BLOB NOT NULL;")
	if err != nil {
		log.Fatal("Could not create table on database: ", err)
	}
}

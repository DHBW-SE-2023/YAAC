package yaac_demon

import (
	"log"
	"time"

	database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gocv.io/x/gocv"
)

func StartDemon(mvvm shared.MVVM, duration time.Duration) {
	// Run forever
	for {
		time.Sleep(duration)

		newMails, err := mvvm.GetMailsToday()
		log.Println("New mails: ", len(newMails))
		if err != nil {
			log.Println("ERROR: Could not get mails for today: ", err)
			continue
		}

		for _, mail := range newMails {
			list, err := TableToAttendanceList(mvvm, mail)
			if err != nil {
				log.Println("ERROR: Could not process image from mail received at ", mail.ReceivedAt)
				continue
			}
			_, err = mvvm.InsertList(list)
			if err != nil {
				log.Printf("ERROR: Could not add list for mail received at %v: %v\n", mail.ReceivedAt, err)
				continue
			}

			mvvm.NotifyNewList(list)
		}
	}
}

func TableToAttendanceList(mvvm shared.MVVM, mail shared.MailData) (shared.AttendanceList, error) {
	table, err := mvvm.NewTable(mail.Image)
	if err != nil {
		return shared.AttendanceList{}, err
	}

	course, err := mvvm.CourseByName(table.Course)
	if err != nil {
		return shared.AttendanceList{}, err
	}

	img, err := gocv.IMEncode(".png", table.Image)
	if err != nil {
		return shared.AttendanceList{}, err
	}

	list := shared.AttendanceList{
		CourseID:   course.ID,
		ReceivedAt: mail.ReceivedAt,
		Image:      img.GetBytes(),
	}

	for _, row := range table.Rows {
		if row.FullName == "" || row.FirstName == "" || row.LastName == "" {
			continue
		}

		// var student shared.Student
		// students, _ := mvvm.CourseStudents(shared.Course{Model: gorm.Model{ID: list.ID}})
		// for _, element := range students {
		// 	if element.LastName == strings.TrimSpace(row.LastName) {
		// 		student = element
		// 	}
		// }

		students, err := mvvm.Students(shared.Student{LastName: row.LastName})
		if err != nil {
			return shared.AttendanceList{}, err
		}

		var student shared.Student

		if len(students) == 0 {
			student, err = mvvm.InsertStudent(shared.Student{
				FirstName: row.FirstName,
				LastName:  row.LastName,
				CourseID:  course.ID,
			})

			if err != nil {
				continue
			}
		} else if len(students) != 1 { // If there are more students with the same name, we don't know what to do
			continue
		} else {
			student = students[0]
		}

		attendance := shared.Attendance{
			StudentID:    student.ID,
			IsAttending:  row.Valid,
			NameROI:      database.Rectangle(row.NameROI),
			SignatureROI: database.Rectangle(row.SignatureROI),
			TotalROI:     database.Rectangle(row.TotalROI),
		}

		list.Attendancies = append(list.Attendancies, attendance)
	}

	return list, nil
}

func UploadImage(mvvm shared.MVVM, img []byte, course *shared.Course) (*shared.AttendanceList, error) {
	list, err := TableToAttendanceList(mvvm, shared.MailData{Image: img, ReceivedAt: time.Now()})
	if err != nil {
		return nil, err
	}

	if course != nil {
		list.CourseID = course.ID
	}

	list, err = mvvm.InsertList(list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

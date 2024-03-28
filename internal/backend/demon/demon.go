package yaac_demon

import (
	"log"
	"time"

	database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func StartDemon(mvvm shared.MVVM, duration time.Duration) {
	// Run forever
	for {
		newMails, err := mvvm.GetMailsToday()
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

		time.Sleep(duration * time.Second)
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

	list := shared.AttendanceList{
		CourseID:   course.ID,
		ReceivedAt: mail.ReceivedAt,
		Image:      mail.Image,
	}

	for _, row := range table.Rows {
		if row.FullName == "" || row.FirstName == "" || row.LastName == "" {
			continue
		}

		students, err := mvvm.Students(shared.Student{CourseID: list.CourseID, FirstName: row.FirstName, LastName: row.LastName})
		if err != nil {
			return shared.AttendanceList{}, err
		}

		var student shared.Student

		if len(students) == 0 {
			student, err = mvvm.InsertStudent(shared.Student{
				FirstName: row.FirstName,
				LastName:  row.LastName,
			})

			if err != nil {
				continue
			}
		} else if len(students) != 1 { // If there are more students with the same name, we don't know what to do
			continue
		}

		student = students[0]

		attendance := shared.Attendance{
			StudentID:    student.CourseID,
			IsAttending:  row.Valid,
			NameROI:      database.Rectangle(row.NameROI),
			SignatureROI: database.Rectangle(row.SignatureROI),
			TotalROI:     database.Rectangle(row.TotalROI),
		}

		list.Attendancies = append(list.Attendancies, attendance)
	}

	return list, nil
}

func UploadImage(mvvm shared.MVVM, img []byte) (*shared.AttendanceList, error) {
	list, err := TableToAttendanceList(mvvm, shared.MailData{Image: img, ReceivedAt: time.Now()})
	if err != nil {
		return nil, err
	}

	list, err = mvvm.InsertList(list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

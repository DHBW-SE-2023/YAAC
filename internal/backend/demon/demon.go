package yaac_backend_demon

import (
	"fmt"
	"log"
	"time"

	database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gocv.io/x/gocv"
)

// Trigger a single runthrough of the demon. This runs independet from StartDemon
func SingleDemonRunthrough(mvvm shared.MVVM) {
	newMails, err := mvvm.GetMailsToday()
	if err != nil {
		log.Println("ERROR: Could not get mails for today: ", err)
		mvvm.NotifyError("demon", fmt.Errorf("e-mails konnten nicht abgerufen werden"))
		return
	}
	log.Println("New mails: ", len(newMails))

	for _, mail := range newMails {
		list, err := TableToAttendanceList(mvvm, mail)
		if err != nil {
			log.Printf("ERROR: Could not process image for mail received at %v: %v\n", mail.ReceivedAt, err)
			mvvm.NotifyError("demon", fmt.Errorf("bild konnte nicht verarbeitet werden"))
			continue
		}
		_, err = mvvm.InsertList(list)
		if err != nil {
			log.Printf("ERROR: Could not add list for mail received at %v: %v\n", mail.ReceivedAt, err)
			mvvm.NotifyError("demon", fmt.Errorf("neue anwesenheitsliste für e-mail vom %v konnte nicht hinzugefügt werden", mail.ReceivedAt))
			continue
		}

		mvvm.NotifyNewList(list)
	}
}

// Start the demon. This runs forever and calls a SingleDemonRunthrough every duration
func StartDemon(mvvm shared.MVVM, duration time.Duration) {
	// Run forever
	for {
		time.Sleep(duration)
		SingleDemonRunthrough(mvvm)
	}
}

// Convert MailData to an AttendanceList.
// If the students extracted from the list do not exist, they are created.
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

		students, err := mvvm.Students(shared.Student{LastName: row.LastName})
		if err != nil {
			return shared.AttendanceList{}, err
		}

		var student shared.Student

		if len(students) == 0 {
			student, err = mvvm.InsertStudent(shared.Student{
				FirstName: row.FirstName,
				LastName:  row.LastName,
				FullName:  row.FullName,
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

// Upload an image. This is basically SingleDemonRunthrough but it takes in an image and a course instead of implicitly reading the mails.
// course may be nil.
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

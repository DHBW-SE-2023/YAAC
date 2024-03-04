package yaac_demon

import (
	"log"
	"time"

	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func StartDemon(mvvm shared.MVVM, duration time.Duration) {
	// Run forever
	for {
		newMails := mvvm.MailsRefresh()

		for _, mail := range newMails {
			table, err := mvvm.ValidateTable(mail.Image)
			if err != nil {
				log.Fatalf("Could not process mail received at %v", mail.ReceivedAt)
				continue
			}

			list := shared.AttendanceList{
				ReceivedAt: mail.ReceivedAt,
				Image:      mail.Image,
			}

			for _, row := range table.Rows {
				students, err := mvvm.Students(shared.Student{FirstName: row.FirstName, LastName: row.LastName})
				if err != nil {
					continue
				}

				var student shared.Student

				// TODO: Is this correct: No student with the name is present. Consider adding them
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
					StudentID:   student.CourseID,
					IsAttending: row.Valid,
				}

				list.Attendancies = append(list.Attendancies, attendance)
			}

			_, err = mvvm.InsertList(list)
			if err != nil {
				log.Fatalf("Could not add list for mail received at %v: %v", mail.ReceivedAt, err)
				continue
			}

			// TODO: Notify frontend, that a new attendance list was created
		}

		time.Sleep(duration * time.Second)
	}
}
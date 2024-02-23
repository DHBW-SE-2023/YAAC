package yaac_shared

import "time"

type MVVM interface {
	MailFormUpdated(data EmailData)
	ValidateTable(img []byte)
	InsertList(list AttendanceList) AttendanceList
	UpdateList(list AttendanceList) AttendanceList
	LatestList(course Course, date time.Time) AttendanceList
	Courses() []Course
}

// models/meeting.go
package models

import (
	"time"
)

type ZoomMeeting struct {
	ID          uint      `gorm:"primaryKey"`         // ID unik untuk masing-masing record
	Topic       string    `gorm:"type:varchar(255);not null"` // Judul atau topik meeting
	StartTime   time.Time `gorm:"not null"`           // Waktu mulai meeting
	Duration    int       `gorm:"not null"`           // Durasi meeting dalam menit
	Password    string    `gorm:"type:varchar(255)"`  // Password meeting
	MeetingID  string   `gorm:"type:varchar(255);unique;not null"` // ID dari Zoom meeting (dari response API Zoom)
	ScheduledAt time.Time `gorm:"autoCreateTime"`     // Waktu pencatatan data meeting
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`     // Waktu update data meeting
	Agenda string `gorm:"type:varchar(255)"`
	Timezone string `gorm:"type:varchar(255)"`
	ScheduleFor string `gorm:"type:varchar(255)"`
	URL string `gorm:"type:text"`
}
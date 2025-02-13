// dtos/meetingDTO.go
package dtos

import "time"

type CreateZoomMeetingRequest struct {
	Agenda      string    `json:"agenda"`      // Agenda atau deskripsi meeting
	Topic       string    `json:"topic"`       // Judul/topik meeting
	StartTime   time.Time `json:"start_time"`  // Waktu mulai meeting (format: RFC3339)
	Duration    int       `json:"duration"`    // Durasi meeting dalam menit
	Password    string    `json:"password"`    // Password untuk meeting
	ScheduleFor string    `json:"schedule_for"`// Email orang yang menjadwalkan meeting
	Timezone    string    `json:"timezone"`    // Timezone untuk meeting (contoh: "Asia/Jakarta")
}

type UpdateMeetingDTO struct {
	Agenda      string    `json:"agenda"`      // Agenda atau deskripsi meeting
	Topic       string    `json:"topic"`       // Judul/topik meeting
	StartTime   time.Time `json:"start_time"`  // Waktu mulai meeting (format: RFC3339)
	Duration    int       `json:"duration"`    // Durasi meeting dalam menit
	Password    string    `json:"password"`    // Password untuk meeting
	ScheduleFor string    `json:"schedule_for"`// Email orang yang menjadwalkan meeting
	Timezone    string    `json:"timezone"`    // Timezone untuk meeting (contoh: "Asia/Jakarta")
}

type GetZoomMeeting struct {
	ID int `json:"id"` //
	Agenda          string `json:"agenda"`
	Duration        int    `json:"duration"`
	Password        string `json:"password"`
	ScheduleFor     string `json:"schedule_for"`
	StartTime       string `json:"start_time"`
	Topic           string `json:"topic"`
	Timezone        string `json:"timezone"`
	URL string `json:"url"`
}


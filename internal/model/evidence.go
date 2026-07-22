package model

import "time"

type Evidence struct {
	EvidenceID  uint      `json:"evidence_id"`
	ScheduleID  uint      `json:"schedule_id"`
	ObjectKey   string    `json:"object_key"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

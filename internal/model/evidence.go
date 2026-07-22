package model

type Evidence struct {
	EvidenceID uint `db:"evidence_id"`
	ScheduleID uint `db:"schedule_id"`
	FileID     uint `db:"file_id"`
}

package dto

type CreateEvidenceRequest struct {
	ScheduleID  uint   `json:"schedule_id"`
	ObjectKey   string `json:"object_key"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

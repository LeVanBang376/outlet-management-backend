package model

import "time"

type File struct {
	FileID      uint      `db:"file_id"`
	ObjectKey   string    `db:"object_key"`
	FileName    string    `db:"file_name"`
	ContentType string    `db:"content_type"`
	Size        int64     `db:"size"`
	CreatedAt   time.Time `db:"created_at"`
}

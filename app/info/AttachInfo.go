package info

import (
	"time"
)

type Attach struct {
	AttachId     string    `db:"id"`
	NoteId       string    `db:"note_id"`
	UploadUserId string    `db:"upload_user_id"`
	Name         string    `db:"name"`
	Title        string    `db:"title"`
	Size         int64     `db:"size"`
	Type         string    `db:"type"`
	Path         string    `db:"path"`
	CreatedTime  time.Time `db:"created_time"`
}

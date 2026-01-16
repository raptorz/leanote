package info

import (
	"time"
)

type NoteImage struct {
	NoteImageId string    `db:"id"`
	NoteId      string    `db:"note_id"`
	UserId      string    `db:"user_id"`
	FileId      string    `db:"file_id"`
	Path        string    `db:"path"`
	CreatedTime time.Time `db:"created_time"`
}

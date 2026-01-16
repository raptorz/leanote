package info

import (
	"time"
)

type File struct {
	FileId         string    `db:"id"`
	UserId         string    `db:"user_id"`
	Name           string    `db:"name"`
	Title          string    `db:"title"`
	Size           int64     `db:"size"`
	Path           string    `db:"path"`
	MimeType       string    `db:"mime_type"`
	CreatedTime    time.Time `db:"created_time"`
	AlbumId        string    `db:"album_id"`
	IsDefaultAlbum bool      `db:"is_default_album"`
	FromFileId     string    `db:"from_file_id"`
}

package info

import (
	"time"
)

type Album struct {
	AlbumId     string    `db:"id"`
	UserId      string    `db:"user_id"`
	Title       string    `db:"title"`
	Type        int       `db:"type"`
	Seq         int       `db:"seq"`
	CreatedTime time.Time `db:"created_time"`
	IsDeleted   bool      `db:"is_deleted"`
}

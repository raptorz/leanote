package info

import (
	"time"
)

type Tag struct {
	UserId string   `db:"id"`
	Tags   []string `db:"tags"`
}

type NoteTag struct {
	TagId       string    `db:"id"`
	UserId      string    `db:"user_id"`
	Tag         string    `db:"tag"`
	Usn         int       `db:"usn"`
	Count       int       `db:"count"`
	CreatedTime time.Time `db:"created_time"`
	UpdatedTime time.Time `db:"updated_time"`
	IsDeleted   bool      `db:"is_deleted"`
}

type TagCount struct {
	TagCountId string `db:"id"`
	UserId     string `db:"user_id"`
	Tag        string `db:"tag"`
	IsBlog     bool   `db:"is_blog"`
	Count      int    `db:"count"`
}

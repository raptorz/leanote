package info

import (
	"time"
)

type Group struct {
	GroupId     string    `db:"id"`
	UserId      string    `db:"user_id"`
	Title       string    `db:"title"`
	UserCount   int       `db:"user_count"`
	CreatedTime time.Time `db:"created_time"`
	IsDeleted   bool      `db:"is_deleted"`

	Users []User `db:"-"`
}

type GroupUser struct {
	GroupUserId string    `db:"id"`
	GroupId     string    `db:"group_id"`
	UserId      string    `db:"user_id"`
	CreatedTime time.Time `db:"created_time"`
}

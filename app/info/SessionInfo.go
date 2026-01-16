package info

import (
	"time"
)

type Session struct {
	Id         string    `db:"id"`
	SessionId  string    `db:"session_id"`

	LoginTimes int       `db:"login_times"`
	Captcha    string    `db:"captcha"`

	UserId string `db:"user_id"`

	CreatedTime time.Time `db:"created_time"`
	UpdatedTime time.Time `db:"updated_time"`
}

type DBSession struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Data      string    `db:"data"`
	ExpiresAt time.Time `db:"expires_at"`
}

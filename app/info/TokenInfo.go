package info

import (
	"time"
)

const (
	TokenPwd = iota
	TokenActiveEmail
	TokenUpdateEmail
)

const (
	PwdOverHours         = 2.0
	ActiveEmailOverHours = 48.0
	UpdateEmailOverHours = 2.0
)

type Token struct {
	TokenId     string    `db:"id"`
	UserId      string    `db:"user_id"`
	Email       string    `db:"email"`
	Token       string    `db:"token"`
	Type        int       `db:"type"`
	CreatedTime time.Time `db:"created_time"`
}

package info

import (
	"time"
)

type EmailLog struct {
	LogId      string    `db:"id"`
	ToEmail    string    `db:"to_email"`
	Subject    string    `db:"subject"`
	Content    string    `db:"content"`
	CreatedTime time.Time `db:"created_time"`
}

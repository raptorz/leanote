package info

import (
	"time"
)

type Report struct {
	ReportId  string    `db:"id"`
	UserId    string    `db:"user_id"`
	TargetId  string    `db:"target_id"`
	TargetType string   `db:"target_type"`
	Reason    string    `db:"reason"`
	CreatedTime time.Time `db:"created_time"`
}

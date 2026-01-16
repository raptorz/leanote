package info

import (
	"time"
)

type Config struct {
	ConfigId    string              `db:"id"`
	UserId      string              `db:"user_id"`
	Key         string              `db:"key"`
	ValueStr    string              `db:"value_str"`
	ValueArr    []string            `db:"value_arr"`
	ValueMap    map[string]string   `db:"value_map"`
	ValueArrMap []map[string]string `db:"value_arr_map"`
	IsArr       bool                `db:"is_arr"`
	IsMap       bool                `db:"is_map"`
	IsArrMap    bool                `db:"is_arr_map"`

	UpdatedTime time.Time `db:"updated_time"`
}

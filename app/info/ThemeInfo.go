package info

import (
	"time"
)

type Theme struct {
	ThemeId     string                 `db:"id"`
	Title       string                 `db:"title"`
	ThemeUrl    string                 `db:"theme_url"`
	CssUrl      string                 `db:"css_url"`
	PreviewImg  string                 `db:"preview_img"`
	Description string                 `db:"description"`
	Info        map[string]interface{} `db:"info"`
	CreatedTime time.Time              `db:"created_time"`
	IsDeleted   bool                   `db:"is_deleted"`
	IsDefault   bool                   `db:"is_default"`
	Path        string                 `db:"path"`
	Name        string                 `db:"name"`
	Author      string                 `db:"author"`
	AuthorUrl   string                 `db:"author_url"`
	Version     string                 `db:"version"`
}

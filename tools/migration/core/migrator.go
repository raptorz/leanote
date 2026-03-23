package core

import (
	"github.com/leanote/leanote/app/db"
)

// Migrator 迁移器接口
type Migrator interface {
	Migrate() (*db.MigrationStats, error)
}

// 获取表迁移顺序
func GetTableMigrationOrder() []string {
	return []string{
		"users",                  // 第一步，无依赖
		"notebooks",              // 依赖users
		"notes",                  // 依赖notebooks和users
		"note_contents",          // 依赖notes
		"note_content_histories", // 依赖notes
		"tags",                   // 依赖users
		"note_tags",              // 依赖notes和tags
		"tag_counts",             // 依赖tags
		"attachs",                // 依赖notes和users
		"share_notebooks",        // 依赖notebooks和users
		"share_notes",            // 依赖notes和users
		"groups",                 // 依赖users
		"group_users",            // 依赖groups和users
		"blogs",                  // 依赖notes和users
		"user_blogs",             // 依赖users和blogs
		"blog_singles",           // 依赖users
		"themes",                 // 无依赖
		"files",                  // 依赖users和albums
		"albums",                 // 依赖users
		"note_images",            // 依赖notes和files
		"tokens",                 // 依赖users
		"suggestions",            // 依赖users
		"configs",                // 依赖users
		"email_logs",             // 无依赖
		"sessions",               // 依赖users
		"blog_likes",             // 依赖blogs和users
		"blog_comments",          // 依赖blogs和users
		"reports",                // 依赖users
	}
}

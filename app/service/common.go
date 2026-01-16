package service

import (
	"database/sql"
	"github.com/leanote/leanote/app/db"
	"strings"
)

// service 通用方法

// 分页, 排序处理
func parsePageAndSort(pageNumber, pageSize int, sortField string, isAsc bool) (skipNum int, sortFieldR string) {
	skipNum = (pageNumber - 1) * pageSize
	if skipNum < 0 {
		skipNum = 0
	}

	// 默认排序字段
	sortFieldR = "updated_time"
	if sortField != "" {
		// 转换字段名到数据库列名
		switch sortField {
		case "UpdatedTime":
			sortFieldR = "updated_time"
		case "CreatedTime":
			sortFieldR = "created_time"
		case "Title":
			sortFieldR = "title"
		case "PublicTime":
			sortFieldR = "public_time"
		default:
			sortFieldR = sortField
		}
	}

	if !isAsc {
		sortFieldR = sortFieldR + " DESC"
	}
	return
}

// incrUsn 增加用户USN
func incrUsn(userId string) int {
	var usn int
	err := db.DB.QueryRow("SELECT usn FROM users WHERE id = $1", userId).Scan(&usn)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		return 0
	}
	usn += 1
	_, err = db.DB.Exec("UPDATE users SET usn = $1 WHERE id = $2", usn, userId)
	if err != nil {
		return 0
	}
	return usn
}

// GetUrTitle 获取URL标题（简化版）
func GetUrTitle(userId string, title string, types string, id string) string {
	// 简化实现：将标题转换为小写并用连字符替换空格
	// TODO: 实现完整的URL标题生成逻辑
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

package service

import (
	"database/sql"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"time"
)

/*
每添加,更新note时, 都将tag添加到tags表中
*/
type TagService struct {
}

/*
func (this *TagService) GetTags(userId string) []string {
	tag := info.Tag{}
	db.Get(db.Tags, userId, &tag)
	LogJ(tag)
	return tag.Tags
}
*/

func (this *TagService) AddTagsI(userId string, tags interface{}) bool {
	if ts, ok2 := tags.([]string); ok2 {
		return this.AddTags(userId, ts)
	}
	return false
}
func (this *TagService) AddTags(userId string, tags []string) bool {
	// TODO: 需要实现PostgreSQL版本的标签添加
	// 暂时返回true，后续实现
	return true
}

//---------------------------
// v2
// 第二版标签, 单独一张表, 每一个tag一条记录

// 添加或更新标签, 先查下是否存在, 不存在则添加, 存在则更新
// 都要统计下tag的note数
// 什么时候调用? 笔记添加Tag, 删除Tag时
// 删除note时, 都可以调用
// 万能
func (this *TagService) AddOrUpdateTag(userId string, tag string) info.NoteTag {
	noteTag := info.NoteTag{}

	// 查询是否已存在
	query := "SELECT tag_id, user_id, tag, count, usn, created_time, updated_time, is_deleted FROM note_tags WHERE user_id = $1 AND tag = $2"
	err := db.DB.QueryRow(query, userId, tag).Scan(
		&noteTag.TagId, &noteTag.UserId, &noteTag.Tag, &noteTag.Count,
		&noteTag.Usn, &noteTag.CreatedTime, &noteTag.UpdatedTime, &noteTag.IsDeleted,
	)

	// 存在, 则更新之
	if err == nil && noteTag.TagId != "" {
		// 统计note数
		count := noteService.CountNoteByTag(userId, tag)
		noteTag.Count = count
		noteTag.UpdatedTime = time.Now()

		// 之前删除过的, 现在要添加回来了
		if noteTag.IsDeleted {
			Log("之前删除过的, 现在要添加回来了:  " + tag)
			noteTag.Usn = incrUsn(userId)
			noteTag.IsDeleted = false
		}

		// 更新标签
		updateQuery := "UPDATE note_tags SET count = $1, updated_time = $2, usn = $3, is_deleted = $4 WHERE tag_id = $5 AND user_id = $6"
		_, err := db.DB.Exec(updateQuery, noteTag.Count, noteTag.UpdatedTime, noteTag.Usn, noteTag.IsDeleted, noteTag.TagId, userId)
		if err != nil {
			Log(err.Error())
		}
		return noteTag
	}

	// 不存在, 则创建之
	if err == sql.ErrNoRows || noteTag.TagId == "" {
		noteTag.TagId = db.NewUUID()
		noteTag.Count = 1
		noteTag.Tag = tag
		noteTag.UserId = userId
		noteTag.CreatedTime = time.Now()
		noteTag.UpdatedTime = noteTag.CreatedTime
		noteTag.Usn = incrUsn(userId)
		noteTag.IsDeleted = false

		insertQuery := `INSERT INTO note_tags (tag_id, user_id, tag, count, usn, created_time, updated_time, is_deleted) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := db.DB.Exec(insertQuery,
			noteTag.TagId, noteTag.UserId, noteTag.Tag, noteTag.Count,
			noteTag.Usn, noteTag.CreatedTime, noteTag.UpdatedTime, noteTag.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
		}
	}

	return noteTag
}

// 得到标签, 按更新时间来排序
func (this *TagService) GetTags(userId string) []info.NoteTag {
	tags := []info.NoteTag{}
	query := "SELECT tag_id, user_id, tag, count, usn, created_time, updated_time, is_deleted FROM note_tags WHERE user_id = $1 AND is_deleted = false ORDER BY updated_time DESC"

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		Log(err.Error())
		return tags
	}
	defer rows.Close()

	for rows.Next() {
		var noteTag info.NoteTag
		err := rows.Scan(
			&noteTag.TagId, &noteTag.UserId, &noteTag.Tag, &noteTag.Count,
			&noteTag.Usn, &noteTag.CreatedTime, &noteTag.UpdatedTime, &noteTag.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		tags = append(tags, noteTag)
	}
	return tags
}

// 删除标签
// 也删除所有的笔记含该标签的
// 返回noteId => usn
func (this *TagService) DeleteTag(userId string, tag string) map[string]int {
	usn := incrUsn(userId)
	query := "UPDATE note_tags SET usn = $1, is_deleted = true WHERE user_id = $2 AND tag = $3"
	_, err := db.DB.Exec(query, usn, userId, tag)
	if err != nil {
		Log(err.Error())
		return map[string]int{}
	}
	return noteService.UpdateNoteToDeleteTag(userId, tag)
}

// 删除标签, 供API调用
func (this *TagService) DeleteTagApi(userId string, tag string, usn int) (ok bool, msg string, toUsn int) {
	noteTag := info.NoteTag{}
	query := "SELECT tag_id, usn FROM note_tags WHERE user_id = $1 AND tag = $2"
	err := db.DB.QueryRow(query, userId, tag).Scan(&noteTag.TagId, &noteTag.Usn)

	if err == sql.ErrNoRows || noteTag.TagId == "" {
		return false, "notExists", 0
	}
	if noteTag.Usn > usn {
		return false, "conflict", 0
	}
	toUsn = incrUsn(userId)
	updateQuery := "UPDATE note_tags SET usn = $1, is_deleted = true WHERE user_id = $2 AND tag = $3"
	_, err = db.DB.Exec(updateQuery, toUsn, userId, tag)
	if err != nil {
		Log(err.Error())
		return false, "", 0
	}
	return true, "", toUsn
}

// 重新统计标签的count
func (this *TagService) reCountTagCount(userId string, tags []string) {
	if tags == nil {
		return
	}
	for _, tag := range tags {
		this.AddOrUpdateTag(userId, tag)
	}
}

// 同步用
func (this *TagService) GeSyncTags(userId string, afterUsn, maxEntry int) []info.NoteTag {
	noteTags := []info.NoteTag{}
	query := "SELECT tag_id, user_id, tag, count, usn, created_time, updated_time, is_deleted FROM note_tags WHERE user_id = $1 AND usn > $2 ORDER BY usn LIMIT $3"

	rows, err := db.DB.Query(query, userId, afterUsn, maxEntry)
	if err != nil {
		Log(err.Error())
		return noteTags
	}
	defer rows.Close()

	for rows.Next() {
		var noteTag info.NoteTag
		err := rows.Scan(
			&noteTag.TagId, &noteTag.UserId, &noteTag.Tag, &noteTag.Count,
			&noteTag.Usn, &noteTag.CreatedTime, &noteTag.UpdatedTime, &noteTag.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		noteTags = append(noteTags, noteTag)
	}
	return noteTags
}

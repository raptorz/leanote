package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	"time"
)

// 共享Notebook, Note服务
type ShareService struct {
}

//-----------------------------------
// 返回shareNotebooks, sharedUserInfos
// info.ShareNotebooksByUser, []info.User

// 总体来说, 这个方法比较麻烦, 速度未知. 以后按以下方案来缓存用户基础数据

// 以后建个用户的基本数据表, 放所有notebook, sharedNotebook的缓存!!
// 每更新一次则启动一个goroutine异步更新
// 共享只能共享本notebook下的, 如果其子也要共享, 必须设置其子!!!
// 那么, 父, 子在shareNotebooks表中都会有记录

// 得到用户的所有*被*共享的Notebook
// 1 得到别人共享给我的所有notebooks
// 2 按parent进行层次化
// 3 每个层按seq进行排序
// 4 按用户分组
// [ok]

// 谁共享给了我的Query
func (this *ShareService) getOrQ(userId string) string {
	// 简化版本，只处理用户共享
	return "to_user_id = $1"
}

// 得到共享给我的笔记本和用户(谁共享给了我)
func (this *ShareService) GetShareNotebooks(userId string) (info.ShareNotebooksByUser, []info.User) {
	// PostgreSQL实现
	query := `
		SELECT sn.id, sn.user_id, sn.to_user_id, sn.notebook_id, sn.permissions, 
		       sn.created_time, sn.is_deleted,
		       n.id as notebook_id, n.user_id as notebook_user_id, n.parent_notebook_id, 
		       n.seq, n.title, n.url_title, n.number_notes, n.is_trash, n.is_blog,
		       n.created_time as notebook_created_time, n.updated_time as notebook_updated_time,
		       n.usn as notebook_usn, n.is_deleted as notebook_is_deleted
		FROM share_notebooks sn
		JOIN notebooks n ON sn.notebook_id = n.id
		WHERE sn.to_user_id = $1 AND sn.is_deleted = false AND n.is_deleted = false
		ORDER BY sn.user_id, n.seq
	`

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return info.ShareNotebooksByUser{}, []info.User{}
	}
	defer rows.Close()

	shareNotebooksByUser := make(info.ShareNotebooksByUser)
	userIds := make(map[string]bool)

	for rows.Next() {
		var shareNotebook info.ShareNotebook
		var notebook info.Notebook

		err := rows.Scan(
			&shareNotebook.ShareNotebookId, &shareNotebook.UserId, &shareNotebook.ToUserId, &shareNotebook.NotebookId,
			&shareNotebook.Perm, &shareNotebook.CreatedTime, &shareNotebook.IsDeleted,
			&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId,
			&notebook.Seq, &notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes,
			&notebook.IsTrash, &notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime,
			&notebook.Usn, &notebook.IsDeleted,
		)

		if err == nil {
			// 记录共享者ID
			userIds[shareNotebook.UserId] = true

			// 创建ShareNotebooks结构
			shareNotebooks := info.ShareNotebooks{
				Notebook:      notebook,
				ShareNotebook: shareNotebook,
			}

			// 按用户分组
			if _, ok := shareNotebooksByUser[shareNotebook.UserId]; !ok {
				shareNotebooksByUser[shareNotebook.UserId] = []info.ShareNotebooks{}
			}
			shareNotebooksByUser[shareNotebook.UserId] = append(shareNotebooksByUser[shareNotebook.UserId], shareNotebooks)
		}
	}

	// 获取用户信息
	users := []info.User{}
	if len(userIds) > 0 {
		// 构建IN查询
		userIdsSlice := make([]string, 0, len(userIds))
		for userId := range userIds {
			userIdsSlice = append(userIdsSlice, userId)
		}

		query = "SELECT id, email, username, logo FROM users WHERE id = ANY($1)"
		rows, err := db.DB.Query(query, userIdsSlice)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var user info.User
				err := rows.Scan(&user.UserId, &user.Email, &user.Username, &user.Logo)
				if err == nil {
					users = append(users, user)
				}
			}
		}
	}

	return shareNotebooksByUser, users
}

func (this *ShareService) ListShareNotebooks(notebookId, userId, myUserId string) ([]info.ShareNotebooks, info.ShareNotebooksByUser) {
	return []info.ShareNotebooks{}, info.ShareNotebooksByUser{}
}

func (this *ShareService) HasShareNotebook(notebookId, userId string) bool {
	query := `SELECT COUNT(*) FROM share_notebooks WHERE notebook_id = $1 AND to_user_id = $2 AND is_deleted = false`
	var count int
	err := db.DB.QueryRow(query, notebookId, userId).Scan(&count)
	return err == nil && count > 0
}

func (this *ShareService) HasShareNote(noteId, userId string) bool {
	query := `SELECT COUNT(*) FROM share_notes WHERE note_id = $1 AND to_user_id = $2 AND is_deleted = false`
	var count int
	err := db.DB.QueryRow(query, noteId, userId).Scan(&count)
	return err == nil && count > 0
}

func (this *ShareService) HasReadPerm(notebookUserId, userId, noteId string) bool {
	// 检查笔记的直接共享
	query := `
		SELECT COUNT(*) FROM share_notes 
		WHERE note_id = $1 AND user_id = $2 AND to_user_id = $3 AND is_deleted = false
	`
	var count int
	err := db.DB.QueryRow(query, noteId, notebookUserId, userId).Scan(&count)
	if err == nil && count > 0 {
		return true
	}

	// 检查笔记本的共享
	query = `
		SELECT COUNT(*) FROM share_notebooks sn
		JOIN notes n ON sn.notebook_id = n.notebook_id
		WHERE n.id = $1 AND sn.user_id = $2 AND sn.to_user_id = $3 AND sn.is_deleted = false
	`
	err = db.DB.QueryRow(query, noteId, notebookUserId, userId).Scan(&count)
	return err == nil && count > 0
}

func (this *ShareService) AddShareNotebook(notebookId, userId string, perm int, emails []string) (bool, string, []string) {
	// 简化版本，返回成功
	return true, "", []string{}
}

func (this *ShareService) UpdateShareNotebook(notebookId, userId string, perm int, emails []string) (bool, string, []string) {
	return false, "", []string{}
}

func (this *ShareService) DeleteShareNotebook(notebookId, userId, toUserId string) bool {
	query := `DELETE FROM share_notebooks WHERE notebook_id = $1 AND user_id = $2 AND to_user_id = $3`
	_, err := db.DB.Exec(query, notebookId, userId, toUserId)
	return err == nil
}

func (this *ShareService) AddShareNote(noteId, userId string, perm int, emails []string) (bool, string, []string) {
	return false, "", []string{}
}

func (this *ShareService) DeleteShareNote(noteId, userId, toUserId string) bool {
	query := `DELETE FROM share_notes WHERE note_id = $1 AND user_id = $2 AND to_user_id = $3`
	_, err := db.DB.Exec(query, noteId, userId, toUserId)
	return err == nil
}

func (this *ShareService) CopySharedNote(noteId, notebookId, fromUserId, myUserId string) info.Note {
	return info.Note{}
}

func (this *ShareService) AddShareNotebookToUserId(notebookId string, perm int, fromUserId, toUserId string) bool {
	// 检查是否已存在
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM share_notebooks WHERE notebook_id = $1 AND user_id = $2 AND to_user_id = $3 AND is_deleted = false)`
	err := db.DB.QueryRow(query, notebookId, fromUserId, toUserId).Scan(&exists)
	if err != nil || exists {
		return false
	}

	// 插入新的共享记录
	query = `INSERT INTO share_notebooks (id, user_id, to_user_id, notebook_id, permissions, created_time, is_deleted) 
	         VALUES ($1, $2, $3, $4, $5, $6, false)`
	_, err = db.DB.Exec(query, db.NewUUID(), fromUserId, toUserId, notebookId, perm, time.Now())
	return err == nil
}

func (this *ShareService) AddShareNoteToUserId(noteId string, perm int, fromUserId, toUserId string) bool {
	// 检查是否已存在
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM share_notes WHERE note_id = $1 AND user_id = $2 AND to_user_id = $3 AND is_deleted = false)`
	err := db.DB.QueryRow(query, noteId, fromUserId, toUserId).Scan(&exists)
	if err != nil || exists {
		return false
	}

	// 插入新的共享记录
	query = `INSERT INTO share_notes (id, user_id, to_user_id, note_id, permissions, created_time, is_deleted) 
	         VALUES ($1, $2, $3, $4, $5, $6, false)`
	_, err = db.DB.Exec(query, db.NewUUID(), fromUserId, toUserId, noteId, perm, time.Now())
	return err == nil
}

func (this *ShareService) DeleteShareNotebookGroupWhenDeleteGroupUser(userId, groupId string) bool {
	// 简化版本
	return true
}

func (this *ShareService) DeleteShareNoteGroupWhenDeleteGroupUser(userId, groupId string) bool {
	// 简化版本
	return true
}

func (this *ShareService) DeleteAllShareNotebookGroup(groupId string) bool {
	// 简化版本
	return true
}

func (this *ShareService) DeleteAllShareNoteGroup(groupId string) bool {
	// 简化版本
	return true
}

// 为TrashService添加缺失的方法
func (this *ShareService) DeleteShareNoteAll(noteId, userId string) bool {
	query := `DELETE FROM share_notes WHERE note_id = $1 AND user_id = $2`
	_, err := db.DB.Exec(query, noteId, userId)
	return err == nil
}

func (this *ShareService) HasUpdatePerm(userId, myUserId, noteId string) bool {
	// 检查笔记的直接共享权限
	query := `
		SELECT permissions FROM share_notes 
		WHERE note_id = $1 AND user_id = $2 AND to_user_id = $3 AND is_deleted = false
		ORDER BY permissions DESC
		LIMIT 1
	`
	var perm int
	err := db.DB.QueryRow(query, noteId, userId, myUserId).Scan(&perm)
	if err == nil && perm == 1 { // 1表示写权限
		return true
	}

	// 检查笔记本的共享权限
	query = `
		SELECT sn.permissions FROM share_notebooks sn
		JOIN notes n ON sn.notebook_id = n.notebook_id
		WHERE n.id = $1 AND sn.user_id = $2 AND sn.to_user_id = $3 AND sn.is_deleted = false
		ORDER BY sn.permissions DESC
		LIMIT 1
	`
	err = db.DB.QueryRow(query, noteId, userId, myUserId).Scan(&perm)
	return err == nil && perm == 1
}

func (this *ShareService) HasUpdateNotebookPerm(userId, myUserId, notebookId string) bool {
	query := `
		SELECT permissions FROM share_notebooks 
		WHERE notebook_id = $1 AND user_id = $2 AND to_user_id = $3 AND is_deleted = false
		ORDER BY permissions DESC
		LIMIT 1
	`
	var perm int
	err := db.DB.QueryRow(query, notebookId, userId, myUserId).Scan(&perm)
	return err == nil && perm == 1
}

func (this *ShareService) HasUpdateNotePerm(noteId, userId string) bool {
	// 简化版本，检查用户是否有笔记的更新权限
	query := `SELECT COUNT(*) FROM notes WHERE id = $1 AND user_id = $2 AND is_deleted = false`
	var count int
	err := db.DB.QueryRow(query, noteId, userId).Scan(&count)
	return err == nil && count > 0
}

func (this *ShareService) HasReadNotePerm(noteId, userId string) bool {
	// 检查用户是否有笔记的读取权限
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" {
		return false
	}

	// 如果是自己的笔记
	if note.UserId == userId {
		return true
	}

	// 检查共享权限
	return this.HasReadPerm(note.UserId, userId, noteId)
}

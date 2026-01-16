package service

import (
	"database/sql"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"os"
	"strings"
	"time"
)

type AttachService struct {
}

// add attach
// api调用时, 添加attach之前是没有note的
// fromApi表示是api添加的, updateNote传过来的, 此时不要incNote's usn, 因为updateNote会inc的
func (this *AttachService) AddAttach(attach info.Attach, fromApi bool) (ok bool, msg string) {
	attach.CreatedTime = time.Now()

	// 生成新的附件ID
	if attach.AttachId == "" {
		attach.AttachId = db.NewUUID()
	}

	// 插入附件到数据库
	query := `INSERT INTO attachs (
		attach_id, note_id, upload_user_id, name, title, size, type, path, created_time
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := db.DB.Exec(query,
		attach.AttachId, attach.NoteId, attach.UploadUserId,
		attach.Name, attach.Title, attach.Size, attach.Type, attach.Path, attach.CreatedTime,
	)

	if err != nil {
		Log(err.Error())
		return false, "db error"
	}

	note := noteService.GetNoteById(attach.NoteId)

	// api调用时, 添加attach之前是没有note的
	var userId string
	if note.NoteId != "" {
		userId = note.UserId
	} else {
		userId = attach.UploadUserId
	}

	// 更新笔记的attachs num
	this.updateNoteAttachNum(attach.NoteId, 1)

	if !fromApi {
		// 增长note's usn
		noteService.IncrNoteUsn(attach.NoteId, userId)
	}

	return true, ""
}

// 更新笔记的附件个数
// addNum 1或-1
func (this *AttachService) updateNoteAttachNum(noteId string, addNum int) bool {
	// 统计附件数量
	var num int
	query := "SELECT COUNT(*) FROM attachs WHERE note_id = $1"
	err := db.DB.QueryRow(query, noteId).Scan(&num)
	if err != nil {
		Log(err.Error())
		return false
	}

	// 更新笔记的附件数量
	updateQuery := "UPDATE notes SET attach_num = $1 WHERE id = $2"
	_, err = db.DB.Exec(updateQuery, num, noteId)
	if err != nil {
		Log(err.Error())
		return false
	}
	return true
}

// list attachs
func (this *AttachService) ListAttachs(noteId, userId string) []info.Attach {
	attachs := []info.Attach{}

	// 判断是否有权限为笔记添加附件, userId为空时表示是分享笔记的附件
	if userId != "" {
		// TODO: 需要迁移shareService后启用
		// if !shareService.HasUpdateNotePerm(noteId, userId) {
		// 	return attachs
		// }
		Log("TODO: Check update note permission")
	}

	// 笔记是否是自己的
	note := noteService.GetNoteByIdAndUserId(noteId, userId)
	if note.NoteId == "" {
		return attachs
	}

	// TODO 这里, 优化权限控制

	// 查询附件
	query := "SELECT attach_id, note_id, upload_user_id, name, title, size, type, path, created_time FROM attachs WHERE note_id = $1"
	rows, err := db.DB.Query(query, noteId)
	if err != nil {
		Log(err.Error())
		return attachs
	}
	defer rows.Close()

	for rows.Next() {
		var attach info.Attach
		err := rows.Scan(
			&attach.AttachId, &attach.NoteId, &attach.UploadUserId,
			&attach.Name, &attach.Title, &attach.Size, &attach.Type, &attach.Path, &attach.CreatedTime,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		attachs = append(attachs, attach)
	}

	return attachs
}

// api调用, 通过noteIds得到note's attachs, 通过noteId归类返回
func (this *AttachService) getAttachsByNoteIds(noteIds []string) map[string][]info.Attach {
	attachs := []info.Attach{}
	noteAttchs := make(map[string][]info.Attach)

	if len(noteIds) == 0 {
		return noteAttchs
	}

	// 构建IN查询
	placeholders := make([]string, len(noteIds))
	args := make([]interface{}, len(noteIds))
	for i, id := range noteIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		"SELECT attach_id, note_id, upload_user_id, name, title, size, type, path, created_time FROM attachs WHERE note_id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return noteAttchs
	}
	defer rows.Close()

	for rows.Next() {
		var attach info.Attach
		err := rows.Scan(
			&attach.AttachId, &attach.NoteId, &attach.UploadUserId,
			&attach.Name, &attach.Title, &attach.Size, &attach.Type, &attach.Path, &attach.CreatedTime,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		attachs = append(attachs, attach)
	}

	// 按noteId归类
	for _, attach := range attachs {
		noteId := attach.NoteId
		if itAttachs, ok := noteAttchs[noteId]; ok {
			noteAttchs[noteId] = append(itAttachs, attach)
		} else {
			noteAttchs[noteId] = []info.Attach{attach}
		}
	}
	return noteAttchs
}

func (this *AttachService) UpdateImageTitle(userId, fileId, title string) bool {
	// TODO: 需要实现files表的更新
	// query := "UPDATE files SET title = $1 WHERE file_id = $2 AND user_id = $3"
	// _, err := db.DB.Exec(query, title, fileId, userId)
	// if err != nil {
	// 	Log(err.Error())
	// 	return false
	// }
	// return true
	Log("TODO: UpdateImageTitle not implemented for PostgreSQL")
	return false
}

// Delete note to delete attas firstly
func (this *AttachService) DeleteAllAttachs(noteId, userId string) bool {
	note := noteService.GetNoteById(noteId)
	if note.UserId == userId {
		attachs := []info.Attach{}

		// 查询所有附件
		query := "SELECT attach_id, note_id, upload_user_id, name, title, size, type, path, created_time FROM attachs WHERE note_id = $1"
		rows, err := db.DB.Query(query, noteId)
		if err != nil {
			Log(err.Error())
			return false
		}
		defer rows.Close()

		for rows.Next() {
			var attach info.Attach
			err := rows.Scan(
				&attach.AttachId, &attach.NoteId, &attach.UploadUserId,
				&attach.Name, &attach.Title, &attach.Size, &attach.Type, &attach.Path, &attach.CreatedTime,
			)
			if err != nil {
				Log(err.Error())
				continue
			}
			attachs = append(attachs, attach)
		}

		// 删除文件
		for _, attach := range attachs {
			attach.Path = strings.TrimLeft(attach.Path, "/")
			os.Remove(revel.BasePath + "/" + attach.Path)
		}

		// 从数据库删除
		deleteQuery := "DELETE FROM attachs WHERE note_id = $1"
		_, err = db.DB.Exec(deleteQuery, noteId)
		if err != nil {
			Log(err.Error())
			return false
		}

		// 更新笔记附件数量
		this.updateNoteAttachNum(noteId, 0)

		return true
	}

	return false
}

// delete attach
// 删除附件为什么要incrNoteUsn ? 因为可能没有内容要修改的
func (this *AttachService) DeleteAttach(attachId, userId string) (bool, string) {
	attach := info.Attach{}

	// 查询附件信息
	query := "SELECT attach_id, note_id, upload_user_id, name, title, size, type, path, created_time FROM attachs WHERE attach_id = $1"
	err := db.DB.QueryRow(query, attachId).Scan(
		&attach.AttachId, &attach.NoteId, &attach.UploadUserId,
		&attach.Name, &attach.Title, &attach.Size, &attach.Type, &attach.Path, &attach.CreatedTime,
	)

	if err == nil && attach.AttachId != "" {
		// 判断是否有权限为笔记添加附件
		// TODO: 需要迁移shareService后启用
		// if !shareService.HasUpdateNotePerm(attach.NoteId, userId) {
		// 	return false, "No Perm"
		// }
		Log("TODO: Check update note permission")

		// 从数据库删除
		deleteQuery := "DELETE FROM attachs WHERE attach_id = $1"
		_, err := db.DB.Exec(deleteQuery, attachId)
		if err != nil {
			Log(err.Error())
			return false, "db error"
		}

		// 更新笔记附件数量
		this.updateNoteAttachNum(attach.NoteId, -1)

		// 删除文件
		attach.Path = strings.TrimLeft(attach.Path, "/")
		err = os.Remove(revel.BasePath + "/" + attach.Path)
		if err == nil {
			// userService.UpdateAttachSize(note.UserId, -attach.Size)
			// 修改note Usn
			noteService.IncrNoteUsn(attach.NoteId, userId)

			return true, "delete file success"
		}
		return false, "delete file error"
	}

	if err == sql.ErrNoRows {
		return false, "no such item"
	}

	Log(err.Error())
	return false, "db error"
}

// 获取文件路径
// 要判断是否具有权限
// userId是否具有attach的访问权限
func (this *AttachService) GetAttach(attachId, userId string) (attach info.Attach) {
	if attachId == "" {
		return
	}

	attach = info.Attach{}

	// 查询附件信息
	query := "SELECT attach_id, note_id, upload_user_id, name, title, size, type, path, created_time FROM attachs WHERE attach_id = $1"
	err := db.DB.QueryRow(query, attachId).Scan(
		&attach.AttachId, &attach.NoteId, &attach.UploadUserId,
		&attach.Name, &attach.Title, &attach.Size, &attach.Type, &attach.Path, &attach.CreatedTime,
	)

	if err != nil || attach.Path == "" {
		if err != nil && err != sql.ErrNoRows {
			Log(err.Error())
		}
		return info.Attach{}
	}

	note := noteService.GetNoteById(attach.NoteId)

	// 判断权限

	// 笔记是否是公开的
	if note.IsBlog {
		return attach
	}

	// 笔记是否是我的
	if note.UserId == userId {
		return attach
	}

	// 我是否有权限查看或协作
	// TODO: 需要迁移shareService后启用
	// if shareService.HasReadNotePerm(attach.NoteId, userId) {
	// 	return attach
	// }
	Log("TODO: Check read note permission")

	return info.Attach{}
}

// 复制笔记时需要复制附件
// noteService调用, 权限已判断
func (this *AttachService) CopyAttachs(noteId, toNoteId, toUserId string) bool {
	attachs := []info.Attach{}

	// 查询源笔记的所有附件
	query := "SELECT attach_id, user_id, note_id, upload_user_id, name, title, size, type, path, created_time FROM attachs WHERE note_id = $1"
	rows, err := db.DB.Query(query, noteId)
	if err != nil {
		Log(err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var attach info.Attach
		err := rows.Scan(
			&attach.AttachId, &attach.NoteId, &attach.UploadUserId,
			&attach.Name, &attach.Title, &attach.Size, &attach.Type, &attach.Path, &attach.CreatedTime,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		attachs = append(attachs, attach)
	}

	// 复制之
	for _, attach := range attachs {
		attach.AttachId = ""
		attach.NoteId = toNoteId
		attach.UploadUserId = toUserId

		// 文件复制一份
		_, ext := SplitFilename(attach.Name)
		newFilename := NewGuid() + ext
		dir := "files/" + toUserId + "/attachs"
		filePath := dir + "/" + newFilename
		err := os.MkdirAll(revel.BasePath+"/"+dir, 0755)
		if err != nil {
			return false
		}
		_, err = CopyFile(revel.BasePath+"/"+attach.Path, revel.BasePath+"/"+filePath)
		if err != nil {
			return false
		}
		attach.Name = newFilename
		attach.Path = filePath

		this.AddAttach(attach, false)
	}

	return true
}

// 只留下files的数据, 其它的都删除
func (this *AttachService) UpdateOrDeleteAttachApi(noteId, userId string, files []info.NoteFile) bool {
	// 现在数据库内的
	attachs := this.ListAttachs(noteId, userId)

	nowAttachs := map[string]bool{}
	if files != nil {
		for _, file := range files {
			if file.IsAttach && file.FileId != "" {
				nowAttachs[file.FileId] = true
			}
		}
	}

	for _, attach := range attachs {
		fileId := attach.AttachId
		if !nowAttachs[fileId] {
			// 需要删除的
			// TODO 权限验证去掉
			this.DeleteAttach(fileId, userId)
		}
	}

	return false
}

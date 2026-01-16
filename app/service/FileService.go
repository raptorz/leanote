package service

import (
	"database/sql"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/lib/pq"
	"strings"
	"time"
)

const DEFAULT_ALBUM_ID = "52d3e8ac99c37b7f0d000001"

type FileService struct {
}

// add Image
func (this *FileService) AddImage(image info.File, albumId, userId string, needCheckSize bool) (ok bool, msg string) {
	image.CreatedTime = time.Now()
	if albumId != "" {
		image.AlbumId = albumId
	} else {
		image.AlbumId = DEFAULT_ALBUM_ID
		image.IsDefaultAlbum = true
	}
	image.UserId = userId

	ok = db.Insert(db.Files, image)
	return
}

// list images
// if albumId == "" get default album images
func (this *FileService) ListImagesWithPage(userId, albumId, key string, pageNumber, pageSize int) info.Page {
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, "CreatedTime", false)
	files := []info.File{}

	// Build PostgreSQL query
	query := "SELECT * FROM files WHERE user_id = $1"
	args := []interface{}{userId}
	argIndex := 2

	if albumId != "" {
		query += fmt.Sprintf(" AND album_id = $%d", argIndex)
		args = append(args, albumId)
		argIndex++
	} else {
		query += " AND is_default_album = true"
	}
	if key != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+key+"%")
		argIndex++
	}

	// Get total count
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)
	var count int
	err := db.DB.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		Log(err.Error())
		return info.Page{Count: 0, List: []info.File{}}
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s OFFSET $%d LIMIT $%d", sortFieldR, argIndex, argIndex+1)
	args = append(args, skipNum, pageSize)

	// Execute query
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return info.Page{Count: count, List: []info.File{}}
	}
	defer rows.Close()

	for rows.Next() {
		var file info.File
		err := rows.Scan(&file.FileId, &file.UserId, &file.Name, &file.Title, &file.Size, &file.Path, &file.MimeType, &file.CreatedTime, &file.AlbumId, &file.IsDefaultAlbum)
		if err != nil {
			Log(err.Error())
			continue
		}
		files = append(files, file)
	}

	return info.Page{Count: count, List: files}
}

func (this *FileService) UpdateImageTitle(userId, fileId, title string) bool {
	return db.UpdateByIdAndUserIdField(db.Files, fileId, userId, "Title", title)
}

// get all images names
// for upgrade
func (this *FileService) GetAllImageNamesMap(userId string) (m map[string]bool) {
	files := []info.File{}
	query := "SELECT name FROM files WHERE user_id = $1"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		Log(err.Error())
		return make(map[string]bool)
	}
	defer rows.Close()

	for rows.Next() {
		var file info.File
		err := rows.Scan(&file.Name)
		if err != nil {
			Log(err.Error())
			continue
		}
		files = append(files, file)
	}

	m = make(map[string]bool)
	if len(files) == 0 {
		return
	}

	for _, file := range files {
		m[file.Name] = true
	}
	return
}

// delete image
func (this *FileService) DeleteImage(userId, fileId string) (bool, string) {
	file := info.File{}
	db.GetByIdAndUserId(db.Files, fileId, userId, &file)

	if file.FileId != "" {
		if db.DeleteByIdAndUserId(db.Files, fileId, userId) {
			// delete image
			// TODO: 简化版本，不实际删除文件
			Log("Would delete file: " + file.Path)
			return true, ""
		}
		return false, "db error"
	}
	return false, "no such item"
}

// update image title
func (this *FileService) UpdateImage(userId, fileId, title string) bool {
	return db.UpdateByIdAndUserIdField(db.Files, fileId, userId, "Title", title)
}

func (this *FileService) GetFileBase64(userId, fileId string) (str string, mine string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	path := this.GetFile(userId, fileId)

	if path == "" {
		return "", ""
	}

	// 简化版本，返回空
	return "", ""

	// 简化版本，返回空
	return "", ""
}

// 得到图片base64, 图片要在之前添加data:image/png;base64,
func (this *FileService) GetImageBase64(userId, fileId string) string {

	str, mime := this.GetFileBase64(userId, fileId)
	if str == "" {
		return ""
	}
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;base64,%s", mime, str)
	default:
	}
	return "data:image/png;base64," + str
}

// 获取文件路径
// 要判断是否具有权限
// userId是否具有fileId的访问权限
func (this *FileService) GetFile(userId, fileId string) string {
	if fileId == "" {
		return ""
	}

	file := info.File{}
	query := "SELECT * FROM files WHERE id = $1"
	err := db.DB.QueryRow(query, fileId).Scan(
		&file.FileId, &file.UserId, &file.Name, &file.Title, &file.Size,
		&file.Path, &file.MimeType, &file.CreatedTime, &file.AlbumId, &file.IsDefaultAlbum, &file.FromFileId)
	if err != nil {
		return ""
	}
	path := file.Path
	if path == "" {
		return ""
	}

	// 1. 判断权限

	// 是否是我的文件
	if userId != "" && file.UserId == userId {
		return path
	}

	// 得到使用过该fileId的所有笔记NoteId
	// 这些笔记是否有public的, 若有则ok
	// 这些笔记(笔记本)是否有共享给我的, 若有则ok

	noteIds := noteImageService.GetNoteIds(fileId)
	if noteIds != nil && len(noteIds) > 0 {
		// 这些笔记是否有public的
		// Check if any note is a blog
		query := "SELECT COUNT(*) FROM notes WHERE id = ANY($1) AND is_blog = true"
		var count int
		err := db.DB.QueryRow(query, pq.Array(noteIds)).Scan(&count)
		if err == nil && count > 0 {
			return path
		}

		// 2014/12/28 修复, 如果是分享给用户组, 那就不行, 这里可以实现
		for _, noteId := range noteIds {
			note := noteService.GetNoteById(noteId)
			if shareService.HasReadPerm(note.UserId, userId, noteId) {
				return path
			}
		}
		/*
			// 若有共享给我的笔记?
			// 对该笔记可读?
			if db.Has(db.ShareNotes, bson.M{"ToUserId": bson.ObjectIdHex(userId), "NoteId": bson.M{"$in": noteIds}}) {
				return path
			}

			// 笔记本是否共享给我?
			// 通过笔记得到笔记本
			notes := []info.Note{}
			db.ListByQWithFields(db.Notes, bson.M{"_id": bson.M{"$in": noteIds}}, []string{"NotebookId"}, &notes)
			if notes != nil && len(notes) > 0 {
				notebookIds := make([]bson.ObjectId, len(notes))
				for i := 0; i < len(notes); i++ {
					notebookIds[i] = notes[i].NotebookId
				}

				if db.Has(db.ShareNotebooks, bson.M{"ToUserId": bson.ObjectIdHex(userId), "NotebookId": bson.M{"$in": notebookIds}}) {
					return path
				}
			}
		*/
	}

	// 可能是刚复制到owner上, 但内容又没有保存, 所以没有note->imageId的映射, 此时看是否有fromFileId
	if file.FromFileId != "" {
		fromFile := info.File{}
		query := "SELECT * FROM files WHERE id = $1"
		err := db.DB.QueryRow(query, file.FromFileId).Scan(
			&fromFile.FileId, &fromFile.UserId, &fromFile.Name, &fromFile.Title, &fromFile.Size,
			&fromFile.Path, &fromFile.MimeType, &fromFile.CreatedTime, &fromFile.AlbumId, &fromFile.IsDefaultAlbum, &fromFile.FromFileId)
		if err == nil && fromFile.UserId == userId {
			return fromFile.Path
		}
	}

	return ""
}

// 复制共享的笔记时, 复制其中的图片到我本地
// 复制图片
func (this *FileService) CopyImage(userId, fileId, toUserId string) (bool, string) {
	// 是否已经复制过了
	file2 := info.File{}
	query := "SELECT * FROM files WHERE user_id = $1 AND from_file_id = $2"
	err := db.DB.QueryRow(query, toUserId, fileId).Scan(
		&file2.FileId, &file2.UserId, &file2.Name, &file2.Title, &file2.Size,
		&file2.Path, &file2.MimeType, &file2.CreatedTime, &file2.AlbumId, &file2.IsDefaultAlbum, &file2.FromFileId)
	if err == nil && file2.FileId != "" {
		return true, file2.FileId
	}

	// 复制之
	file := info.File{}
	query = "SELECT * FROM files WHERE id = $1 AND user_id = $2"
	err = db.DB.QueryRow(query, fileId, userId).Scan(
		&file.FileId, &file.UserId, &file.Name, &file.Title, &file.Size,
		&file.Path, &file.MimeType, &file.CreatedTime, &file.AlbumId, &file.IsDefaultAlbum, &file.FromFileId)

	if err == sql.ErrNoRows || file.FileId == "" || file.UserId != userId {
		return false, ""
	}

	_, ext := SplitFilename(file.Name)
	guid := NewGuid()
	newFilename := guid + ext

	// TODO 统一目录格式
	// dir := "files/" + toUserId + "/images"
	dir := "files/" + GetRandomFilePath(toUserId, guid) + "/images"
	filePath := dir + "/" + newFilename
	// 简化版本，不实际创建目录和复制文件
	Log("Would copy file from " + file.Path + " to " + filePath)

	fileInfo := info.File{
		Name:       newFilename,
		Title:      file.Title,
		Path:       filePath,
		Size:       file.Size,
		FromFileId: file.FileId,
		UserId:     toUserId,
		MimeType:   file.MimeType,
	}
	id := db.NewUUID()
	fileInfo.FileId = id
	fileId = id
	Ok, _ := this.AddImage(fileInfo, "", toUserId, false)

	if Ok {
		return Ok, id
	}
	return false, ""
}

// 是否是我的文件
func (this *FileService) IsMyFile(userId, fileId string) bool {
	// Check if file exists and belongs to user
	query := "SELECT COUNT(*) FROM files WHERE id = $1 AND user_id = $2"
	var count int
	err := db.DB.QueryRow(query, fileId, userId).Scan(&count)
	if err != nil {
		Log(err.Error())
		return false
	}
	return count > 0
}

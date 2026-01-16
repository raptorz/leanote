package service

import (
	"database/sql"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"strings"
	"time"
)

type NoteService struct {
}

// isBlog 判断笔记本是否是博客
func (this *NoteService) isBlog(notebookId string) bool {
	var isBlog bool
	query := `SELECT is_blog FROM notebooks WHERE id = $1`
	err := db.DB.QueryRow(query, notebookId).Scan(&isBlog)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return isBlog
}

// reCountNotebookNumberNotes 重新统计笔记本下的笔记数目
func (this *NoteService) reCountNotebookNumberNotes(notebookId string) bool {
	var count int
	err := db.DB.QueryRow(
		"SELECT COUNT(*) FROM notes WHERE notebook_id = $1 AND is_trash = false AND (is_deleted = false OR is_deleted IS NULL)",
		notebookId,
	).Scan(&count)
	if err != nil {
		Log(err.Error())
		return false
	}

	_, err = db.DB.Exec(
		"UPDATE notebooks SET number_notes = $1, updated_time = $2 WHERE id = $3",
		count, time.Now(), notebookId,
	)
	if err != nil {
		Log(err.Error())
		return false
	}

	return true
}

// updateNoteInDB 更新笔记到数据库
func (this *NoteService) updateNoteInDB(noteId, userId string, updates map[string]interface{}) bool {
	if len(updates) == 0 {
		return true
	}

	// 构建UPDATE语句
	setParts := make([]string, 0, len(updates))
	values := make([]interface{}, 0, len(updates)+2)

	for key, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, len(values)+1))
		values = append(values, value)
	}

	query := fmt.Sprintf("UPDATE notes SET %s WHERE id = $%d AND user_id = $%d",
		strings.Join(setParts, ", "), len(values)+1, len(values)+2)
	values = append(values, noteId, userId)

	_, err := db.DB.Exec(query, values...)
	if err != nil {
		Log(err.Error())
		return false
	}
	return true
}

// 通过id, userId得到note
func (this *NoteService) GetNote(noteId, userId string) (note info.Note) {
	note = info.Note{}
	query := `SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE id = $1 AND user_id = $2`

	err := db.DB.QueryRow(query, noteId, userId).Scan(
		&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
		&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
		&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
		&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
		&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
		&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
		&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
		&note.Usn, &note.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return
}

// fileService调用
// 不能是已经删除了的, life bug, 客户端删除后, 竟然还能在web上打开
func (this *NoteService) GetNoteById(noteId string) (note info.Note) {
	note = info.Note{}
	if noteId == "" {
		return
	}
	query := `SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE id = $1 AND (is_deleted = false OR is_deleted IS NULL)`

	err := db.DB.QueryRow(query, noteId).Scan(
		&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
		&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
		&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
		&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
		&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
		&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
		&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
		&note.Usn, &note.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return
}
func (this *NoteService) GetNoteByIdAndUserId(noteId, userId string) (note info.Note) {
	note = info.Note{}
	if noteId == "" || userId == "" {
		return
	}
	query := `SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE id = $1 AND user_id = $2 AND (is_deleted = false OR is_deleted IS NULL)`

	err := db.DB.QueryRow(query, noteId, userId).Scan(
		&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
		&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
		&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
		&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
		&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
		&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
		&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
		&note.Usn, &note.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return
}

// 得到blog, blogService用
// 不要传userId, 因为是公开的
func (this *NoteService) GetBlogNote(noteId string) (note info.Note) {
	note = info.Note{}
	query := `SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE id = $1 AND is_blog = true AND is_trash = false AND (is_deleted = false OR is_deleted IS NULL)`

	err := db.DB.QueryRow(query, noteId).Scan(
		&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
		&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
		&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
		&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
		&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
		&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
		&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
		&note.Usn, &note.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return
}

// 通过id, userId得到noteContent
func (this *NoteService) GetNoteContent(noteContentId, userId string) (noteContent info.NoteContent) {
	noteContent = info.NoteContent{}
	query := `SELECT note_id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id
		FROM note_contents WHERE note_id = $1 AND user_id = $2`

	err := db.DB.QueryRow(query, noteContentId, userId).Scan(
		&noteContent.NoteId, &noteContent.UserId, &noteContent.IsBlog,
		&noteContent.Content, &noteContent.Abstract, &noteContent.CreatedTime,
		&noteContent.UpdatedTime, &noteContent.UpdatedUserId,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return
}

// 得到笔记和内容
func (this *NoteService) GetNoteAndContent(noteId, userId string) (noteAndContent info.NoteAndContent) {
	note := this.GetNote(noteId, userId)
	noteContent := this.GetNoteContent(noteId, userId)
	return info.NoteAndContent{note, noteContent}
}

func (this *NoteService) GetNoteBySrc(src, userId string) (note info.Note) {
	note = info.Note{}
	if src == "" {
		return
	}

	query := `SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE user_id = $1 AND src = $2 AND (is_deleted = false OR is_deleted IS NULL)
		ORDER BY usn DESC LIMIT 1`

	err := db.DB.QueryRow(query, userId, src).Scan(
		&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
		&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
		&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
		&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
		&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
		&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
		&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
		&note.Usn, &note.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return
}

func (this *NoteService) GetNoteAndContentBySrc(src, userId string) (noteId string, noteAndContent info.NoteAndContentSep) {
	note := this.GetNoteBySrc(src, userId)
	if note.NoteId != "" {
		noteId = note.NoteId
		noteContent := this.GetNoteContent(note.NoteId, userId)
		return noteId, info.NoteAndContentSep{note, noteContent}
	}
	return
}

// 获取同步的笔记
// > afterUsn的笔记
func (this *NoteService) GetSyncNotes(userId string, afterUsn, maxEntry int) []info.ApiNote {
	notes := []info.Note{}
	query := `SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE user_id = $1 AND usn > $2 AND (is_deleted = false OR is_deleted IS NULL)
		ORDER BY usn LIMIT $3`

	rows, err := db.DB.Query(query, userId, afterUsn, maxEntry)
	if err != nil {
		Log(err.Error())
		return this.ToApiNotes(notes)
	}
	defer rows.Close()

	for rows.Next() {
		var note info.Note
		err := rows.Scan(
			&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
			&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
			&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
			&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
			&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
			&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
			&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
			&note.Usn, &note.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, note)
	}

	return this.ToApiNotes(notes)
}

// note与apiNote的转换
func (this *NoteService) ToApiNotes(notes []info.Note) []info.ApiNote {
	// 2, 得到所有图片, 附件信息
	// 查images表, attachs表
	if len(notes) > 0 {
		noteIds := make([]string, len(notes))
		for i, note := range notes {
			noteIds[i] = note.NoteId
		}
		noteFilesMap := this.getFiles(noteIds)
		// 生成info.ApiNote
		apiNotes := make([]info.ApiNote, len(notes))
		for i, note := range notes {
			noteId := note.NoteId
			apiNotes[i] = this.ToApiNote(&note, noteFilesMap[noteId])
		}
		return apiNotes
	}
	// 返回空的
	return []info.ApiNote{}
}

// note与apiNote的转换
func (this *NoteService) ToApiNote(note *info.Note, files []info.NoteFile) info.ApiNote {
	apiNote := info.ApiNote{
		NoteId:      note.NoteId,
		NotebookId:  note.NotebookId,
		UserId:      note.UserId,
		Title:       note.Title,
		Tags:        note.Tags,
		IsMarkdown:  note.IsMarkdown,
		IsBlog:      note.IsBlog,
		IsTrash:     note.IsTrash,
		IsDeleted:   note.IsDeleted,
		Usn:         note.Usn,
		CreatedTime: note.CreatedTime,
		UpdatedTime: note.UpdatedTime,
		PublicTime:  note.PublicTime,
		Files:       files,
	}
	return apiNote
}

// getDirtyNotes, 把note的图片, 附件信息都发送给客户端
// 客户端保存到本地, 再获取图片, 附件

// 得到所有图片, 附件信息
// 查images表, attachs表
// [待测]
func (this *NoteService) getFiles(noteIds []string) map[string][]info.NoteFile {
	// TODO: 需要迁移noteImageService和attachService后启用
	// noteImages := noteImageService.getImagesByNoteIds(noteIds)
	// noteAttachs := attachService.getAttachsByNoteIds(noteIds)

	noteFilesMap := map[string][]info.NoteFile{}

	for _, noteId := range noteIds {
		noteFiles := []info.NoteFile{}
		// images
		// TODO: 需要迁移后启用
		// if images, ok := noteImages[noteId]; ok {
		// 	for _, image := range images {
		// 		noteFiles = append(noteFiles, info.NoteFile{
		// 			FileId: image.FileId,
		// 			Type:   image.Type,
		// 		})
		// 	}
		// }

		// attach
		// TODO: 需要迁移后启用
		// if attachs, ok := noteAttachs[noteId]; ok {
		// 	for _, attach := range attachs {
		// 		noteFiles = append(noteFiles, info.NoteFile{
		// 			FileId:   attach.AttachId,
		// 			Type:     attach.Type,
		// 			Title:    attach.Title,
		// 			IsAttach: true,
		// 		})
		// 	}
		// }

		noteFilesMap[noteId] = noteFiles
	}

	return noteFilesMap
}

// 列出note, 排序规则, 还有分页
// CreatedTime, UpdatedTime, title 来排序
func (this *NoteService) ListNotes(userId, notebookId string,
	isTrash bool, pageNumber, pageSize int, sortField string, isAsc bool, isBlog bool) (count int, notes []info.Note) {
	notes = []info.Note{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)

	// 构建查询条件
	whereClauses := []string{"user_id = $1", "is_trash = $2", "(is_deleted = false OR is_deleted IS NULL)"}
	args := []interface{}{userId, isTrash}
	argIndex := 3

	if isBlog {
		whereClauses = append(whereClauses, fmt.Sprintf("is_blog = $%d", argIndex))
		args = append(args, true)
		argIndex++
	}
	if notebookId != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("notebook_id = $%d", argIndex))
		args = append(args, notebookId)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// 总记录数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notes WHERE %s", whereClause)
	err := db.DB.QueryRow(countQuery, args...).Scan(&count)
	if err != nil {
		Log(err.Error())
		return 0, notes
	}

	// 排序
	orderBy := sortFieldR
	if !isAsc {
		orderBy += " DESC"
	}

	// 查询数据
	query := fmt.Sprintf(`SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, pageSize, skipNum)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return count, notes
	}
	defer rows.Close()

	for rows.Next() {
		var note info.Note
		err := rows.Scan(
			&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
			&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
			&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
			&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
			&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
			&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
			&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
			&note.Usn, &note.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, note)
	}
	return
}

// 通过noteIds来查询
// ShareService调用
func (this *NoteService) ListNotesByNoteIdsWithPageSort(noteIds []string, userId string,
	pageNumber, pageSize int, sortField string, isAsc bool, isBlog bool) (notes []info.Note) {
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	notes = []info.Note{}
	if len(noteIds) == 0 {
		return notes
	}

	// 构建IN查询
	placeholders := make([]string, len(noteIds))
	args := make([]interface{}, len(noteIds))
	for i, id := range noteIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// 排序
	orderBy := sortFieldR
	if !isAsc {
		orderBy += " DESC"
	}

	query := fmt.Sprintf(`SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE id IN (%s) AND is_trash = false
		ORDER BY %s LIMIT $%d OFFSET $%d`,
		strings.Join(placeholders, ", "), orderBy, len(args)+1, len(args)+2)

	args = append(args, pageSize, skipNum)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return notes
	}
	defer rows.Close()

	for rows.Next() {
		var note info.Note
		err := rows.Scan(
			&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
			&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
			&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
			&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
			&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
			&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
			&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
			&note.Usn, &note.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, note)
	}
	return
}

// shareService调用
func (this *NoteService) ListNotesByNoteIds(noteIds []string) (notes []info.Note) {
	notes = []info.Note{}
	if len(noteIds) == 0 {
		return notes
	}

	// 构建IN查询
	placeholders := make([]string, len(noteIds))
	args := make([]interface{}, len(noteIds))
	for i, id := range noteIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
		FROM notes WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return notes
	}
	defer rows.Close()

	for rows.Next() {
		var note info.Note
		err := rows.Scan(
			&note.NoteId, &note.UserId, &note.CreatedUserId, &note.NotebookId,
			&note.Title, &note.Desc, &note.Src, &note.ImgSrc,
			&note.Tags, &note.IsTrash, &note.IsBlog, &note.UrlTitle,
			&note.IsRecommend, &note.IsTop, &note.HasSelfDefined,
			&note.ReadNum, &note.LikeNum, &note.CommentNum, &note.IsMarkdown,
			&note.AttachNum, &note.CreatedTime, &note.UpdatedTime,
			&note.RecommendTime, &note.PublicTime, &note.UpdatedUserId,
			&note.Usn, &note.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, note)
	}
	return
}

// blog需要
func (this *NoteService) ListNoteContentsByNoteIds(noteIds []string) (notes []info.NoteContent) {
	notes = []info.NoteContent{}
	if len(noteIds) == 0 {
		return notes
	}

	// 构建IN查询
	placeholders := make([]string, len(noteIds))
	args := make([]interface{}, len(noteIds))
	for i, id := range noteIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT note_id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id
		FROM note_contents WHERE note_id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return notes
	}
	defer rows.Close()

	for rows.Next() {
		var noteContent info.NoteContent
		err := rows.Scan(
			&noteContent.NoteId, &noteContent.UserId, &noteContent.IsBlog,
			&noteContent.Content, &noteContent.Abstract, &noteContent.CreatedTime,
			&noteContent.UpdatedTime, &noteContent.UpdatedUserId,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, noteContent)
	}
	return
}

// 只得到abstract, 不需要content
func (this *NoteService) ListNoteAbstractsByNoteIds(noteIds []string) (notes []info.NoteContent) {
	notes = []info.NoteContent{}
	if len(noteIds) == 0 {
		return notes
	}

	// 构建IN查询
	placeholders := make([]string, len(noteIds))
	args := make([]interface{}, len(noteIds))
	for i, id := range noteIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT note_id, abstract FROM note_contents WHERE note_id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return notes
	}
	defer rows.Close()

	for rows.Next() {
		var noteContent info.NoteContent
		err := rows.Scan(&noteContent.NoteId, &noteContent.Abstract)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, noteContent)
	}
	return
}

func (this *NoteService) ListNoteContentByNoteIds(noteIds []string) (notes []info.NoteContent) {
	notes = []info.NoteContent{}
	if len(noteIds) == 0 {
		return notes
	}

	// 构建IN查询
	placeholders := make([]string, len(noteIds))
	args := make([]interface{}, len(noteIds))
	for i, id := range noteIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT note_id, abstract, content FROM note_contents WHERE note_id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return notes
	}
	defer rows.Close()

	for rows.Next() {
		var noteContent info.NoteContent
		err := rows.Scan(&noteContent.NoteId, &noteContent.Abstract, &noteContent.Content)
		if err != nil {
			Log(err.Error())
			continue
		}
		notes = append(notes, noteContent)
	}
	return
}

// 添加笔记
// 首先要判断Notebook是否是Blog, 是的话设为blog
// [ok]

func (this *NoteService) AddNote(note info.Note, fromApi bool) info.Note {
	if note.NoteId == "" {
		note.NoteId = db.NewUUID()
	}

	// 关于创建时间, 可能是客户端发来, 此时判断时间是否有
	note.CreatedTime = FixUrlTime(note.CreatedTime)
	note.UpdatedTime = FixUrlTime(note.UpdatedTime)

	note.IsTrash = false
	note.UpdatedUserId = note.UserId
	note.UrlTitle = GetUrTitle(note.UserId, note.Title, "note", note.NoteId)
	note.Usn = incrUsn(note.UserId)

	notebookId := note.NotebookId

	// api会传IsBlog, web不会传
	if !fromApi {
		// 设为blog
		note.IsBlog = this.isBlog(notebookId)
	}
	//	if note.IsBlog {
	note.PublicTime = note.UpdatedTime
	//	}

	// 插入笔记
	query := `INSERT INTO notes (
		id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
	)`

	_, err := db.DB.Exec(query,
		note.NoteId, note.UserId, note.CreatedUserId, note.NotebookId,
		note.Title, note.Desc, note.Src, note.ImgSrc,
		note.Tags, note.IsTrash, note.IsBlog, note.UrlTitle,
		note.IsRecommend, note.IsTop, note.HasSelfDefined,
		note.ReadNum, note.LikeNum, note.CommentNum, note.IsMarkdown,
		note.AttachNum, note.CreatedTime, note.UpdatedTime,
		note.RecommendTime, note.PublicTime, note.UpdatedUserId,
		note.Usn, note.IsDeleted,
	)
	if err != nil {
		Log(err.Error())
		return info.Note{}
	}

	// tag1 - TODO: 需要迁移tagService后启用
	// tagService.AddTags(note.UserId, note.Tags)

	// recount notebooks' notes number
	this.reCountNotebookNumberNotes(notebookId)

	return note
}

// 添加共享d笔记
func (this *NoteService) AddSharedNote(note info.Note, myUserId string) info.Note {
	// 判断我是否有权限添加 - TODO: 需要迁移shareService后启用
	// if shareService.HasUpdateNotebookPerm(note.UserId, myUserId, note.NotebookId) {
	// 	note.CreatedUserId = myUserId // 是我给共享我的人创建的
	// 	return this.AddNote(note, false)
	// }

	// 临时实现：直接添加
	note.CreatedUserId = myUserId
	return this.AddNote(note, false)
}

// 添加笔记本内容
// [ok]
func (this *NoteService) AddNoteContent(noteContent info.NoteContent) info.NoteContent {

	noteContent.CreatedTime = FixUrlTime(noteContent.CreatedTime)
	noteContent.UpdatedTime = FixUrlTime(noteContent.UpdatedTime)

	noteContent.UpdatedUserId = noteContent.UserId

	// 插入笔记内容
	query := `INSERT INTO note_contents (
		note_id, user_id, is_blog, content, abstract, created_time, updated_time, updated_user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.DB.Exec(query,
		noteContent.NoteId, noteContent.UserId, noteContent.IsBlog,
		noteContent.Content, noteContent.Abstract, noteContent.CreatedTime,
		noteContent.UpdatedTime, noteContent.UpdatedUserId,
	)
	if err != nil {
		Log(err.Error())
		return info.NoteContent{}
	}

	// 更新笔记图片 - TODO: 需要迁移noteImageService后启用
	// noteImageService.UpdateNoteImages(noteContent.UserId, noteContent.NoteId, "", noteContent.Content)

	return noteContent
}

// API, abstract, desc需要这里获取
// 不需要
/*
func (this *NoteService) AddNoteAndContentApi(note info.Note, noteContent info.NoteContent, myUserId bson.ObjectId) info.Note {
	if(note.NoteId.Hex() == "") {
		noteId := bson.NewObjectId();
		note.NoteId = noteId;
	}
	note.CreatedTime = time.Now()
	note.UpdatedTime = note.CreatedTime
	note.IsTrash = false
	note.UpdatedUserId = note.UserId
	note.UrlTitle = GetUrTitle(note.UserId.Hex(), note.Title, "note")
	note.Usn = userService.IncrUsn(note.UserId.Hex())

	// desc这里获取
	desc := SubStringHTMLToRaw(noteContent.Content, 50)
	note.Desc = desc;

	// 设为blog
	notebookId := note.NotebookId.Hex()
	note.IsBlog = notebookService.IsBlog(notebookId)

	if note.IsBlog {
		note.PublicTime = note.UpdatedTime
	}

	db.Insert(db.Notes, note)

	// tag1, 不需要了
//	tagService.AddTags(note.UserId.Hex(), note.Tags)

	// recount notebooks' notes number
	notebookService.ReCountNotebookNumberNotes(notebookId)

	// 这里, 添加到内容中
	abstract := SubStringHTML(noteContent.Content, 200, "")
	noteContent.Abstract = abstract
	this.AddNoteContent(noteContent)

	return note
}*/

// 添加笔记和内容
// 这里使用 info.NoteAndContent 接收?
func (this *NoteService) AddNoteAndContentForController(note info.Note, noteContent info.NoteContent, updatedUserId string) info.Note {
	if note.UserId != updatedUserId {
		// TODO: 需要迁移shareService后启用
		// if !shareService.HasUpdateNotebookPerm(note.UserId, updatedUserId, note.NotebookId) {
		// 	Log("NO AUTH11")
		// 	return info.Note{}
		// } else {
		// 	Log("HAS AUTH -----------")
		// }
		Log("TODO: Check share permission")
	}
	return this.AddNoteAndContent(note, noteContent, updatedUserId)
}
func (this *NoteService) AddNoteAndContent(note info.Note, noteContent info.NoteContent, myUserId string) info.Note {
	if note.NoteId == "" {
		note.NoteId = db.NewUUID()
	}
	noteContent.NoteId = note.NoteId
	if note.UserId != myUserId {
		note = this.AddSharedNote(note, myUserId)
	} else {
		note = this.AddNote(note, false)
	}
	if note.NoteId != "" {
		this.AddNoteContent(noteContent)
	}
	return note
}

func (this *NoteService) AddNoteAndContentApi(note info.Note, noteContent info.NoteContent, myUserId string) info.Note {
	if note.NoteId == "" {
		note.NoteId = db.NewUUID()
	}
	noteContent.NoteId = note.NoteId
	if note.UserId != myUserId {
		note = this.AddSharedNote(note, myUserId)
	} else {
		note = this.AddNote(note, true)
	}
	if note.NoteId != "" {
		this.AddNoteContent(noteContent)
	}
	return note
}

// 修改笔记
// 这里没有判断usn
func (this *NoteService) UpdateNote(updatedUserId, noteId string, needUpdate map[string]interface{}, usn int) (bool, string, int) {
	// 是否存在
	note := this.GetNoteById(noteId)
	if note.NoteId == "" {
		return false, "notExists", 0
	}

	userId := note.UserId
	// updatedUserId 要修改userId的note, 此时需要判断是否有修改权限
	if userId != updatedUserId {
		// TODO: 需要迁移shareService后启用
		// if !shareService.HasUpdatePerm(userId, updatedUserId, noteId) {
		// 	Log("NO AUTH2")
		// 	return false, "noAuth", 0
		// } else {
		// 	Log("HAS AUTH -----------")
		// }
		Log("TODO: Check update permission")
	}

	/*
		// 这里不再判断, 因为controller已经判断了, 删除附件会新增, 所以不用判断
		if usn > 0 && note.Usn != usn {
			Log("有冲突!!")
			Log(note.Usn)
			Log(usn)
			return false, "conflict", 0
		}
	*/

	// 是否已自定义
	if note.IsBlog && note.HasSelfDefined {
		delete(needUpdate, "ImgSrc")
		delete(needUpdate, "Desc")
	}

	needUpdate["updated_user_id"] = updatedUserId

	// 可以将时间传过来
	updatedTime, ok := needUpdate["UpdatedTime"].(time.Time)
	if ok {
		needUpdate["updated_time"] = FixUrlTime(updatedTime)
		delete(needUpdate, "UpdatedTime")
	} else {
		needUpdate["updated_time"] = time.Now()
	}

	afterUsn := incrUsn(userId)
	needUpdate["usn"] = afterUsn

	needRecountTags := false

	// 是否修改了isBlog
	// 也要修改noteContents的IsBlog
	if isBlog, ok := needUpdate["IsBlog"]; ok {
		isBlog2 := isBlog.(bool)
		if note.IsBlog != isBlog2 {
			this.UpdateNoteContentIsBlog(noteId, userId, isBlog2)

			// 重新发布成博客
			if !note.IsBlog {
				needUpdate["PublicTime"] = needUpdate["UpdatedTime"]
			}

			needRecountTags = true
		}
	}

	// 添加tag2
	// TODO 这个tag去掉, 添加tag另外添加, 不要这个
	if _, ok := needUpdate["Tags"]; ok {
		// TODO: 需要迁移tagService后启用
		// tagService.AddTagsI(userId, tags)

		// 如果是博客, 标签改了, 那么重新计算
		if note.IsBlog {
			needRecountTags = true
		}
	}

	// 执行更新
	ok = this.updateNoteInDB(noteId, userId, needUpdate)
	if !ok {
		return ok, "", 0
	}

	if needRecountTags {
		// 重新计算tags
		// TODO: 需要迁移blogService后启用
		// go (func() {
		// 	blogService.ReCountBlogTags(userId)
		// })()
	}

	// 重新获取之
	note = this.GetNoteById(noteId)

	hasRecount := false

	// 如果修改了notebookId, 则更新notebookId'count
	// 两方的notebook也要修改
	notebookIdI := needUpdate["NotebookId"]
	if notebookIdI != nil {
		notebookId := notebookIdI.(string)
		if notebookId != "" {
			this.reCountNotebookNumberNotes(note.NotebookId)
			this.reCountNotebookNumberNotes(notebookId)
			hasRecount = true
		}
	}

	// 不要多次更新, isTrash = false, = true都要重新统计
	if isTrashI, ok := needUpdate["IsTrash"]; ok {
		// 如果是垃圾, 则删除之共享
		isTrash := isTrashI.(bool)
		if isTrash {
			// TODO: 需要迁移shareService后启用
			// shareService.DeleteShareNoteAll(noteId, userId)
		}
		if !hasRecount {
			this.reCountNotebookNumberNotes(note.NotebookId)
		}
	}

	return true, "", afterUsn
}

// 当设置/取消了笔记为博客
func (this *NoteService) UpdateNoteContentIsBlog(noteId, userId string, isBlog bool) {
	query := "UPDATE note_contents SET is_blog = $1, updated_time = $2 WHERE note_id = $3 AND user_id = $4"
	_, err := db.DB.Exec(query, isBlog, time.Now(), noteId, userId)
	if err != nil {
		Log(err.Error())
	}
}

// 附件修改, 增加noteIncr
func (this *NoteService) IncrNoteUsn(noteId, userId string) int {
	afterUsn := incrUsn(userId)
	query := "UPDATE notes SET updated_time = $1, usn = $2 WHERE id = $3 AND user_id = $4"
	_, err := db.DB.Exec(query, time.Now(), afterUsn, noteId, userId)
	if err != nil {
		Log(err.Error())
	}
	return afterUsn
}

// 这里要判断权限, 如果userId != updatedUserId, 那么需要判断权限
// [ok] TODO perm还没测 [del]
func (this *NoteService) UpdateNoteTitle(userId, updatedUserId, noteId, title string) bool {
	// updatedUserId 要修改userId的note, 此时需要判断是否有修改权限
	if userId != updatedUserId {
		// TODO: 需要迁移shareService后启用
		// if !shareService.HasUpdatePerm(userId, updatedUserId, noteId) {
		// 	println("NO AUTH")
		// 	return false
		// }
		Log("TODO: Check update permission")
	}

	afterUsn := incrUsn(userId)
	query := "UPDATE notes SET updated_user_id = $1, title = $2, updated_time = $3, usn = $4 WHERE id = $5 AND user_id = $6"
	_, err := db.DB.Exec(query, updatedUserId, title, time.Now(), afterUsn, noteId, userId)
	if err != nil {
		Log(err.Error())
		return false
	}
	return true
}

// 修改笔记本内容
// [ok] TODO perm未测
// hasBeforeUpdateNote 之前是否更新过note其它信息, 如果有更新, usn不用更新
// TODO abstract这里生成
func (this *NoteService) UpdateNoteContent(updatedUserId, noteId, content, abstract string,
	hasBeforeUpdateNote bool,
	usn int, updatedTime time.Time) (bool, string, int) {
	// 是否已自定义
	note := this.GetNoteById(noteId)
	if note.NoteId == "" {
		return false, "notExists", 0
	}
	userId := note.UserId
	// updatedUserId 要修改userId的note, 此时需要判断是否有修改权限
	if userId != updatedUserId {
		// TODO: 需要迁移shareService后启用
		// if !shareService.HasUpdatePerm(userId, updatedUserId, noteId) {
		// 	Log("NO AUTH")
		// 	return false, "noAuth", 0
		// }
		Log("TODO: Check update permission")
	}

	updatedTime = FixUrlTime(updatedTime)

	// usn, 修改笔记不可能单独修改内容
	afterUsn := 0
	// 如果之前没有修改note其它信息, 那么usn++
	if !hasBeforeUpdateNote {
		// 需要验证
		if usn >= 0 && note.Usn != usn {
			return false, "conflict", 0
		}
		afterUsn = incrUsn(userId)
		// 更新notes表的usn
		query := "UPDATE notes SET usn = $1 WHERE id = $2 AND user_id = $3"
		_, err := db.DB.Exec(query, afterUsn, noteId, userId)
		if err != nil {
			Log(err.Error())
		}
	}

	// 更新note_contents表
	var query string
	var args []interface{}
	if note.IsBlog && note.HasSelfDefined {
		query = "UPDATE note_contents SET updated_user_id = $1, content = $2, updated_time = $3 WHERE note_id = $4 AND user_id = $5"
		args = []interface{}{updatedUserId, content, updatedTime, noteId, userId}
	} else {
		query = "UPDATE note_contents SET updated_user_id = $1, content = $2, abstract = $3, updated_time = $4 WHERE note_id = $5 AND user_id = $6"
		args = []interface{}{updatedUserId, content, abstract, updatedTime, noteId, userId}
	}

	_, err := db.DB.Exec(query, args...)
	if err != nil {
		Log(err.Error())
		return false, "", 0
	}

	// 这里, 添加历史记录 - TODO: 需要迁移noteContentHistoryService后启用
	// noteContentHistoryService.AddHistory(noteId, userId, info.EachHistory{UpdatedUserId: updatedUserId,
	// 	Content:     content,
	// 	UpdatedTime: time.Now(),
	// })

	// 更新笔记图片 - TODO: 需要迁移noteImageService后启用
	// noteImageService.UpdateNoteImages(userId, noteId, note.ImgSrc, content)

	return true, "", afterUsn
}

// ?????
// 这种方式太恶心, 改动很大
// 通过content修改笔记的imageIds列表
// src="http://localhost:9000/file/outputImage?fileId=541ae75499c37b6b79000005&noteId=541ae63c19807a4bb9000000"
func (this *NoteService) updateNoteImages(noteId string, content string) bool {
	return true
}

// 更新tags
// [ok] [del]
func (this *NoteService) UpdateTags(noteId string, userId string, tags []string) bool {
	afterUsn := incrUsn(userId)
	query := "UPDATE notes SET tags = $1, usn = $2 WHERE id = $3 AND user_id = $4"
	_, err := db.DB.Exec(query, tags, afterUsn, noteId, userId)
	if err != nil {
		Log(err.Error())
		return false
	}
	return true
}

func (this *NoteService) ToBlog(userId, noteId string, isBlog, isTop bool) bool {
	if isTop {
		isBlog = true
	}
	if !isBlog {
		isTop = false
	}

	afterUsn := incrUsn(userId)
	var query string
	var args []interface{}

	if isBlog {
		query = "UPDATE notes SET is_blog = $1, is_top = $2, public_time = $3, usn = $4 WHERE id = $5 AND user_id = $6"
		args = []interface{}{isBlog, isTop, time.Now(), afterUsn, noteId, userId}
	} else {
		query = "UPDATE notes SET is_blog = $1, is_top = $2, has_self_defined = $3, usn = $4 WHERE id = $5 AND user_id = $6"
		args = []interface{}{isBlog, isTop, false, afterUsn, noteId, userId}
	}

	_, err := db.DB.Exec(query, args...)
	if err != nil {
		Log(err.Error())
		return false
	}

	// 重新计算tags
	go (func() {
		this.UpdateNoteContentIsBlog(noteId, userId, isBlog)

		// TODO: 需要迁移blogService后启用
		// blogService.ReCountBlogTags(userId)
	})()
	return true
}

// 移动note
// trash, 正常的都可以用
// 1. 要检查下notebookId是否是自己的
// 2. 要判断之前是否是blog, 如果不是, 那么notebook是否是blog?
func (this *NoteService) MoveNote(noteId, notebookId, userId string) info.Note {
	if notebookService.IsMyNotebook(notebookId, userId) {
		note := this.GetNote(noteId, userId)
		preNotebookId := note.NotebookId

		afterUsn := incrUsn(userId)
		query := "UPDATE notes SET is_trash = false, notebook_id = $1, usn = $2 WHERE id = $3 AND user_id = $4"
		_, err := db.DB.Exec(query, notebookId, afterUsn, noteId, userId)

		if err == nil {
			// 更新blog状态
			this.updateToNotebookBlog(noteId, notebookId, userId)

			// recount notebooks' notes number
			this.reCountNotebookNumberNotes(notebookId)
			// 之前不是trash才统计, trash本不在统计中的
			if !note.IsTrash && preNotebookId != notebookId {
				this.reCountNotebookNumberNotes(preNotebookId)
			}
		} else {
			Log(err.Error())
		}

		return this.GetNote(noteId, userId)
	}
	return info.Note{}
}

// 如果自己的blog状态是true, 不用改变,
// 否则, 如果notebookId的blog是true, 则改为true之
// 返回blog状态
// move, copy时用
func (this *NoteService) updateToNotebookBlog(noteId, notebookId, userId string) bool {
	if this.IsBlog(noteId) {
		return true
	}
	if notebookService.IsBlog(notebookId) {
		query := "UPDATE notes SET is_blog = true, public_time = $1 WHERE id = $2 AND user_id = $3"
		_, err := db.DB.Exec(query, time.Now(), noteId, userId)
		if err != nil {
			Log(err.Error())
			return false
		}
		return true
	}
	return false
}

// 判断是否是blog
func (this *NoteService) IsBlog(noteId string) bool {
	var isBlog bool
	query := "SELECT is_blog FROM notes WHERE id = $1"
	err := db.DB.QueryRow(query, noteId).Scan(&isBlog)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return isBlog
}

// 复制note
// 正常的可以用
// 先查, 再新建
// 要检查下notebookId是否是自己的
func (this *NoteService) CopyNote(noteId, notebookId, userId string) info.Note {
	if notebookService.IsMyNotebook(notebookId, userId) {
		note := this.GetNote(noteId, userId)
		noteContent := this.GetNoteContent(noteId, userId)

		// 重新生成noteId
		note.NoteId = db.NewUUID()
		note.NotebookId = notebookId

		noteContent.NoteId = note.NoteId
		note = this.AddNoteAndContent(note, noteContent, note.UserId)

		// 更新blog状态
		isBlog := this.updateToNotebookBlog(note.NoteId, notebookId, userId)

		// recount
		this.reCountNotebookNumberNotes(notebookId)

		note.IsBlog = isBlog

		return note
	}

	return info.Note{}
}

// 复制别人的共享笔记给我
// 将别人可用的图片转为我的图片, 复制图片
func (this *NoteService) CopySharedNote(noteId, notebookId, fromUserId, myUserId string) info.Note {
	// 判断是否共享了给我
	// Log(notebookService.IsMyNotebook(notebookId, myUserId))
	// TODO: 需要迁移shareService后启用
	// if notebookService.IsMyNotebook(notebookId, myUserId) && shareService.HasReadPerm(fromUserId, myUserId, noteId) {
	if notebookService.IsMyNotebook(notebookId, myUserId) {
		note := this.GetNote(noteId, fromUserId)
		if note.NoteId == "" {
			return info.Note{}
		}
		noteContent := this.GetNoteContent(noteId, fromUserId)

		// 重新生成noteId
		note.NoteId = db.NewUUID()
		note.NotebookId = notebookId
		note.UserId = myUserId
		note.IsTop = false
		note.IsBlog = false // 别人的可能是blog

		note.ImgSrc = "" // 为什么清空, 因为图片需要复制, 先清空

		// content
		noteContent.NoteId = note.NoteId
		noteContent.UserId = note.UserId

		// 复制图片, 把note的图片都copy给我, 且修改noteContent图片路径 - TODO: 需要迁移noteImageService后启用
		// noteContent.Content = noteImageService.CopyNoteImages(noteId, fromUserId, note.NoteId, noteContent.Content, myUserId)

		// 复制附件 - TODO: 需要迁移attachService后启用
		// attachService.CopyAttachs(noteId, note.NoteId, myUserId)

		// 添加之
		note = this.AddNoteAndContent(note, noteContent, note.UserId)

		// 更新blog状态
		isBlog := this.updateToNotebookBlog(note.NoteId, notebookId, myUserId)

		// recount
		this.reCountNotebookNumberNotes(notebookId)

		note.IsBlog = isBlog
		return note
	}

	return info.Note{}
}

// 通过noteId得到notebookId
// shareService call
// [ok]
func (this *NoteService) GetNotebookId(noteId string) string {
	var notebookId string
	query := "SELECT notebook_id FROM notes WHERE id = $1"
	err := db.DB.QueryRow(query, noteId).Scan(&notebookId)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return notebookId
}

// ------------------
// 搜索Note, 博客使用了
func (this *NoteService) SearchNote(key, userId string, pageNumber, pageSize int, sortField string, isAsc, isBlog bool) (count int, notes []info.Note) {
	notes = []info.Note{}
	_, _ = parsePageAndSort(pageNumber, pageSize, sortField, isAsc)

	// TODO: Implement PostgreSQL search
	// For now, return empty results
	return 0, []info.Note{}
	/*
		// 利用标题和desc, 不用content
		orQ := []bson.M{
			bson.M{"Title": bson.M{"$regex": bson.RegEx{".*?" + key + ".*", "i"}}},
			bson.M{"Desc": bson.M{"$regex": bson.RegEx{".*?" + key + ".*", "i"}}},
		}
		// 不是trash的
		query := bson.M{"UserId": bson.ObjectIdHex(userId),
			"IsTrash":   false,
			"IsDeleted": false, // 不能搜索已删除了的
			"$or":       orQ,
		}
		if isBlog {
			query["IsBlog"] = true
		}
		q := db.Notes.Find(query)

		// 总记录数
		count, _ = q.Count()

		q.Sort(sortFieldR).
			Skip(skipNum).
			Limit(pageSize).
			All(&notes)

		// 如果 < pageSize 那么搜索content, 且id不在这些id之间的
		if len(notes) < pageSize {
			notes = this.searchNoteFromContent(notes, userId, key, pageSize, sortFieldR, isBlog)
		}
		return
	*/
}

// 搜索noteContents, 补集pageSize个
func (this *NoteService) searchNoteFromContent(notes []info.Note, userId, key string, pageSize int, sortField string, isBlog bool) []info.Note {
	// TODO: Implement PostgreSQL version
	return notes
}
func (this *NoteService) SearchNoteByTags(tags []string, userId string, pageNumber, pageSize int, sortField string, isAsc bool) (count int, notes []info.Note) {
	_, _ = parsePageAndSort(pageNumber, pageSize, sortField, isAsc)

	// TODO: Implement PostgreSQL version
	// For now, return empty results
	return 0, []info.Note{}
}

// ------------
// 统计
func (this *NoteService) CountNote(userId string) int {
	query := "SELECT COUNT(*) FROM notes WHERE is_trash = false AND is_deleted = false"
	var args []interface{}
	if userId != "" {
		query += " AND user_id = $1"
		args = append(args, userId)
	}
	var count int
	err := db.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		Log(err.Error())
		return 0
	}
	return count
}
func (this *NoteService) CountBlog(userId string) int {
	query := "SELECT COUNT(*) FROM notes WHERE is_blog = true AND is_trash = false AND is_deleted = false"
	var args []interface{}
	if userId != "" {
		query += " AND user_id = $1"
		args = append(args, userId)
	}
	var count int
	err := db.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		Log(err.Error())
		return 0
	}
	return count
}

// 通过标签来查询
func (this *NoteService) CountNoteByTag(userId string, tag string) int {
	if tag == "" {
		return 0
	}
	// PostgreSQL使用数组操作符
	query := "SELECT COUNT(*) FROM notes WHERE user_id = $1 AND is_deleted = false AND $2 = ANY(tags)"
	var count int
	err := db.DB.QueryRow(query, userId, tag).Scan(&count)
	if err != nil {
		Log(err.Error())
		return 0
	}
	return count
}

// 删除tag
// 返回所有note的Usn
func (this *NoteService) UpdateNoteToDeleteTag(userId string, targetTag string) map[string]int {
	// TODO: Implement PostgreSQL version
	// For now, return empty map
	return map[string]int{}
}

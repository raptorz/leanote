package service

import (
	"database/sql"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"sort"
	"strings"
	"time"
)

// 笔记本

type NotebookService struct {
}

// 排序
func sortSubNotebooks(eachNotebooks info.SubNotebooks) info.SubNotebooks {
	// 遍历子, 则子往上进行排序
	for _, eachNotebook := range eachNotebooks {
		if eachNotebook.Subs != nil && len(eachNotebook.Subs) > 0 {
			eachNotebook.Subs = sortSubNotebooks(eachNotebook.Subs)
		}
	}

	// 子排完了, 本层排
	sort.Sort(&eachNotebooks)
	return eachNotebooks
}

// 整理(成有关系)并排序
// GetNotebooks()调用
// ShareService调用
func ParseAndSortNotebooks(userNotebooks []info.Notebook, noParentDelete, needSort bool) info.SubNotebooks {
	// 整理成info.Notebooks
	// 第一遍, 建map
	// notebookId => info.Notebooks
	userNotebooksMap := make(map[string]*info.Notebooks, len(userNotebooks))
	for _, each := range userNotebooks {
		newNotebooks := info.Notebooks{Subs: info.SubNotebooks{}}
		newNotebooks.NotebookId = each.NotebookId
		newNotebooks.Title = each.Title
		//		newNotebooks.Title = html.EscapeString(each.Title)
		newNotebooks.Title = strings.Replace(strings.Replace(each.Title, "<script>", "", -1), "</script", "", -1)
		newNotebooks.Seq = each.Seq
		newNotebooks.UserId = each.UserId
		newNotebooks.ParentNotebookId = each.ParentNotebookId
		newNotebooks.NumberNotes = each.NumberNotes
		newNotebooks.IsTrash = each.IsTrash
		newNotebooks.IsBlog = each.IsBlog

		// 存地址
		userNotebooksMap[each.NotebookId] = &newNotebooks
	}

	// 第二遍, 追加到父下

	// 需要删除的id
	needDeleteNotebookId := map[string]bool{}
	for id, each := range userNotebooksMap {
		// 如果有父, 那么追加到父下, 并剪掉当前, 那么最后就只有根的元素
		if each.ParentNotebookId != "" {
			if userNotebooksMap[each.ParentNotebookId] != nil {
				userNotebooksMap[each.ParentNotebookId].Subs = append(userNotebooksMap[each.ParentNotebookId].Subs, each) // Subs是存地址
				// 并剪掉
				// bug
				needDeleteNotebookId[id] = true
				// delete(userNotebooksMap, id)
			} else if noParentDelete {
				// 没有父, 且设置了要删除
				needDeleteNotebookId[id] = true
				// delete(userNotebooksMap, id)
			}
		}
	}

	// 第三遍, 得到所有根
	final := make(info.SubNotebooks, len(userNotebooksMap)-len(needDeleteNotebookId))
	i := 0
	for id, each := range userNotebooksMap {
		if !needDeleteNotebookId[id] {
			final[i] = each
			i++
		}
	}

	// 最后排序
	if needSort {
		return sortSubNotebooks(final)
	}
	return final
}

// 得到某notebook
func (this *NotebookService) GetNotebook(notebookId, userId string) info.Notebook {
	notebook := info.Notebook{}
	query := `SELECT id, user_id, parent_notebook_id, seq, title, url_title, number_notes, 
		is_trash, is_blog, created_time, updated_time, usn, is_deleted 
		FROM notebooks WHERE id = $1 AND user_id = $2`

	err := db.DB.QueryRow(query, notebookId, userId).Scan(
		&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId, &notebook.Seq,
		&notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes, &notebook.IsTrash,
		&notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime, &notebook.Usn,
		&notebook.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return notebook
}

func (this *NotebookService) GetNotebookById(notebookId string) info.Notebook {
	notebook := info.Notebook{}
	query := `SELECT id, user_id, parent_notebook_id, seq, title, url_title, number_notes, 
		is_trash, is_blog, created_time, updated_time, usn, is_deleted 
		FROM notebooks WHERE id = $1`

	err := db.DB.QueryRow(query, notebookId).Scan(
		&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId, &notebook.Seq,
		&notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes, &notebook.IsTrash,
		&notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime, &notebook.Usn,
		&notebook.IsDeleted,
	)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return notebook
}

func (this *NotebookService) GetNotebookByUserIdAndUrlTitle(userId, notebookIdOrUrlTitle string) info.Notebook {
	notebook := info.Notebook{}
	if IsObjectId(notebookIdOrUrlTitle) {
		return this.GetNotebookById(notebookIdOrUrlTitle)
	} else {
		query := `SELECT id, user_id, parent_notebook_id, seq, title, url_title, number_notes, 
			is_trash, is_blog, created_time, updated_time, usn, is_deleted 
			FROM notebooks WHERE user_id = $1 AND url_title = $2`

		err := db.DB.QueryRow(query, userId, encodeValue(notebookIdOrUrlTitle)).Scan(
			&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId, &notebook.Seq,
			&notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes, &notebook.IsTrash,
			&notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime, &notebook.Usn,
			&notebook.IsDeleted,
		)
		if err != nil && err != sql.ErrNoRows {
			Log(err.Error())
		}
		return notebook
	}
}

// 同步的方法
func (this *NotebookService) GeSyncNotebooks(userId string, afterUsn, maxEntry int) []info.Notebook {
	notebooks := []info.Notebook{}
	query := `SELECT id, user_id, parent_notebook_id, seq, title, url_title, number_notes, 
		is_trash, is_blog, created_time, updated_time, usn, is_deleted 
		FROM notebooks WHERE user_id = $1 AND usn > $2 
		ORDER BY usn LIMIT $3`

	rows, err := db.DB.Query(query, userId, afterUsn, maxEntry)
	if err != nil {
		Log(err.Error())
		return notebooks
	}
	defer rows.Close()

	for rows.Next() {
		var notebook info.Notebook
		err := rows.Scan(
			&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId, &notebook.Seq,
			&notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes, &notebook.IsTrash,
			&notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime, &notebook.Usn,
			&notebook.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		notebooks = append(notebooks, notebook)
	}

	return notebooks
}

// 得到用户下所有的notebook
// 排序好之后返回
// [ok]
func (this *NotebookService) GetNotebooks(userId string) info.SubNotebooks {
	userNotebooks := []info.Notebook{}
	query := `SELECT id, user_id, parent_notebook_id, seq, title, url_title, number_notes, 
		is_trash, is_blog, created_time, updated_time, usn, is_deleted 
		FROM notebooks WHERE user_id = $1 AND (is_deleted = false OR is_deleted IS NULL)`

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		Log(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var notebook info.Notebook
		err := rows.Scan(
			&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId, &notebook.Seq,
			&notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes, &notebook.IsTrash,
			&notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime, &notebook.Usn,
			&notebook.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		userNotebooks = append(userNotebooks, notebook)
	}

	if len(userNotebooks) == 0 {
		return nil
	}

	return ParseAndSortNotebooks(userNotebooks, true, true)
}

// share调用, 不需要删除没有父的notebook
// 不需要排序, 因为会重新排序
// 通过notebookIds得到notebooks, 并转成层次有序
func (this *NotebookService) GetNotebooksByNotebookIds(notebookIds []string) info.SubNotebooks {
	if len(notebookIds) == 0 {
		return nil
	}

	userNotebooks := []info.Notebook{}
	// 构建IN查询
	placeholders := make([]string, len(notebookIds))
	args := make([]interface{}, len(notebookIds))
	for i, id := range notebookIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, user_id, parent_notebook_id, seq, title, url_title, number_notes, 
		is_trash, is_blog, created_time, updated_time, usn, is_deleted 
		FROM notebooks WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		Log(err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var notebook info.Notebook
		err := rows.Scan(
			&notebook.NotebookId, &notebook.UserId, &notebook.ParentNotebookId, &notebook.Seq,
			&notebook.Title, &notebook.UrlTitle, &notebook.NumberNotes, &notebook.IsTrash,
			&notebook.IsBlog, &notebook.CreatedTime, &notebook.UpdatedTime, &notebook.Usn,
			&notebook.IsDeleted,
		)
		if err != nil {
			Log(err.Error())
			continue
		}
		userNotebooks = append(userNotebooks, notebook)
	}

	if len(userNotebooks) == 0 {
		return nil
	}

	return ParseAndSortNotebooks(userNotebooks, false, false)
}

// 添加
func (this *NotebookService) AddNotebook(notebook info.Notebook) (bool, info.Notebook) {

	if notebook.NotebookId == "" {
		notebook.NotebookId = db.NewUUID()
	}

	notebook.UrlTitle = GetUrTitle(notebook.UserId, notebook.Title, "notebook", notebook.NotebookId)
	notebook.Usn = UserS.IncrUsn(notebook.UserId)
	now := time.Now()
	notebook.CreatedTime = now
	notebook.UpdatedTime = now

	query := `INSERT INTO notebooks (
		id, user_id, parent_notebook_id, seq, title, url_title, number_notes,
		is_trash, is_blog, created_time, updated_time, usn, is_deleted
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := db.DB.Exec(query,
		notebook.NotebookId, notebook.UserId, notebook.ParentNotebookId, notebook.Seq,
		notebook.Title, notebook.UrlTitle, notebook.NumberNotes, notebook.IsTrash,
		notebook.IsBlog, notebook.CreatedTime, notebook.UpdatedTime, notebook.Usn,
		notebook.IsDeleted,
	)
	if err != nil {
		Log(err.Error())
		return false, notebook
	}
	return true, notebook
}

// 更新笔记, api
func (this *NotebookService) UpdateNotebookApi(userId, notebookId, title, parentNotebookId string, seq, usn int) (bool, string, info.Notebook) {
	if notebookId == "" {
		return false, "notebookIdNotExists", info.Notebook{}
	}

	// 先判断usn是否和数据库的一样, 如果不一样, 则冲突, 不保存
	notebook := this.GetNotebookById(notebookId)
	// 不存在
	if notebook.NotebookId == "" {
		return false, "notExists", notebook
	} else if notebook.Usn != usn {
		return false, "conflict", notebook
	}
	notebook.Usn = incrUsn(userId)
	notebook.Title = title

	// 更新数据库
	query := `UPDATE notebooks SET title = $1, usn = $2, seq = $3, updated_time = $4, parent_notebook_id = $5 
		WHERE id = $6 AND user_id = $7`

	var parentId interface{}
	if parentNotebookId != "" {
		parentId = parentNotebookId
	} else {
		parentId = nil
	}

	_, err := db.DB.Exec(query, title, notebook.Usn, seq, time.Now(), parentId, notebookId, userId)
	if err != nil {
		Log(err.Error())
		return false, "", notebook
	}

	return true, "", this.GetNotebookById(notebookId)
}

// 判断是否是blog
func (this *NotebookService) IsBlog(notebookId string) bool {
	var isBlog bool
	query := `SELECT is_blog FROM notebooks WHERE id = $1`
	err := db.DB.QueryRow(query, notebookId).Scan(&isBlog)
	if err != nil && err != sql.ErrNoRows {
		Log(err.Error())
	}
	return isBlog
}

// 判断是否是我的notebook
func (this *NotebookService) IsMyNotebook(notebookId, userId string) bool {
	var count int
	query := `SELECT COUNT(*) FROM notebooks WHERE id = $1 AND user_id = $2`
	err := db.DB.QueryRow(query, notebookId, userId).Scan(&count)
	if err != nil {
		Log(err.Error())
		return false
	}
	return count > 0
}

// 更新笔记本信息
// 太广, 不用
/*
func (this *NotebookService) UpdateNotebook(notebook info.Notebook) bool {
	return db.UpdateByIdAndUserId2(db.Notebooks, notebook.NotebookId, notebook.UserId, notebook)
}
*/

// 更新笔记本标题
// [ok]
func (this *NotebookService) UpdateNotebookTitle(notebookId, userId, title string) bool {
	usn := UserS.IncrUsn(userId)
	query := `UPDATE notebooks SET title = $1, usn = $2, updated_time = $3 WHERE id = $4 AND user_id = $5`
	_, err := db.DB.Exec(query, title, usn, time.Now(), notebookId, userId)
	if err != nil {
		Log(err.Error())
		return false
	}
	return true
}

// 更新notebook
func (this *NotebookService) UpdateNotebook(userId, notebookId string, needUpdate map[string]interface{}) bool {
	needUpdate["updated_time"] = time.Now()
	needUpdate["usn"] = incrUsn(userId)

	// 构建UPDATE语句
	setParts := make([]string, 0, len(needUpdate))
	values := make([]interface{}, 0, len(needUpdate)+2)

	for key, value := range needUpdate {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, len(values)+1))
		values = append(values, value)
	}

	query := fmt.Sprintf("UPDATE notebooks SET %s WHERE id = $%d AND user_id = $%d",
		strings.Join(setParts, ", "), len(values)+1, len(values)+2)
	values = append(values, notebookId, userId)

	_, err := db.DB.Exec(query, values...)
	if err != nil {
		Log(err.Error())
		return false
	}
	return true
}

// ToBlog or Not
func (this *NotebookService) ToBlog(userId, notebookId string, isBlog bool) bool {
	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		Log(err.Error())
		return false
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				Log(err.Error())
			}
		}
	}()

	// 更新笔记本
	usn := UserS.IncrUsn(userId)
	_, err = tx.Exec("UPDATE notebooks SET is_blog = $1, usn = $2 WHERE id = $3 AND user_id = $4",
		isBlog, usn, notebookId, userId)
	if err != nil {
		Log(err.Error())
		return false
	}

	// 更新笔记
	updateData := map[string]interface{}{
		"is_blog": isBlog,
		"usn":     usn,
	}
	if isBlog {
		updateData["public_time"] = time.Now()
	} else {
		updateData["has_self_defined"] = false
	}

	// 构建UPDATE语句
	setParts := make([]string, 0, len(updateData))
	values := make([]interface{}, 0, len(updateData)+2)

	for key, value := range updateData {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, len(values)+1))
		values = append(values, value)
	}

	query := fmt.Sprintf("UPDATE notes SET %s WHERE user_id = $%d AND notebook_id = $%d",
		strings.Join(setParts, ", "), len(values)+1, len(values)+2)
	values = append(values, userId, notebookId)

	_, err = tx.Exec(query, values...)
	if err != nil {
		Log(err.Error())
		return false
	}

	// 更新note_contents
	// 先查询该notebook下的所有笔记ID
	rows, err := tx.Query("SELECT id FROM notes WHERE user_id = $1 AND notebook_id = $2", userId, notebookId)
	if err != nil {
		Log(err.Error())
		return false
	}
	defer rows.Close()

	var noteIds []string
	for rows.Next() {
		var noteId string
		if err := rows.Scan(&noteId); err != nil {
			Log(err.Error())
			continue
		}
		noteIds = append(noteIds, noteId)
	}

	if len(noteIds) > 0 {
		// 构建IN查询
		placeholders := make([]string, len(noteIds))
		args := make([]interface{}, len(noteIds)+1)
		for i, id := range noteIds {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args[i+1] = id
		}
		args[0] = isBlog

		query = fmt.Sprintf("UPDATE note_contents SET is_blog = $1 WHERE note_id IN (%s)",
			strings.Join(placeholders, ", "))
		_, err = tx.Exec(query, args...)
		if err != nil {
			Log(err.Error())
			return false
		}
	}

	// 重新计算tags
	// TODO: 需要迁移blogService后启用
	// go (func() {
	// 	blogService.ReCountBlogTags(userId)
	// })()

	return true
}

// 查看是否有子notebook
// 先查看该notebookId下是否有notes, 没有则删除
func (this *NotebookService) DeleteNotebook(userId, notebookId string) (bool, string) {
	// 检查是否有子笔记本
	var childCount int
	err := db.DB.QueryRow(
		"SELECT COUNT(*) FROM notebooks WHERE parent_notebook_id = $1 AND user_id = $2 AND (is_deleted = false OR is_deleted IS NULL)",
		notebookId, userId,
	).Scan(&childCount)
	if err != nil {
		Log(err.Error())
		return false, "查询失败"
	}

	if childCount == 0 { // 无子笔记本
		// 检查是否有笔记（不包括回收站）
		var noteCount int
		err = db.DB.QueryRow(
			"SELECT COUNT(*) FROM notes WHERE notebook_id = $1 AND user_id = $2 AND is_trash = false AND (is_deleted = false OR is_deleted IS NULL)",
			notebookId, userId,
		).Scan(&noteCount)
		if err != nil {
			Log(err.Error())
			return false, "查询失败"
		}

		if noteCount == 0 { // 没有笔记
			// 软删除
			usn := UserS.IncrUsn(userId)
			_, err := db.DB.Exec(
				"UPDATE notebooks SET is_deleted = true, usn = $1, updated_time = $2 WHERE id = $3 AND user_id = $4",
				usn, time.Now(), notebookId, userId,
			)
			if err != nil {
				Log(err.Error())
				return false, "删除失败"
			}
			return true, ""
		}
		return false, "笔记本下有笔记"
	} else {
		return false, "笔记本下有子笔记本"
	}
}

// API调用, 删除笔记本, 不作笔记控制
func (this *NotebookService) DeleteNotebookForce(userId, notebookId string, usn int) (bool, string) {
	notebook := this.GetNotebookById(notebookId)
	// 不存在
	if notebook.NotebookId == "" {
		return false, "notExists"
	} else if notebook.Usn != usn {
		return false, "conflict"
	}

	_, err := db.DB.Exec("DELETE FROM notebooks WHERE id = $1 AND user_id = $2", notebookId, userId)
	if err != nil {
		Log(err.Error())
		return false, "删除失败"
	}
	return true, ""
}

// 排序
// 传入 notebookId => Seq
// 为什么要传入userId, 防止修改其它用户的信息 (恶意)
// [ok]
func (this *NotebookService) SortNotebooks(userId string, notebookId2Seqs map[string]int) bool {
	if len(notebookId2Seqs) == 0 {
		return false
	}

	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		Log(err.Error())
		return false
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				Log(err.Error())
			}
		}
	}()

	for notebookId, seq := range notebookId2Seqs {
		usn := UserS.IncrUsn(userId)
		_, err = tx.Exec(
			"UPDATE notebooks SET seq = $1, usn = $2, updated_time = $3 WHERE id = $4 AND user_id = $5",
			seq, usn, time.Now(), notebookId, userId,
		)
		if err != nil {
			Log(err.Error())
			return false
		}
	}

	return true
}

// 排序和设置父
func (this *NotebookService) DragNotebooks(userId string, curNotebookId string, parentNotebookId string, siblings []string) bool {
	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		Log(err.Error())
		return false
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				Log(err.Error())
			}
		}
	}()

	usn := UserS.IncrUsn(userId)
	// 如果没parentNotebookId, 则parentNotebookId设空
	if parentNotebookId == "" {
		_, err = tx.Exec(
			"UPDATE notebooks SET parent_notebook_id = NULL, usn = $1, updated_time = $2 WHERE id = $3 AND user_id = $4",
			usn, time.Now(), curNotebookId, userId,
		)
	} else {
		_, err = tx.Exec(
			"UPDATE notebooks SET parent_notebook_id = $1, usn = $2, updated_time = $3 WHERE id = $4 AND user_id = $5",
			parentNotebookId, usn, time.Now(), curNotebookId, userId,
		)
	}

	if err != nil {
		Log(err.Error())
		return false
	}

	// 排序
	for seq, notebookId := range siblings {
		usn = incrUsn(userId)
		_, err = tx.Exec(
			"UPDATE notebooks SET seq = $1, usn = $2, updated_time = $3 WHERE id = $4 AND user_id = $5",
			seq, usn, time.Now(), notebookId, userId,
		)
		if err != nil {
			Log(err.Error())
			return false
		}
	}

	return true
}

// 重新统计笔记本下的笔记数目
// noteSevice: AddNote, CopyNote, CopySharedNote, MoveNote
// trashService: DeleteNote (recove不用, 都统一在MoveNote里了)
func (this *NotebookService) ReCountNotebookNumberNotes(notebookId string) bool {
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

func (this *NotebookService) ReCountAll() {
	// 得到所有笔记本
	rows, err := db.DB.Query("SELECT id FROM notebooks")
	if err != nil {
		Log(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var notebookId string
		if err := rows.Scan(&notebookId); err != nil {
			Log(err.Error())
			continue
		}
		this.ReCountNotebookNumberNotes(notebookId)
	}
}

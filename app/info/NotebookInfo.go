package info

import (
	"time"
)

// 在数据库中每个
// 修改字段必须要在NotebookService中修改ParseAndSortNotebooks(没有匿名字段), 以后重构
type Notebook struct {
	NotebookId       string    `db:"id"`
	UserId           string    `db:"user_id"`
	ParentNotebookId string    `db:"parent_notebook_id"`
	Seq              int       `db:"seq"`
	Title            string    `db:"title"`
	UrlTitle         string    `db:"url_title"`
	NumberNotes      int       `db:"number_notes"`
	IsTrash          bool      `db:"is_trash"`
	IsBlog           bool      `db:"is_blog"`
	CreatedTime      time.Time `db:"created_time"`
	UpdatedTime      time.Time `db:"updated_time"`

	Usn       int  `db:"usn"`
	IsDeleted bool `db:"is_deleted"`
}

// 仅仅是为了返回前台
type SubNotebooks []*Notebooks // 存地址, 为了生成tree
type Notebooks struct {
	Notebook
	Subs SubNotebooks // 子notebook 在数据库中是没有的
}

// SubNotebook sort
func (this SubNotebooks) Len() int {
	return len(this)
}
func (this SubNotebooks) Less(i, j int) bool {
	return (*this[i]).Seq < (*this[j]).Seq
}
func (this SubNotebooks) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

/*
修改方案, 因为要共享notebook的问题, 所以还是每个notebook一条记录
{
	notebookId,
	title,
	seq,
	parentNoteBookId, // 上级
	userId
}

得到所有该用户的notebook, 然后组装成tree返回之
更新顺序
添加notebook
更新notebook
删除notebook
*/

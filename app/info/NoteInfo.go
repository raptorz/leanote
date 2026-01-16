package info

import (
	"time"
)

// ObjectId - 为了兼容MongoDB的bson.ObjectId
type ObjectId string

// Hex - 返回ObjectId的十六进制字符串表示
func (o ObjectId) Hex() string {
	return string(o)
}

// 只存笔记基本信息
// 内容不存放
type Note struct {
	NoteId        string `db:"id"`
	UserId        string `db:"user_id"`
	CreatedUserId string `db:"created_user_id"`
	NotebookId    string `db:"notebook_id"`
	Title         string `db:"title"`
	Desc          string `db:"description"`

	Src string `db:"src"`

	ImgSrc string   `db:"img_src"`
	Tags   []string `db:"tags"`

	IsTrash bool `db:"is_trash"`

	IsBlog         bool   `db:"is_blog"`
	UrlTitle       string `db:"url_title"`
	IsRecommend    bool   `db:"is_recommend"`
	IsTop          bool   `db:"is_top"`
	HasSelfDefined bool   `db:"has_self_defined"`

	ReadNum    int `db:"read_num"`
	LikeNum    int `db:"like_num"`
	CommentNum int `db:"comment_num"`

	IsMarkdown bool `db:"is_markdown"`

	AttachNum int `db:"attach_num"`

	CreatedTime   time.Time `db:"created_time"`
	UpdatedTime   time.Time `db:"updated_time"`
	RecommendTime time.Time `db:"recommend_time"`
	PublicTime    time.Time `db:"public_time"`
	UpdatedUserId string    `db:"updated_user_id"`

	Usn int `db:"usn"`

	IsDeleted bool `db:"is_deleted"`
}

// 内容
type NoteContent struct {
	NoteId string `db:"note_id"`
	UserId string `db:"user_id"`

	IsBlog bool `db:"is_blog"`

	Content  string `db:"content"`
	Abstract string `db:"abstract"`

	CreatedTime   time.Time `db:"created_time"`
	UpdatedTime   time.Time `db:"updated_time"`
	UpdatedUserId string    `db:"updated_user_id"`
}

// 基本信息和内容在一起
type NoteAndContent struct {
	Note
	NoteContent
}

// 历史记录
// 每一个历史记录对象
type EachHistory struct {
	UpdatedUserId string    `db:"updated_user_id"`
	UpdatedTime   time.Time `db:"updated_time"`
	Content       string    `db:"content"`
}
type NoteContentHistory struct {
	NoteContentHistoryId string        `db:"id"`
	NoteId               string        `db:"note_id"`
	UserId               string        `db:"user_id"`
	Histories            []EachHistory `db:"histories"`
	CreatedTime          time.Time     `db:"created_time"`
}

// 为了NoteController接收参数

// 更新note或content
// 肯定会传userId(谁的), NoteId
// 会传Title, Content, Tags, 一种或几种
type NoteOrContent struct {
	NotebookId string
	NoteId     string
	UserId     string
	Title      string
	Desc       string
	Src        string
	ImgSrc     string
	Tags       string
	Content    string
	Abstract   string
	IsNew      bool
	IsMarkdown bool
	FromUserId string // 为共享而新建
	IsBlog     bool   // 是否是blog, 更新note不需要修改, 添加note时才有可能用到, 此时需要判断notebook是否设为Blog
}

// 分开的
type NoteAndContentSep struct {
	NoteInfo        Note
	NoteContentInfo NoteContent
}

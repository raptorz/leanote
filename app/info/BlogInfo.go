package info

import (
	"time"
)

type BlogItem struct {
	Note
	Abstract string
	Content  string
	HasMore  bool
	User     User
}

type UserBlogBase struct {
	Logo     string `db:"logo"`
	Title    string `db:"title"`
	SubTitle string `db:"subtitle"`
}

type UserBlogComment struct {
	CanComment  bool   `db:"can_comment"`
	CommentType string `db:"comment_type"`
	DisqusId    string `db:"disqus_id"`
}

type UserBlogStyle struct {
	Style string `db:"style"`
	Css   string `db:"css"`
}

type UserBlog struct {
	UserId   string `db:"id"`
	Logo     string `db:"logo"`
	Title    string `db:"title"`
	SubTitle string `db:"subtitle"`
	AboutMe  string `db:"about_me"`

	CanComment bool `db:"can_comment"`

	CommentType string `db:"comment_type"`
	DisqusId    string `db:"disqus_id"`

	Style string `db:"style"`
	Css   string `db:"css"`

	ThemeId   string `db:"theme_id"`
	ThemePath string `db:"-"`

	CateIds []string            `db:"cate_ids"`
	Singles []map[string]string `db:"singles"`

	PerPageSize int    `db:"per_page_size"`
	SortField   string `db:"sort_field"`
	IsAsc       bool   `db:"is_asc"`

	SubDomain string `db:"sub_domain"`
	Domain    string `db:"domain"`
}

type BlogStat struct {
	NoteId     string `db:"note_id"`
	ReadNum    int    `db:"read_num"`
	LikeNum    int    `db:"like_num"`
	CommentNum int    `db:"comment_num"`
}

type BlogSingle struct {
	SingleId    string    `db:"id"`
	UserId      string    `db:"user_id"`
	Title       string    `db:"title"`
	UrlTitle    string    `db:"slug"`
	Content     string    `db:"content"`
	UpdatedTime time.Time `db:"updated_time"`
	CreatedTime time.Time `db:"created_time"`
	IsPublished bool      `db:"is_published"`
	PublishedTime time.Time `db:"published_time"`
}

type BlogLike struct {
	LikeId      string    `db:"id"`
	BlogId      string    `db:"blog_id"`
	UserId      string    `db:"user_id"`
	CreatedTime time.Time `db:"created_time"`
}

type BlogComment struct {
	CommentId   string    `db:"id"`
	BlogId      string    `db:"blog_id"`
	NoteId      string    `db:"note_id"`
	UserId      string    `db:"user_id"`
	Content     string    `db:"content"`
	ToCommentId string    `db:"to_comment_id"`
	ToUserId    string    `db:"to_user_id"`
	LikeNum     int       `db:"like_num"`
	LikeUserIds []string  `db:"like_user_ids"`
	CreatedTime time.Time `db:"created_time"`
}

type BlogCommentPublic struct {
	BlogComment
	IsILikeIt bool
}

type BlogUrls struct {
	IndexUrl    string
	CateUrl     string
	SearchUrl   string
	SingleUrl   string
	PostUrl     string
	ArchiveUrl  string
	TagsUrl     string
	TagPostsUrl string
}

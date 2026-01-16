package info

import (
	"time"
)

// 第三方类型
const (
	ThirdGithub = iota
	ThirdQQ
)

type User struct {
	UserId      string    `db:"id"`
	Email       string    `db:"email"`
	Verified    bool      `db:"verified"`
	Username    string    `db:"username"`
	UsernameRaw string    `db:"username_raw"`
	Pwd         string    `db:"pwd" json:"-"`
	CreatedTime time.Time `db:"created_time"`

	Logo  string `db:"logo"`
	Theme string `db:"theme"`

	NotebookWidth int  `db:"notebook_width"`
	NoteListWidth int  `db:"note_list_width"`
	MdEditorWidth int  `db:"md_editor_width"`
	LeftIsMin     bool `db:"left_is_min"`

	ThirdUserId   string `db:"third_user_id"`
	ThirdUsername string `db:"third_username"`
	ThirdType     int    `db:"third_type"`

	ImageNum   int    `db:"image_num" json:"-"`
	ImageSize  int64  `db:"image_size" json:"-"`
	AttachNum  int    `db:"attach_num" json:"-"`
	AttachSize int64  `db:"attach_size" json:"-"`
	FromUserId string `db:"from_user_id"`

	AccountType      string    `db:"account_type" json:"-"`
	AccountStartTime time.Time `db:"account_start_time" json:"-"`
	AccountEndTime   time.Time `db:"account_end_time" json:"-"`

	MaxImageNum      int   `db:"max_image_num" json:"-"`
	MaxImageSize     int64 `db:"max_image_size" json:"-"`
	MaxAttachNum     int   `db:"max_attach_num" json:"-"`
	MaxAttachSize    int64 `db:"max_attach_size" json:"-"`
	MaxPerAttachSize int64 `db:"max_per_attach_size" json:"-"`

	Usn            int       `db:"usn"`
	FullSyncBefore time.Time `db:"full_sync_before"`
	IsDeleted      bool      `db:"is_deleted"`
}

type UserAccount struct {
	AccountType      string    `db:"account_type" json:"-"`
	AccountStartTime time.Time `db:"account_start_time" json:"-"`
	AccountEndTime   time.Time `db:"account_end_time" json:"-"`

	MaxImageNum      int   `db:"max_image_num" json:"-"`
	MaxImageSize     int64 `db:"max_image_size" json:"-"`
	MaxAttachNum     int   `db:"max_attach_num" json:"-"`
	MaxAttachSize    int64 `db:"max_attach_size" json:"-"`
	MaxPerAttachSize int64 `db:"max_per_attach_size" json:"-"`
}

// note主页需要
type UserAndBlogUrl struct {
	User
	BlogUrl string `BlogUrl`
	PostUrl string `PostUrl`
}

type UserAndBlog struct {
	UserId    string   `db:"id"`
	Email     string   `db:"email"`
	Username  string   `db:"username"`
	Logo      string   `db:"logo"`
	BlogTitle string   `db:"blog_title"`
	BlogLogo  string   `db:"blog_logo"`
	BlogUrl   string   `db:"blog_url"`

	BlogUrls
}

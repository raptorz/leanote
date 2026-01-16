package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
)

// blog
type BlogService struct {
}

// 得到博客统计信息
// ReadNum, LikeNum, CommentNum
func (this *BlogService) GetBlogStat(noteId string) (stat info.BlogStat) {
	note := noteService.GetBlogNote(noteId)
	stat = info.BlogStat{note.NoteId, note.ReadNum, note.LikeNum, note.CommentNum}
	return
}

// 通过id或urlTitle得到博客
func (this *BlogService) GetBlogByIdAndUrlTitle(userId string, noteIdOrUrlTitle string) (blog info.BlogItem) {
	// TODO: 实现PostgreSQL版本
	return info.BlogItem{}
}

// 通过id得到博客
func (this *BlogService) GetBlog(noteId string) (blog info.BlogItem) {
	// TODO: 实现PostgreSQL版本
	return info.BlogItem{}
}

// 得到用户的所有博客笔记本
func (this *BlogService) GetBlogNotebooks(userId string) (notebooks []info.Notebook) {
	// TODO: 实现PostgreSQL版本
	return []info.Notebook{}
}

// 重新统计博客的标签
func (this *BlogService) ReCountBlogTags(userId string) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

// 得到博客的标签和统计
func (this *BlogService) GetBlogTags(userId string) (tagsMap map[string]int) {
	// TODO: 实现PostgreSQL版本
	return map[string]int{}
}

// 得到博客列表
func (this *BlogService) ListBlogs(userId, notebookId, tag string, page, pageSize int, sortField string, isAsc bool) (info.Page, []info.BlogItem) {
	// TODO: 实现PostgreSQL版本
	return info.Page{}, []info.BlogItem{}
}

// 平台 lea+
// 博客列表
func (this *BlogService) ListAllBlogs(userId, tag string, keywords string, isRecommend bool, page, pageSize int, sorterField string, isAsc bool) (info.Page, []info.BlogItem) {
	// TODO: 实现PostgreSQL版本
	return info.Page{}, []info.BlogItem{}
}

// 上一篇文章, 下一篇文章
func (this *BlogService) PreNextBlog(userId string, sorterField string, isAsc bool, noteId string, baseTime interface{}) (info.Post, info.Post) {
	// TODO: 实现PostgreSQL版本
	return info.Post{}, info.Post{}
}

// ------------------------
// 博客设置
func (this *BlogService) fixUserBlog(userBlog *info.UserBlog) {
	// TODO: 实现PostgreSQL版本
}

func (this *BlogService) GetUserBlog(userId string) info.UserBlog {
	// TODO: 实现PostgreSQL版本
	return info.UserBlog{}
}

// 修改之
func (this *BlogService) UpdateUserBlog(userBlog info.UserBlog) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

// 修改之UserBlogBase
func (this *BlogService) UpdateUserBlogBase(userId string, userBlog info.UserBlogBase) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

func (this *BlogService) UpdateUserBlogComment(userId string, userBlog info.UserBlogComment) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

func (this *BlogService) UpdateUserBlogStyle(userId string, userBlog info.UserBlogStyle) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

// 分页与排序
func (this *BlogService) UpdateUserBlogPaging(userId string, perPageSize int, sortField string, isAsc bool) (ok bool, msg string) {
	// TODO: 实现PostgreSQL版本
	return false, ""
}

// AddOrUpdateSingle creates or updates a single blog post
func (this *BlogService) AddOrUpdateSingle(userId, noteId, title, content string) bool {
	// TODO: 实现PostgreSQL版本
	// This should create or update an "About Me" blog post
	return false
}

// 通过子域名得到博客用户
func (this *BlogService) GetUserBlogBySubDomain(subDomain string) info.UserBlog {
	// TODO: 实现PostgreSQL版本
	return info.UserBlog{}
}

func (this *BlogService) GetUserBlogByDomain(domain string) info.UserBlog {
	// TODO: 实现PostgreSQL版本
	return info.UserBlog{}
}

//---------------------
// 后台管理

// 推荐博客
func (this *BlogService) SetRecommend(noteId string, isRecommend bool) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

// 得到评论
func (this *BlogService) ListComments(noteId string, page, pageSize int, isAll bool) ([]info.UserAndBlog, bool) {
	// TODO: 实现PostgreSQL版本
	return []info.UserAndBlog{}, false
}

func (this *BlogService) IsILikeIt(noteId, userId string) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

// 阅读次数统计+1
func (this *BlogService) IncReadNum(noteId string) bool {
	note := noteService.GetNoteById(noteId)
	if note.IsBlog {
		query := "UPDATE notes SET read_num = read_num + 1 WHERE id = $1"
		_, err := db.DB.Exec(query, noteId)
		if err != nil {
			Log(err.Error())
			return false
		}
		return true
	}
	return false
}

// 点赞
func (this *BlogService) LikeBlog(noteId, userId string) (ok bool, isLike bool) {
	// TODO: 实现PostgreSQL版本
	return false, false
}

// 得到点赞用户
func (this *BlogService) GetLikeUsers(noteId string, page, pageSize int, isAll bool) ([]info.UserAndBlog, bool) {
	// TODO: 实现PostgreSQL版本
	return []info.UserAndBlog{}, false
}

// 添加评论
func (this *BlogService) AddComment(noteId, userId, content string) (bool, info.BlogComment) {
	// TODO: 实现PostgreSQL版本
	return false, info.BlogComment{}
}

// 删除评论
func (this *BlogService) DeleteComment(noteId, commentId, myUserId string) bool {
	// TODO: 实现PostgreSQL版本
	return false
}

// 辅助函数
func (this *BlogService) FixNote(note info.Note) info.Post {
	// TODO: 实现PostgreSQL版本
	return info.Post{}
}

func (this *BlogService) ToPost(note info.Note) info.Post {
	// TODO: 实现PostgreSQL版本
	return info.Post{}
}

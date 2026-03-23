package db

import (
	"fmt"
	"strconv"

	"github.com/leanote/leanote/app/db/common"
	"github.com/leanote/leanote/app/db/mongodb"
	"github.com/leanote/leanote/app/db/postgres"
	"github.com/revel/revel"
)

var globalDatabase Database

// 向后兼容的全局变量
var DB interface{} // *sql.DB for PostgreSQL, *mgo.Session for MongoDB

var Notebooks, Notes, NoteContents interface{}
var NoteContentHistories, ShareNotes, ShareNotebooks interface{}
var HasShareNotes, Blogs, Users, Groups interface{}
var GroupUsers, Tags, NoteTags, TagCounts interface{}
var UserBlogs, Tokens, Suggestions interface{}
var Albums, Files, Attachs, NoteImages interface{}
var Configs, EmailLogs interface{}
var BlogLikes, BlogComments, Reports interface{}
var BlogSingles, Themes, Sessions interface{}

// InitDatabase 初始化数据库（根据配置选择实现）
func InitDatabase() error {
	config := common.Config{
		Type:     revel.Config().StringDefault("db.type", "postgresql"),
		Host:     revel.Config().StringDefault("db.host", "127.0.0.1"),
		Port:     27017, // 默认MongoDB端口
		Username: revel.Config().StringDefault("db.username", ""),
		Password: revel.Config().StringDefault("db.password", ""),
		Database: revel.Config().StringDefault("db.dbname", "leanote"),
		SSLMode:  revel.Config().StringDefault("db.sslmode", "disable"),
	}

	var database Database
	var err error

	switch config.Type {
	case "mongodb":
		database = &mongodb.MongoDatabase{}
		portStr := revel.Config().StringDefault("db.port", "27017")
		if portStr != "" {
			port, _ := strconv.Atoi(portStr)
			config.Port = port
		}
	case "postgresql":
		database = &postgres.PostgresDatabase{}
		portStr := revel.Config().StringDefault("db.port", "5432")
		if portStr != "" {
			port, _ := strconv.Atoi(portStr)
			config.Port = port
		}
	default:
		return fmt.Errorf("unknown database type: %s", config.Type)
	}

	err = database.Initialize(config)
	if err != nil {
		return err
	}

	globalDatabase = database
	setupBackwardCompatibility()

	return nil
}

// setupBackwardCompatibility 设置向后兼容的全局变量
func setupBackwardCompatibility() {
	switch db := globalDatabase.(type) {
	case *mongodb.MongoDatabase:
		// MongoDB Collection变量
		Notebooks = db.Notebooks
		Notes = db.Notes
		NoteContents = db.NoteContents
		NoteContentHistories = db.NoteContentHistories
		ShareNotes = db.ShareNotes
		ShareNotebooks = db.ShareNotebooks
		HasShareNotes = db.HasShareNotes
		Blogs = db.Blogs
		Users = db.Users
		Groups = db.Groups
		GroupUsers = db.GroupUsers
		Tags = db.Tags
		NoteTags = db.NoteTags
		TagCounts = db.TagCounts
		UserBlogs = db.UserBlogs
		Tokens = db.Tokens
		Suggestions = db.Suggestions
		Albums = db.Albums
		Files = db.Files
		Attachs = db.Attachs
		NoteImages = db.NoteImages
		Configs = db.Configs
		EmailLogs = db.EmailLogs
		BlogLikes = db.BlogLikes
		BlogComments = db.BlogComments
		Reports = db.Reports
		BlogSingles = db.BlogSingles
		Themes = db.Themes
		Sessions = db.Sessions
		DB = db.session

	case *postgres.PostgresDatabase:
		// PostgreSQL 不需要Collection变量
		Notebooks = db{}
		Notes = db{}
		NoteContents = db{}
		NoteContentHistories = db{}
		ShareNotes = db{}
		ShareNotebooks = db{}
		HasShareNotes = db{}
		Blogs = db{}
		Users = db{}
		Groups = db{}
		GroupUsers = db{}
		Tags = db{}
		NoteTags = db{}
		TagCounts = db{}
		UserBlogs = db{}
		Tokens = db{}
		Suggestions = db{}
		Albums = db{}
		Files = db{}
		Attachs = db{}
		NoteImages = db{}
		Configs = db{}
		EmailLogs = db{}
		BlogLikes = db{}
		BlogComments = db{}
		Reports = db{}
		BlogSingles = db{}
		Themes = db{}
		Sessions = db{}
		DB = db.db
	}
}

// GetGlobalDB 获取全局数据库实例
func GetGlobalDB() Database {
	return globalDatabase
}

// NewUUID 兼容旧的NewUUID函数
func NewUUID() string {
	if globalDatabase != nil {
		return globalDatabase.NewID()
	}
	return ""
}

// CheckConnection 兼容旧的CheckConnection函数
func CheckConnection() {
	if globalDatabase != nil {
		globalDatabase.CheckConnection()
	}
}

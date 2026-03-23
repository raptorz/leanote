package mongodb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"

	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/db/common"
	"github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
)

type MongoDatabase struct {
	session *mgo.Session
	config  db.Config
	db      *mgo.Database
	idGen   common.IDGenerator

	// Collection对象（保持向后兼容）
	Notebooks            *mgo.Collection
	Notes                *mgo.Collection
	NoteContents         *mgo.Collection
	NoteContentHistories *mgo.Collection
	ShareNotes           *mgo.Collection
	ShareNotebooks       *mgo.Collection
	HasShareNotes        *mgo.Collection
	Blogs                *mgo.Collection
	Users                *mgo.Collection
	Groups               *mgo.Collection
	GroupUsers           *mgo.Collection
	Tags                 *mgo.Collection
	NoteTags             *mgo.Collection
	TagCounts            *mgo.Collection
	UserBlogs            *mgo.Collection
	Tokens               *mgo.Collection
	Suggestions          *mgo.Collection
	Albums               *mgo.Collection
	Files                *mgo.Collection
	Attachs              *mgo.Collection
	NoteImages           *mgo.Collection
	Configs              *mgo.Collection
	EmailLogs            *mgo.Collection
	BlogLikes            *mgo.Collection
	BlogComments         *mgo.Collection
	Reports              *mgo.Collection
	BlogSingles          *mgo.Collection
	Themes               *mgo.Collection
	Sessions             *mgo.Collection
}

func (m *MongoDatabase) Initialize(config db.DatabaseConfig) error {
	var ok bool
	var url string

	if url == "" {
		url, ok = revel.Config.String("db.url")
		if !ok {
			url, ok = revel.Config.String("db.urlEnv")
			if ok {
				lea.Log("get db conf from urlEnv: " + url)
			}
		} else {
			lea.Log("get db conf from db.url: " + url)
		}

		if ok {
			urls := strings.Split(url, "/")
			m.config.Database = urls[len(urls)-1]

			if strings.Contains(m.config.Database, "?") {
				urls = strings.Split(m.config.Database, "?")
				m.config.Database = urls[0]
			}
		}
	}
	if m.config.Database == "" {
		m.config.Database, _ = revel.Config.String("db.dbname")
	}

	if !ok {
		portStr, _ := revel.Config.String("db.port")
		if portStr != "" {
			port, _ := strconv.Atoi(portStr)
			m.config.Port = port
		} else {
			m.config.Port = 27017
		}
		m.config.Host, _ = revel.Config.String("db.host")
		m.config.Username, _ = revel.Config.String("db.username")
		m.config.Password, _ = revel.Config.String("db.password")
		usernameAndPassword := m.config.Username + ":" + m.config.Password + "@"
		if m.config.Username == "" || m.config.Password == "" {
			usernameAndPassword = ""
		}
		url = "mongodb://" + usernameAndPassword + m.config.Host + ":" + fmt.Sprintf("%d", m.config.Port) + "/" + m.config.Database
	}

	lea.Log(url)

	var err error
	m.session, err = mgo.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	m.session.SetMode(mgo.Monotonic, true)
	m.db = m.session.DB(m.config.Database)

	m.initCollections()

	m.idGen = &common.ObjectIdGenerator{}

	lea.Log("Connected to MongoDB database successfully")
	return nil
}

func (m *MongoDatabase) initCollections() {
	m.Notebooks = m.db.C("notebooks")
	m.Notes = m.db.C("notes")
	m.NoteContents = m.db.C("note_contents")
	m.NoteContentHistories = m.db.C("note_content_histories")
	m.ShareNotes = m.db.C("share_notes")
	m.ShareNotebooks = m.db.C("share_notebooks")
	m.HasShareNotes = m.db.C("has_share_notes")
	m.Blogs = m.db.C("blogs")
	m.Users = m.db.C("users")
	m.Groups = m.db.C("groups")
	m.GroupUsers = m.db.C("group_users")
	m.Tags = m.db.C("tags")
	m.NoteTags = m.db.C("note_tags")
	m.TagCounts = m.db.C("tag_count")
	m.UserBlogs = m.db.C("user_blogs")
	m.Tokens = m.db.C("tokens")
	m.Suggestions = m.db.C("suggestions")
	m.Albums = m.db.C("albums")
	m.Files = m.db.C("files")
	m.Attachs = m.db.C("attachs")
	m.NoteImages = m.db.C("note_images")
	m.Configs = m.db.C("configs")
	m.EmailLogs = m.db.C("email_logs")
	m.BlogLikes = m.db.C("blog_likes")
	m.BlogComments = m.db.C("blog_comments")
	m.Reports = m.db.C("reports")
	m.BlogSingles = m.db.C("blog_singles")
	m.Themes = m.db.C("themes")
	m.Sessions = m.db.C("sessions")
}

func (m *MongoDatabase) Close() error {
	if m.session != nil {
		m.session.Close()
	}
	return nil
}

func (m *MongoDatabase) Ping() error {
	return m.session.Ping()
}

func (m *MongoDatabase) IsConnected() bool {
	return m.session != nil && m.Ping() == nil
}

func (m *MongoDatabase) CheckConnection() {
	err := m.session.Ping()
	if err != nil {
		lea.Log("Lost connection to db!")
		m.session.Refresh()
		err = m.session.Ping()
		if err == nil {
			lea.Log("Reconnect to db successful.")
		} else {
			lea.Log("重连失败!!!! 警告告")
		}
	}
}

func (m *MongoDatabase) NewID() string {
	return m.idGen.Generate()
}

func (m *MongoDatabase) IsValidID(id string) bool {
	return m.idGen.IsValid(id)
}

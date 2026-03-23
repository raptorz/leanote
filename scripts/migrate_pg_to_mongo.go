package main

import (
	"flag"
	"log"
	"os"

	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
)

type MigrationConfig struct {
	BatchSize     int
	StopOnError   bool
	ValidateAfter bool
	Verbose       bool
}

type MigrationStats struct {
	TotalRecords   int
	SuccessRecords int
	FailedRecords  int
	StartTime      time.Time
	EndTime        time.Time
}

type Migrator struct {
	MongoURL    string
	PostgresURL string
	MongoDB     *mgo.Database
	PostgresDB  *sql.DB
	IDMap       map[string]string
	IDMapMutex  sync.RWMutex
	Config      MigrationConfig
	Stats       MigrationStats
}

func main() {
	mongoURL := flag.String("mongo-url", "mongodb://localhost:27017/leanote", "MongoDB connection URL")
	postgresURL := flag.String("postgres-url", "host=localhost port=5432 user=leanote password= dbname=leanote sslmode=disable", "PostgreSQL connection URL")
	batchSize := flag.Int("batch-size", 1000, "Batch size for migration")
	stopOnError := flag.Bool("stop-on-error", false, "Stop on error")
	verbose := flag.Bool("verbose", true, "Verbose output")
	flag.Parse()

	config := MigrationConfig{
		BatchSize:   *batchSize,
		StopOnError: *stopOnError,
		Verbose:     *verbose,
	}

	migrator := &Migrator{
		MongoURL:    *mongoURL,
		PostgresURL: *postgresURL,
		IDMap:       make(map[string]string),
		Config:      config,
	}

	log.Println("Starting MongoDB to PostgreSQL migration...")
	log.Printf("MongoDB URL: %s", *mongoURL)
	log.Printf("PostgreSQL URL: %s", *postgresURL)

	err := migrator.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer migrator.Close()

	log.Println("Connected to both databases successfully")

	migrator.Stats.StartTime = time.Now()

	err = migrator.CreateIDMappingTable()
	if err != nil {
		log.Fatalf("Failed to create ID mapping table: %v", err)
	}

	tables := []string{"users", "notebooks", "notes", "note_contents", "tags"}

	for _, table := range tables {
		log.Printf("Migrating table: %s", table)
		total, success, failed := migrator.MigrateTable(table)
		log.Printf("Table %s: %d total, %d success, %d failed", table, total, success, failed)
		migrator.Stats.TotalRecords += total
		migrator.Stats.SuccessRecords += success
		migrator.Stats.FailedRecords += failed

		if failed > 0 && *stopOnError {
			log.Fatal("Migration stopped due to errors")
		}
	}

	migrator.Stats.EndTime = time.Now()

	log.Println("\n=== Migration Summary ===")
	log.Printf("Duration: %v", migrator.Stats.EndTime.Sub(migrator.Stats.StartTime))
	log.Printf("Total Success: %d records", migrator.Stats.SuccessRecords)
	log.Printf("Total Failed: %d records", migrator.Stats.FailedRecords)

	if migrator.Stats.FailedRecords > 0 {
		log.Println("⚠ Migration completed with errors")
		os.Exit(1)
	} else {
		log.Println("✓ Migration completed successfully!")
	}
}

func (m *Migrator) Connect() error {
	var err error

	session, err := mgo.Dial(m.MongoURL)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	m.MongoDB = session.DB("")

	m.PostgresDB, err = sql.Open("postgres", m.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	err = m.PostgresDB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	return nil
}

func (m *Migrator) Close() error {
	if m.PostgresDB != nil {
		m.PostgresDB.Close()
	}
	if m.MongoDB != nil {
		m.MongoDB.Session.Close()
	}
	return nil
}

func (m *Migrator) CreateIDMappingTable() error {
	_, err := m.PostgresDB.Exec(`
		CREATE TABLE IF NOT EXISTS id_mapping (
			object_id VARCHAR(24) PRIMARY KEY,
			uuid VARCHAR(36) NOT NULL UNIQUE,
			table_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (m *Migrator) MigrateTable(tableName string) (int, int, int) {
	switch tableName {
	case "users":
		return m.MigrateUsers()
	case "notebooks":
		return m.MigrateNotebooks()
	case "notes":
		return m.MigrateNotes()
	case "note_contents":
		return m.MigrateNoteContents()
	case "tags":
		return m.MigrateTags()
	default:
		return 0, 0, 0
	}
}

func (m *Migrator) MigrateUsers() (int, int, int) {
	var users []map[string]interface{}
	err := m.MongoDB.C("users").Find(nil).All(&users)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return 0, 0, 0
	}

	total := len(users)
	success := 0
	failed := 0

	query := `INSERT INTO users (
		id, email, verified, username, username_raw, pwd, created_time,
		logo, theme, notebook_width, note_list_width, md_editor_width, left_is_min,
		third_user_id, third_username, third_type, image_num, image_size,
		attach_num, attach_size, from_user_id, account_type, account_start_time,
		account_end_time, max_image_num, max_image_size, max_attach_num,
		max_attach_size, max_per_attach_size, usn, full_sync_before, is_deleted
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
		$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)`

	for i, user := range users {
		userId := m.MapID(getString(user, "_id"), "users")

		_, err := m.PostgresDB.Exec(query,
			userId,
			getString(user, "Email"),
			getBool(user, "Verified"),
			getString(user, "Username"),
			getString(user, "UsernameRaw"),
			getString(user, "Pwd"),
			user["CreatedTime"],
			getString(user, "Logo"),
			getString(user, "Theme"),
			getInt(user, "NotebookWidth"),
			getInt(user, "NoteListWidth"),
			getInt(user, "MdEditorWidth"),
			getBool(user, "LeftIsMin"),
			getString(user, "ThirdUserId"),
			getString(user, "ThirdUsername"),
			getInt(user, "ThirdType"),
			getInt(user, "ImageNum"),
			getInt64(user, "ImageSize"),
			getInt(user, "AttachNum"),
			getInt64(user, "AttachSize"),
			getString(user, "FromUserId"),
			getString(user, "AccountType"),
			user["AccountStartTime"],
			user["AccountEndTime"],
			getInt(user, "MaxImageNums"),
			getInt64(user, "MaxImageSize"),
			getInt(user, "MaxAttachNum"),
			getInt64(user, "MaxAttachSize"),
			getInt64(user, "MaxPerAttachSize"),
			getInt(user, "Usn"),
			user["FullSyncBefore"],
			getBool(user, "IsDeleted"),
		)

		if err != nil {
			failed++
			log.Printf("Error inserting user %d: %v", i, err)
			if m.Config.StopOnError {
				break
			}
		} else {
			success++
		}

		if m.Config.Verbose && (i+1)%100 == 0 {
			log.Printf("Migrated %d/%d users", i+1, total)
		}
	}

	return total, success, failed
}

func (m *Migrator) MigrateNotebooks() (int, int, int) {
	var notebooks []map[string]interface{}
	err := m.MongoDB.C("notebooks").Find(nil).All(&notebooks)
	if err != nil {
		log.Printf("Error fetching notebooks: %v", err)
		return 0, 0, 0
	}

	total := len(notebooks)
	success := 0
	failed := 0

	query := `INSERT INTO notebooks (
		id, user_id, parent_notebook_id, seq, title, url_title, number_notes,
		is_trash, is_blog, created_time, updated_time, usn, is_deleted
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	for i, nb := range notebooks {
		id := getString(nb, "_id")
		notebookId := m.MapID(id, "notebooks")
		userId := m.MapID(getString(nb, "UserId"), "users")
		parentId := ""
		if getString(nb, "ParentNotebookId") != "" {
			parentId = m.MapID(getString(nb, "ParentNotebookId"), "notebooks")
		}

		_, err := m.PostgresDB.Exec(query,
			notebookId,
			userId,
			parentId,
			getInt(nb, "Seq"),
			getString(nb, "Title"),
			getString(nb, "UrlTitle"),
			getInt(nb, "NumberNotes"),
			getBool(nb, "IsTrash"),
			getBool(nb, "IsBlog"),
			nb["CreatedTime"],
			nb["UpdatedTime"],
			getInt(nb, "Usn"),
			getBool(nb, "IsDeleted"),
		)

		if err != nil {
			failed++
			if m.Config.StopOnError {
				break
			}
		} else {
			success++
		}

		if m.Config.Verbose && (i+1)%100 == 0 {
			log.Printf("Migrated %d/%d notebooks", i+1, total)
		}
	}

	return total, success, failed
}

func (m *Migrator) MigrateNotes() (int, int, int) {
	var notes []map[string]interface{}
	err := m.MongoDB.C("notes").Find(nil).All(&notes)
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		return 0, 0, 0
	}

	total := len(notes)
	success := 0
	failed := 0

	query := `INSERT INTO notes (
		id, user_id, created_user_id, notebook_id, title, description, src, img_src,
		tags, is_trash, is_blog, url_title, is_recommend, is_top, has_self_defined,
		read_num, like_num, comment_num, is_markdown, attach_num, created_time,
		updated_time, recommend_time, public_time, updated_user_id, usn, is_deleted
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)`

	for i, note := range notes {
		id := getString(note, "_id")
		noteId := m.MapID(id, "notes")
		userId := m.MapID(getString(note, "UserId"), "users")
		createdUserId := ""
		if getString(note, "CreatedUserId") != "" {
			createdUserId = m.MapID(getString(note, "CreatedUserId"), "users")
		}
		notebookId := m.MapID(getString(note, "NotebookId"), "notebooks")

		_, err := m.PostgresDB.Exec(query,
			noteId,
			userId,
			createdUserId,
			notebookId,
			getString(note, "Title"),
			getString(note, "Desc"),
			getString(note, "Src"),
			getString(note, "ImgSrc"),
			note["Tags"],
			getBool(note, "IsTrash"),
			getBool(note, "IsBlog"),
			getString(note, "UrlTitle"),
			getBool(note, "IsRecommend"),
			getBool(note, "IsTop"),
			getBool(note, "HasSelfDefined"),
			getInt(note, "ReadNum"),
			getInt(note, "LikeNum"),
			getInt(note, "CommentNum"),
			getBool(note, "IsMarkdown"),
			getInt(note, "AttachNum"),
			note["CreatedTime"],
			note["UpdatedTime"],
			note["RecommendTime"],
			note["PublicTime"],
			getString(note, "UpdatedUserId"),
			getInt(note, "Usn"),
			getBool(note, "IsDeleted"),
		)

		if err != nil {
			failed++
			if m.Config.StopOnError {
				break
			}
		} else {
			success++
		}

		if m.Config.Verbose && (i+1)%100 == 0 {
			log.Printf("Migrated %d/%d notes", i+1, total)
		}
	}

	return total, success, failed
}

func (m *Migrator) MigrateNoteContents() (int, int, int) {
	var noteContents []map[string]interface{}
	err := m.MongoDB.C("note_contents").Find(nil).All(&noteContents)
	if err != nil {
		log.Printf("Error fetching note contents: %v", err)
		return 0, 0, 0
	}

	total := len(noteContents)
	success := 0
	failed := 0

	query := `INSERT INTO note_contents (
		note_id, user_id, is_blog, content, abstract, created_time,
		updated_time, updated_user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for i, nc := range noteContents {
		noteId := m.MapID(getString(nc, "_id"), "notes")
		userId := m.MapID(getString(nc, "UserId"), "users")

		_, err := m.PostgresDB.Exec(query,
			noteId,
			userId,
			getBool(nc, "IsBlog"),
			getString(nc, "Content"),
			getString(nc, "Abstract"),
			nc["CreatedTime"],
			nc["UpdatedTime"],
			getString(nc, "UpdatedUserId"),
		)

		if err != nil {
			failed++
			if m.Config.StopOnError {
				break
			}
		} else {
			success++
		}

		if m.Config.Verbose && (i+1)%100 == 0 {
			log.Printf("Migrated %d/%d note contents", i+1, total)
		}
	}

	return total, success, failed
}

func (m *Migrator) MigrateTags() (int, int, int) {
	var tags []map[string]interface{}
	err := m.MongoDB.C("tags").Find(nil).All(&tags)
	if err != nil {
		log.Printf("Error fetching tags: %v", err)
		return 0, 0, 0
	}

	total := len(tags)
	success := 0
	failed := 0

	query := `INSERT INTO tags (
		id, user_id, tag, usn, count, created_time, updated_time, is_deleted
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for i, tag := range tags {
		tagId := m.MapID(getString(tag, "_id"), "tags")
		userId := m.MapID(getString(tag, "UserId"), "users")

		_, err := m.PostgresDB.Exec(query,
			tagId,
			userId,
			getString(tag, "Tag"),
			getInt(tag, "Usn"),
			getInt(tag, "Count"),
			tag["CreatedTime"],
			tag["UpdatedTime"],
			getBool(tag, "IsDeleted"),
		)

		if err != nil {
			failed++
			if m.Config.StopOnError {
				break
			}
		} else {
			success++
		}

		if m.Config.Verbose && (i+1)%100 == 0 {
			log.Printf("Migrated %d/%d tags", i+1, total)
		}
	}

	return total, success, failed
}

func (m *Migrator) MapID(objectId string, tableName string) string {
	m.IDMapMutex.RLock()
	if uuid, exists := m.IDMap[objectId]; exists {
		m.IDMapMutex.RUnlock()
		return uuid
	}
	m.IDMapMutex.RUnlock()

	uuid := objectId + "-0000-0000-0000-000000000000"

	m.IDMapMutex.Lock()
	m.IDMap[objectId] = uuid
	m.IDMapMutex.Unlock()

	go func() {
		_, _ = m.PostgresDB.Exec(
			"INSERT INTO id_mapping (object_id, uuid, table_name) VALUES ($1, $2, $3)",
			objectId, uuid, tableName,
		)
	}()

	return uuid
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int32:
			return int(v)
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return 0
}

func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return int64(v)
		case int32:
			return int64(v)
		case int64:
			return v
		case float64:
			return int64(v)
		}
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

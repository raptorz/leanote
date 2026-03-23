package mongodb_to_postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/leanote/leanote/migration"
)

type MongoToPostgresMigrator struct {
	MongoURL     string
	PostgresURL  string
	MongoDB      *mgo.Database
	PostgresDB   *sql.DB
	IDMap        map[string]string
	IDMapMutex   sync.RWMutex
	Config       migration.MigrationConfig
	ProgressChan chan migration.MigrationProgress
}

func NewMongoToPostgresMigrator(mongoURL, postgresURL string, config migration.MigrationConfig) *MongoToPostgresMigrator {
	return &MongoToPostgresMigrator{
		MongoURL:     mongoURL,
		PostgresURL:  postgresURL,
		IDMap:        make(map[string]string),
		Config:       config,
		ProgressChan: make(chan migration.MigrationProgress, 100),
	}
}

func (m *MongoToPostgresMigrator) Connect() error {
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

func (m *MongoToPostgresMigrator) Close() error {
	var err error

	if m.PostgresDB != nil {
		if e := m.PostgresDB.Close(); e != nil {
			err = fmt.Errorf("PostgreSQL close error: %v", e)
		}
	}

	if m.MongoDB != nil {
		m.MongoDB.Session.Close()
	}

	return err
}

func (m *MongoToPostgresMigrator) PreMigrate() error {
	if m.Config.Verbose {
		log.Println("Pre-migration: Creating ID mapping table...")
	}

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

func (m *MongoToPostgresMigrator) MigrateTable(tableName string) (migration.MigrationResult, error) {
	startTime := time.Now()
	result := migration.MigrationResult{
		TableName: tableName,
	}

	if m.Config.Verbose {
		log.Printf("Starting migration for table: %s", tableName)
	}

	var migrator func() (int, int, []string)

	switch tableName {
	case "users":
		migrator = m.MigrateUsers
	case "notebooks":
		migrator = m.MigrateNotebooks
	case "notes":
		migrator = m.MigrateNotes
	case "note_contents":
		migrator = m.MigrateNoteContents
	case "tags":
		migrator = m.MigrateTags
	case "attachs":
		migrator = m.MigrateAttachs
	case "files":
		migrator = m.MigrateFiles
	case "albums":
		migrator = m.MigrateAlbums
	default:
		return result, fmt.Errorf("unknown table: %s", tableName)
	}

	result.TotalRecords, result.SuccessRecords, result.Errors = migrator()
	result.FailedRecords = result.TotalRecords - result.SuccessRecords
	result.Duration = time.Since(startTime)

	if m.Config.Verbose {
		log.Printf("Migration complete for table %s: %d/%d records in %v",
			tableName, result.SuccessRecords, result.TotalRecords, result.Duration)
	}

	return result, nil
}

func (m *MongoToPostgresMigrator) PostMigrate() error {
	if m.Config.Verbose {
		log.Println("Post-migration: Creating indexes...")
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_id_mapping_uuid ON id_mapping(uuid)",
		"CREATE INDEX IF NOT EXISTS idx_id_mapping_table ON id_mapping(table_name)",
	}

	for _, idx := range indexes {
		if _, err := m.PostgresDB.Exec(idx); err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	return nil
}

func (m *MongoToPostgresMigrator) Validate() (migration.ValidationResult, error) {
	tables := []string{"users", "notebooks", "notes", "note_contents", "tags", "attachs", "files", "albums"}
	results := make([]migration.ValidationResult, 0, len(tables))

	for _, table := range tables {
		result, err := m.ValidateTable(table)
		if err != nil {
			log.Printf("Validation error for table %s: %v", table, err)
			continue
		}
		results = append(results, result)

		if m.Config.Verbose {
			status := "PASSED"
			if !result.Passed {
				status = "FAILED"
			}
			log.Printf("Validation for %s: %s (Source: %d, Target: %d)",
				table, status, result.SourceCount, result.TargetCount)
		}
	}

	allPassed := true
	for _, result := range results {
		if !result.Passed {
			allPassed = false
			break
		}
	}

	finalResult := migration.ValidationResult{
		Passed: allPassed,
	}

	if allPassed {
		finalResult.DataIntegrity = append(finalResult.DataIntegrity, "All tables validated successfully")
	}

	return finalResult, nil
}

func (m *MongoToPostgresMigrator) ValidateTable(tableName string) (migration.ValidationResult, error) {
	result := migration.ValidationResult{
		Table: tableName,
	}

	mongoCount, err := m.MongoDB.C(tableName).Find(nil).Count()
	if err != nil {
		return result, err
	}
	result.SourceCount = mongoCount

	var pgCount int
	err = m.PostgresDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&pgCount)
	if err != nil {
		return result, err
	}
	result.TargetCount = pgCount

	result.RecordsMatch = mongoCount == pgCount
	result.Passed = result.RecordsMatch

	if !result.RecordsMatch {
		result.DataIntegrity = append(result.DataIntegrity,
			fmt.Sprintf("Record count mismatch: MongoDB=%d, PostgreSQL=%d", mongoCount, pgCount))
	}

	return result, nil
}

func (m *MongoToPostgresMigrator) MapID(objectId string, tableName string) string {
	m.IDMapMutex.RLock()
	if uuid, exists := m.IDMap[objectId]; exists {
		m.IDMapMutex.RUnlock()
		return uuid
	}
	m.IDMapMutex.RUnlock()

	uuid := bson.NewObjectId().Hex()

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

func (m *MongoToPostgresMigrator) SendProgress(tableName string, current, total int, recordID string) {
	if m.ProgressChan != nil {
		progress := migration.MigrationProgress{
			CurrentTable:    tableName,
			Progress:        float64(current) / float64(total) * 100,
			TotalRecords:    total,
			MigratedRecords: current,
			CurrentRecord:   recordID,
		}

		select {
		case m.ProgressChan <- progress:
		default:
		}
	}
}

func (m *MongoToPostgresMigrator) MigrateUsers() (int, int, []string) {
	var users []map[string]interface{}
	err := m.MongoDB.C("users").Find(nil).All(&users)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return 0, 0, []string{err.Error()}
	}

	total := len(users)
	success := 0
	errors := []string{}

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
		id := getString(user, "_id")
		userId := m.MapID(id, "users")

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
			errors = append(errors, err.Error())
			if m.Config.StopOnError {
				break
			}
		} else {
			success++
		}

		m.SendProgress("users", i+1, total, userId)
	}

	return total, success, errors
}

func (m *MongoToPostgresMigrator) MigrateNotebooks() (int, int, []string) {
	return 0, 0, []string{}
}

func (m *MongoToPostgresMigrator) MigrateNotes() (int, int, []string) {
	return 0, 0, []string{}
}

func (m *MongoToPostgresMigrator) MigrateNoteContents() (int, int, []string) {
	return 0, 0, []string{}
}

func (m *MongoToPostgresMigrator) MigrateTags() (int, int, []string) {
	return 0, 0, []string{}
}

func (m *MongoToPostgresMigrator) MigrateAttachs() (int, int, []string) {
	return 0, 0, []string{}
}

func (m *MongoToPostgresMigrator) MigrateFiles() (int, int, []string) {
	return 0, 0, []string{}
}

func (m *MongoToPostgresMigrator) MigrateAlbums() (int, int, []string) {
	return 0, 0, []string{}
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

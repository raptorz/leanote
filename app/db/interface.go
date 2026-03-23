package db

// Row PostgreSQL Row类型
type Row interface{}

// Rows PostgreSQL Rows类型
type Rows interface{}

// Result PostgreSQL Result类型
type Result interface{}

// Transaction 事务接口
type Transaction interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
}

// PaginationResult 分页结果
type PaginationResult struct {
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
	Data       []map[string]interface{}
}

// Database 核心数据库接口，提供统一的数据库操作
type Database interface {
	// 连接管理
	Initialize(config DatabaseConfig) error
	Close() error
	Ping() error
	IsConnected() bool

	// ID生成
	NewID() string
	IsValidID(id string) bool

	// 基础CRUD操作（兼容现有接口）
	Insert(table string, data interface{}) bool
	Update(table string, id string, data interface{}) bool
	Delete(table string, id string) bool

	// 批量操作
	BatchInsert(table string, data []interface{}) (int, error)
	BatchUpdate(table string, ids []string, data interface{}) (int, error)
	BatchDelete(table string, ids []string) (int, error)

	// 查询操作
	Select(table string, fields []string, where string, args ...interface{}) ([]map[string]interface{}, error)
	Join(mainTable, joinTable, joinType, onCondition string, where string, args ...interface{}) ([]map[string]interface{}, error)
	Paginate(table string, page, pageSize int, where string, orderBy string, args ...interface{}) (PaginationResult, error)

	// 事务支持
	Begin() (Transaction, error)

	// MongoDB兼容接口（保留旧代码兼容性）
	InsertOld(collection interface{}, data interface{}) bool
	UpdateByIdAndUserId(collection interface{}, id, userId string, data interface{}) bool
	GetByIdAndUserId(collection interface{}, id, userId string, data interface{})
	GetByQ(collection interface{}, query interface{}, data interface{})
	ListByQ(collection interface{}, query interface{}, data interface{})
	Count(collection interface{}, query interface{}) int
	Has(collection interface{}, query interface{}) bool

	// PostgreSQL特定接口
	QueryRow(query string, args ...interface{}) Row
	Query(query string, args ...interface{}) (Rows, error)
	Exec(query string, args ...interface{}) (Result, error)

	// 检查连接状态
	CheckConnection()
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string // "mongodb" or "postgresql"
	URL      string // 完整连接URL
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string // PostgreSQL专用
	Options  map[string]interface{}
}

// MigrationStats 迁移统计信息
type MigrationStats struct {
	TotalTables    int
	SuccessTables  int
	FailedTables   int
	SkippedTables  int
	TotalRecords   int
	SuccessRecords int
	FailedRecords  int
	StartTime      string
	EndTime        string
	Duration       string
}

// DatabaseError 统一数据库错误类型
type DatabaseError struct {
	Code    string
	Message string
	Err     error
	Table   string
	Query   string
}

func (e *DatabaseError) Error() string {
	if e.Err != nil {
		return "Error: " + e.Message
	}
	return "Error: " + e.Message
}

// IsNotFound 判断是否为记录未找到错误
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return errMsg == "not found" || errMsg == "sql: no rows in result"
}

// IsDuplicate 判断是否为重复键错误
func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return errMsg == "duplicate key" || errMsg == "duplicate" ||
		errMsg == "UNIQUE constraint violated"
}

// IsConnectionError 判断是否为连接错误
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return errMsg == "connection refused" ||
		errMsg == "connection reset" ||
		errMsg == "broken pipe" ||
		errMsg == "no connection"
}

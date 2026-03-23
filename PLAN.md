# Leanote 数据库抽象和迁移项目 - 详细执行计划

## 📅 创建日期
2026-03-23

## 🎯 项目目标
将数据库操作部分独立成一个单独的模块，提供 MongoDB 和 PostgreSQL 两个版本，可以通过配置切换选择，并开发数据迁移模块用于双向迁移。

---

## 🏗️ 整体架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (Controllers)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    服务层 (Services)                         │
│  ← 统一使用 db.Database 接口，不再直接调用原生SQL              │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│               数据库接口层 (db.Database)                    │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  interface.go - 统一接口定义                         │   │
│  └─────────────────────────────────────────────────────┘   │
└──────┬──────────────────────────────────────┬──────────────┘
       │                                      │
┌──────▼──────────────┐              ┌────────▼──────────────┐
│   MongoDB 实现      │              │  PostgreSQL 实现     │
│  app/db/mongodb/   │              │  app/db/postgres/    │
│  - client.go       │              │  - client.go         │
│  - database.go     │              │  - database.go       │
└─────────────────────┘              └──────────────────────┘
       │                                      │
┌──────▼──────────────────────────────────────▼──────────────┐
│                   迁移工具 (migration/)                     │
│  ┌─────────────────┐  ┌─────────────────────────────────┐  │
│  │ MongoDB → PG    │  │     PostgreSQL → MongoDB       │  │
│  │ migrate.go      │  │     migrate.go                 │  │
│  └─────────────────┘  └─────────────────────────────────┘  │
│  ┌─────────────────────────────────────────────────────┐  │
│  │          数据验证模块 (validator/)                    │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 📝 详细任务分解

### 阶段1：完善数据库接口层

#### 1.1 增强接口定义 (`app/db/interface.go`)

新增接口方法：
```go
// 批量操作
BatchInsert(table string, data []interface{}) (int, error)
BatchUpdate(table string, ids []string, data interface{}) (int, error)
BatchDelete(table string, ids []string) (int, error)

// 事务支持
Begin() (Transaction, error)

// 查询增强
Select(table string, fields []string, where string, args ...interface{}) ([]map[string]interface{}, error)
Join(mainTable, joinTable, joinType, onCondition string, where string, args ...interface{}) ([]map[string]struct{}, error)

// 分页查询
Paginate(table string, page, pageSize int, where string, orderBy string, args ...interface{}) (PaginationResult, error)
```

#### 1.2 完善MongoDB实现 (`app/db/mongodb/`)

- 实现 `BatchInsert`、`BatchUpdate`、`BatchDelete`
- 实现 `Begin()` 使用MongoDB Session
- 实现 `Select`、`Join` 使用MongoDB聚合管道
- 实现 `Paginate` 使用skip和limit

#### 1.3 完善PostgreSQL实现 (`app/db/postgres/`)

- 实现批量操作使用COPY或批量INSERT
- 实现 `Begin()` 使用sql.Tx
- 实现 `Select`、`Join` 使用标准SQL JOIN
- 实现 `Paginate` 使用OFFSET和LIMIT

#### 1.4 通用工具 (`app/db/common/`)

- 添加查询构建器 `QueryBuilder`
- 添加条件构建器 `WhereBuilder`
- 添加分页助手 `PaginationHelper`

---

### 阶段2：重构服务层使用统一接口

**策略**：将所有服务层的原生SQL调用改为使用 `db.Database` 接口

#### 2.1 重构UserService (`app/service/UserService.go`)

**改造前**（当前develop分支）：
```go
func AddUser(user info.User) bool {
    query := `INSERT INTO ... VALUES (...)`
    _, err := db.DB.Exec(query, ...)
    return err == nil
}
```

**改造后**：
```go
func AddUser(user info.User) bool {
    if user.UserId == "" {
        user.UserId = db.NewUUID()
    }
    return db.DB.Insert("users", user)
}
```

#### 2.2 重构NoteService (`app/service/NoteService.go`)

复杂查询的改造：
```go
// 改造前
func GetNotes(userId string, tags []string) []info.Note {
    query := `SELECT * FROM notes WHERE user_id = $1 AND tags && $2`
    rows, _ := db.DB.Query(query, userId, pq.Array(tags))
    // ... 处理结果
}

// 改造后
func GetNotes(userId string, tags []string) []info.Note {
    where := db.WhereBuilder{}.Eq("user_id", userId).InArray("tags", tags).Build()
    result := db.DB.Select("notes", nil, where, userId, tags)
    // ... 转换为Note结构
}
```

#### 2.3 重构所有服务文件

需要重构的文件列表：
- `UserService.go`
- `NoteService.go`
- `NotebookService.go`
- `TagService.go`
- `AttachService.go`
- `FileService.go`
- `AuthService.go`
- `GroupService.go`
- `ShareService.go`
- `NoteImageService.go`
- `NoteContentHistoryService.go`
- `SessionService.go`
- `TokenService.go`
- `TrashService.go`
- `PwdService.go`

---

### 阶段3：开发数据迁移工具

**目录结构**：
```
migration/
├── migration.go          # 核心迁移框架
├── types.go             # 通用类型定义
├── mongodb_to_postgres/  # MongoDB → PostgreSQL
│   ├── migrator.go       # 主迁移器
│   ├── tables/           # 各表迁移逻辑
│   │   ├── users.go
│   │   ├── notebooks.go
│   │   ├── notes.go
│   │   └── ...
│   └── validator.go      # 迁移验证
└── postgres_to_mongodb/  # PostgreSQL → MongoDB
    ├── migrator.go
    ├── tables/
    │   ├── users.go
    │   └── ...
    └── validator.go
```

#### 3.1 核心迁移框架 (`migration/migration.go`)

```go
type Migrator interface {
    Connect() error
    Close() error
    PreMigrate() error
    MigrateTable(tableName string) (MigrationResult, error)
    PostMigrate() error
    Validate() (ValidationResult, error)
}

type MigrationResult struct {
    TableName      string

    TotalRecords   int
    SuccessRecords int
    FailedRecords  int
    Errors         []error
    Duration      time.Duration
}

type MigrationProgress struct {
    CurrentTable   string
    Progress       float64
    TotalRecords   int
    MigratedRecords int
}
```

#### 3.2 MongoDB → PostgreSQL迁移器

```go
type MongoToPostgresMigrator struct {
    MongoConfig    DatabaseConfig
    PostgresConfig DatabaseConfig
    MongoDB        *mgo.Database
    PostgresDB     *sql.DB
    IDMap          map[string]string // ObjectId -> UUID
}

func (m *MongoToPostgresMigrator) Migrate() error {
    // 1. 连接两个数据库
    // 2. 按依赖顺序迁移表
    tables := []string{
        "users",
        "notebooks",
        "notes",
        "note_contents",
        "tags",
        "attachs",
        // ... 其他表
    }

    for _, table := range tables {
        result, err := m.MigrateTable(table)
        // 处理结果和错误
    }

    // 3. 验证数据
    return m.Validate()
}
```

#### 3.3 PostgreSQL → MongoDB迁移器

结构类似，方向相反：
```go
type PostgresToMongoMigrator struct {
    PostgresConfig DatabaseConfig
    MongoConfig    DatabaseConfig
    PostgresDB     *sql.DB
    MongoDB        *mgo.Database
    IDMap          map[string]string // UUID -> ObjectId
}
```

#### 3.4 表迁移逻辑示例

```go
func (m *MongoToPostgresMigrator) MigrateUsers() MigrationResult {
    result := MigrationResult{TableName: "users"}

    // 从MongoDB读取
    var mongoUsers []MongoUser
    err := m.MongoDB.C("users").Find(nil).All(&mongoUsers)
    if err != nil {
        result.Errors = append(result.Errors, err)
        return result
    }

    result.TotalRecords = len(mongoUsers)

    // 转换并插入PostgreSQL
    for _, mu := range mongoUsers {
        pgUser := m.convertMongoToPostgresUser(mu)
        err := m.insertPostgresUser(pgUser)
        if err != nil {
            result.FailedRecords++
            result.Errors = append(result.Errors, err)
        } else {
            result.SuccessRecords++
        }
    }

    return result
}
```

---

### 阶段4：实现数据验证

#### 4.1 验证框架

```go
type Validator struct {
    SourceDB     Database
    TargetDB     Database
    Tables       []string
}

type ValidationResult struct {
    Table         string
    SourceCount   int
    TargetCount   int
    MissingInTarget []string
    ExtraInTarget  []string
    DataIntegrity  []string
    Passed       bool
}

func (v *Validator) ValidateAll() []ValidationResult {
    var results []ValidationResult

    for _, table := range v.Tables {
        result := v.ValidateTable(table)
        results = append(results, result)
    }

    return results
}

func (v *Validator) ValidateTable(table string) ValidationResult {
    // 1. 比较记录数
    sourceCount := v.SourceDB.Count(table, nil)
    targetCount := v.TargetDB.Count(table, nil)

    result := ValidationResult{
        Table:       table,
        SourceCount: sourceCount,
        TargetCount: targetCount,
    }

    if sourceCount != targetCount {
        result.Passed = false
        return result
    }

    // 2. 采样比较数据
    result = v.ValidateSampleData(table)

    return result
}
```

#### 4.2 数据一致性检查

- 记录数验证
- 外键关系验证
- 关键字段采样验证
- 总和/平均值等统计验证

---

### 阶段5：配置和文档

#### 5.1 配置文件示例

`conf/app.conf`:
```ini
# 数据库类型选择: mongodb 或 postgresql
db.type=postgresql

# PostgreSQL配置
db.host=127.0.0.1
db.port=5432
db.dbname=leanote
db.username=leanote
db.password=your_password
db.sslmode=disable

# MongoDB配置（可选，用于迁移）
mongodb.host=127.0.0.1
mongodb.port=27017
mongodb.dbname=leanote
```

#### 5.2 迁移工具使用文档

`docs/MIGRATION_GUIDE.md`:
```markdown
# 数据迁移指南

## MongoDB → PostgreSQL 迁移

1. 备份MongoDB数据库
2. 确保PostgreSQL已创建并导入schema
3. 运行迁移工具

```bash
go run migration/mongodb_to_postgres/main.go \
    --mongo-url="mongodb://localhost:27017/leanote" \
    --postgres-url="host=localhost port=5432 user=leanote password=xxx dbname=leanote"
```

4. 验证数据

## PostgreSQL → MongoDB 迁移

（类似步骤）
```

---

## 📊 工作量评估

| 阶段 | 任务 | 预估工时 |
|------|------|---------|
| 1 | 完善数据库接口层 | 8-12小时 |
| 2 | 重构服务层（15个文件） | 16-24小时 |
| 3 | 开发MongoDB→PG迁移工具 | 12-16小时 |
| 4 | 开发PG→MongoDB迁移工具 | 10-14小时 |
| 5 | 实现数据验证 | 6-8小时 |
| 6 | 测试和文档 | 8-12小时 |
| **总计** | | **60-86小时** |

---

## ⚠️ 风险和注意事项

1. **数据一致性风险**
   - 迁移前必须完整备份
   - 建议在测试环境先验证

2. **性能风险**
   - 大数据量迁移需要分批处理
   - 使用批量插入提升性能

3. **ID映射问题**
   - ObjectId和UUID的转换需要妥善处理
   - 保存ID映射表以便回滚

4. **服务层重构风险**
   - 一次只重构一个服务文件
   - 每次重构后运行测试验证

---

## ✅ 验收标准

1. ✅ 数据库接口层完整，支持所有必要的CRUD操作
2. ✅ 服务层不再直接使用原生SQL，完全通过接口操作
3. ✅ 配置切换数据库类型即可运行，无需修改代码
4. ✅ MongoDB→PostgreSQL迁移工具完成，包含数据验证
5. ✅ PostgreSQL→MongoDB迁移工具完成，包含数据验证
6. ✅ 所有测试通过
7. ✅ 完整的使用文档

---

## 📋 执行进度

### 阶段1：完善数据库接口层
- [ ] 1.1 增强接口定义
- [ ] 1.2 完善MongoDB实现
- [ ] 1.3 完善PostgreSQL实现
- [ ] 1.4 通用工具

### 阶段2：重构服务层
- [ ] 2.1 重构UserService
- [ ] 2.2 重构NoteService
- [ ] 2.3 重构所有服务文件

### 阶段3：开发数据迁移工具
- [ ] 3.1 核心迁移框架
- [ ] 3.2 MongoDB → PostgreSQL迁移器
- [ ] 3.3 PostgreSQL → MongoDB迁移器
- [ ] 3.4 表迁移逻辑

### 阶段4：实现数据验证
- [ ] 4.1 验证框架
- [ ] 4.2 数据一致性检查

### 阶段5：配置和文档
- [ ] 5.1 配置文件
- [ ] 5.2 使用文档

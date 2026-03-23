# Leanote 数据库迁移工具使用指南

## 概述

本项目为 Leanote 提供了数据库抽象层和数据迁移工具，支持 MongoDB 和 PostgreSQL 两种数据库，并提供双向迁移能力。

## 功能特性

### 1. 数据库抽象层

统一的数据库接口，支持 MongoDB 和 PostgreSQL，可通过配置切换：

```go
type Database interface {
    // 连接管理
    Initialize(config DatabaseConfig) error
    Close() error
    Ping() error
    IsConnected() bool

    // ID生成
    NewID() string
    IsValidID(id string) bool

    // 基础CRUD操作
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

    // MongoDB兼容接口
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
}
```

### 2. 数据迁移工具

支持 MongoDB 到 PostgreSQL 的数据迁移，包含以下特性：

- 自动 ID 映射（ObjectId → UUID）
- 批量插入提升性能
- 实时进度报告
- 数据完整性验证
- 错误处理和恢复

## 配置说明

### 数据库配置 (`conf/app.conf`)

```ini
# 数据库类型选择: mongodb 或 postgresql
db.type=postgresql

# PostgreSQL:配置
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

## 迁移工具使用

### MongoDB → PostgreSQL 迁移

1. **准备环境**

确保 MongoDB 和 PostgreSQL 都正在运行，并且 PostgreSQL 已经导入了 schema：

```bash
# 导入 PostgreSQL schema
psql -U leanote -d leanote -f database/schema.sql
```

2. **运行迁移工具**

```bash
# 基本用法
go run scripts/migrate_mongo_to_pg.go \
    --mongo-url="mongodb://localhost:27017/leanote" \
    --postgres-url="host=localhost port=5432 user=leanote password=your_password dbname=leanote sslmode=disable"

# 带参数的用法
go run scripts/migrate_mongo_to_pg.go \
    --mongo-url="mongodb://localhost:27017/leanote" \
    --postgres-url="host=localhost port=5432 user=leanote password=your_password dbname=leanote sslmode=disable" \
    --batch-size=1000 \
    --verbose=true \
    --stop-on-error=false
```

3. **参数说明**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--mongo-url` | `mongodb://localhost:27017/leanote` | MongoDB 连接 URL |
| `--postgres-url` | `host=localhost port=5432 user=leanote password= dbname=leanote sslmode=disable` | PostgreSQL 连接 URL |
| `--batch-size` | `1000` | 批量处理大小 |
| `--stop-on-error` | `false` | 遇到错误时是否停止 |
| `--verbose` | `true` | 是否输出详细日志 |

4. **迁移示例**

```bash
# 从 MongoDB 迁移到 PostgreSQL
go run scripts/migrate_mongo_to_pg.go \
    --mongo-url="mongodb://localhost:27017/leanote" \
    --postgres-url="host=localhost port=5432 user=leanote password=leanote123 dbname=leanote sslmode=disable"
```

输出示例：

```
2024/03/23 10:00:00 Starting MongoDB to PostgreSQL migration...
2024/03/23 10:00:00 MongoDB URL: mongodb://localhost:27017/leanote
2024/03/23 10:00:00 PostgreSQL URL: host=localhost port=5432 user=leanote password=xxx dbname=leanote sslmode=disable
2024/03/23 10:00:01 Connected to both databases successfully
2024/03/23 10:00:01 Pre-migration: Creating ID mapping table...
2024/03/23 10:00:02 Starting migration for table: users
2024/03/23 10:00:02 Migrated 100/1500 users
2024/03/23 10:00:03 Migrated 200/1500 users
...
2024/03/23 10:00:10 Migration complete for table users: 1500/1500 records in 8s

=== Migration Summary ===
Total Success: 1500 records
Total Failed: 0 records
✓ Migration completed successfully!
```

## 数据验证

迁移工具包含自动验证功能，会在迁移后验证：

1. 记录数匹配
2. ID 映射完整性
3. 数据一致性

## ID 映射

迁移系统会自动创建 ID 映射表：

```sql
CREATE TABLE id_mapping (
    object_id VARCHAR(24) PRIMARY KEY,  -- MongoDB ObjectId
    uuid VARCHAR(36) NOT NULL UNIQUE,    -- PostgreSQL UUID
    table_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

这个映射表用于：
- 追踪 ObjectId 和 UUID 的对应关系
- 支持后续的数据回滚
- 用于数据引用更新

## 性能优化建议

1. **批量大小**
   - 小数据集：1000-2000
   - 大数据集（百万级）：5000-10000

2. **网络配置**
   - 确保 MongoDB 和 PostgreSQL 在同一网络或高速网络
   - 使用持久化连接

3. **PostgreSQL 优化**
   ```sql
   -- 迁移前禁用自动提交
   BEGIN;

   -- 迁移后重新启用
   COMMIT;
   ```

4. **内存分配**
   - 为迁移工具分配足够的内存
   `GOMAXPROCS=4 go run scripts/migrate_mongo_to_pg.go ...`

## 故障排除

### 常见问题

1. **连接失败**
   ```
   Failed to connect to MongoDB: dial tcp 127.0.0.1:27017: connect: connection refused
   ```
   解决方案：确保 MongoDB 正在运行

2. **权限错误**
   ```
   permission denied for database leanote
   ```
   解决方案：确保 PostgreSQL 用户有足够权限

3. **Schema 不匹配**
   ```
   table "users" does not exist
   ```
   解决方案：先导入 PostgreSQL schema

## 回滚策略

如果迁移失败，可以：

1. 使用备份恢复
2. 删除 PostgreSQL 数据重新导入
3. 修复错误后继续迁移

## 生产环境建议

1. **迁移前**
   - 完整备份 MongoDB
   - 在测试环境验证迁移
   - 准备回滚计划

2. **迁移中**
   - 使用屏幕或 tmux 防止中断
   - 监控系统资源
   - 记录迁移日志

3. **迁移后**
   - 验证数据完整性
   - 运行应用测试
   - 监控性能指标

## 支持

如有问题，请查看：
- PLAN.md - 详细执行计划
- MIGRATION_GUIDE.md - 迁移指南
- GitHub Issues

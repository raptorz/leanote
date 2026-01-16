# MongoDB 到 PostgreSQL 迁移指南

## 概述
本文档说明如何将 Leanote 项目从 MongoDB 迁移到 PostgreSQL。

## 已完成的工作

### 1. 依赖更新
- 从 `go.mod` 中移除了 `gopkg.in/mgo.v2`
- 添加了 `github.com/lib/pq` (PostgreSQL 驱动)
- 添加了 `github.com/google/uuid` (UUID 生成)

### 2. 数据库连接层
- 创建了 `app/db/Postgres.go`，提供 PostgreSQL 连接管理
- 提供了基础的数据库操作函数（Insert, Update, Delete, Get 等）

### 3. 数据模型
- 所有 `app/info/*.go` 文件已更新
- 将 `bson.ObjectId` 替换为 `string` 类型（UUID）
- 更新了所有 struct tags 从 `bson` 改为 `db`

### 4. 数据库 Schema
- 创建了 `database/schema.sql`，包含所有 PostgreSQL 表结构
- 使用 UUID 作为主键类型
- 创建了所有必要的索引

### 5. Service 层示例
- 创建了 `app/service/UserService_pg.go` 作为参考
- 展示了如何将 MongoDB 查询转换为 SQL 查询

## 需要完成的 Service 迁移

以下是需要迁移的 Service 文件列表，按照优先级排序：

### 高优先级（核心功能）
1. `NotebookService.go` - 笔记本管理
2. `NoteService.go` - 笔记管理
3. `NoteContentHistoryService.go` - 笔记历史记录
4. `TagService.go` - 标签管理
5. `ShareService.go` - 分享功能
6. `AttachService.go` - 附件管理
7. `AuthService.go` - 认证功能

### 中优先级
8. `FileService.go` - 文件管理
9. `BlogService.go` - 博客功能
10. `GroupService.go` - 分组功能
11. `SessionService.go` - 会话管理
12. `TokenService.go` - Token 管理
13. `AlbumService.go` - 相册管理
14. `NoteImageService.go` - 笔记图片
15. `ConfigService.go` - 配置管理
16. `ThemeService.go` - 主题管理
17. `EmailService.go` - 邮件服务
18. `SuggestionService.go` - 建议功能

### 低优先级
19. `PwdService.go` - 密码找回
20. `TrashService.go` - 回收站
21. `UpgradeService.go` - 升级服务
22. `common.go` - 通用函数

## Service 迁移步骤

### 1. 理解原有 MongoDB 操作
查看原始 Service 文件，理解以下 MongoDB 操作：
- `db.Insert(collection, data)` - 插入
- `db.Update(collection, query, data)` - 更新
- `db.Delete(collection, query)` - 删除
- `db.Get(collection, id, dest)` - 根据 ID 获取
- `db.GetByQ(collection, query, dest)` - 根据条件获取
- `db.ListByQ(collection, query, dest)` - 列表查询
- `db.Count(collection, query)` - 计数
- `db.Has(collection, query)` - 检查是否存在

### 2. 转换为 PostgreSQL 操作

#### 插入操作
```go
// MongoDB
db.Insert(db.Users, user)

// PostgreSQL
query := `INSERT INTO users (id, email, username, ...) VALUES ($1, $2, $3, ...)`
_, err := db.DB.Exec(query, user.UserId, user.Email, user.Username, ...)
```

#### 查询操作
```go
// MongoDB
var user info.User
db.GetByQ(db.Users, bson.M{"email": email}, &user)

// PostgreSQL
var user info.User
query := `SELECT id, email, username, ... FROM users WHERE email = $1`
err := db.DB.QueryRow(query, email).Scan(&user.UserId, &user.Email, &user.Username, ...)
```

#### 更新操作
```go
// MongoDB
db.Update(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, bson.M{"$set": bson.M{"username": newName}})

// PostgreSQL
query := `UPDATE users SET username = $1 WHERE id = $2`
_, err := db.DB.Exec(query, newName, userId)
```

#### 删除操作
```go
// MongoDB
db.Delete(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)})

// PostgreSQL
query := `DELETE FROM users WHERE id = $1`
_, err := db.DB.Exec(query, userId)
```

#### 列表查询
```go
// MongoDB
var users []info.User
db.ListByQ(db.Users, bson.M{"is_deleted": false}, &users)

// PostgreSQL
query := `SELECT id, email, username, ... FROM users WHERE is_deleted = false`
rows, err := db.DB.Query(query)
if err != nil {
    return nil, err
}
defer rows.Close()

var users []info.User
for rows.Next() {
    var user info.User
    err := rows.Scan(&user.UserId, &user.Email, &user.Username, ...)
    if err != nil {
        return nil, err
    }
    users = append(users, user)
}
```

### 3. 处理特殊类型

#### UUID
```go
// 不再需要 bson.ObjectIdHex() 转换
userId := "123e4567-e89b-12d3-a456-426614174000"

// 生成新 UUID
userId := db.NewUUID()
```

#### 数组字段
```go
// PostgreSQL 使用原生数组类型
query := `SELECT tags FROM notes WHERE id = $1`
err := db.DB.QueryRow(query, noteId).Scan(pq.Array(&tags))

// 插入数组
query := `INSERT INTO notes (tags) VALUES ($1)`
_, err := db.DB.Exec(query, pq.Array([]string{"tag1", "tag2"}))
```

#### JSON 字段
```go
// 对于复杂对象，使用 JSONB
query := `SELECT config FROM users WHERE id = $1`
var config map[string]interface{}
err := db.DB.QueryRow(query, userId).Scan(&config)
```

### 4. 处理事务
```go
tx, err := db.DB.Begin()
if err != nil {
    return err
}

defer func() {
    if err != nil {
        tx.Rollback()
    } else {
        err = tx.Commit()
    }
}()

// 在事务中执行操作
_, err = tx.Exec(query1, args1...)
if err != nil {
    return err
}

_, err = tx.Exec(query2, args2...)
if err != nil {
    return err
}
```

### 5. 处理分页
```go
// MongoDB
db.ListByQWithFields(db.Notes, bson.M{"user_id": userId}, fields, &notes)
db.Coll(db.Notes).Find(query).Limit(pageSize).Skip(pageNum*pageSize).All(&notes)

// PostgreSQL
query := `SELECT id, title, ... FROM notes WHERE user_id = $1 LIMIT $2 OFFSET $3`
rows, err := db.DB.Query(query, userId, pageSize, pageNum*pageSize)
```

## 配置文件更新

在 `conf/app.conf` 中更新数据库配置：

```ini
# PostgreSQL 配置
db.host=localhost
db.port=5432
db.username=leanote
db.password=your_password
db.dbname=leanote

# 或使用连接字符串
# db.url=host=localhost port=5432 user=leanote password=your_password dbname=leanote sslmode=disable
```

## 初始化代码更新

在 `app/service/init.go` 中：

```go
// 将 db.Init() 改为 db.InitPG()
db.InitPG(revel.Config.StringDefault("db.url", ""), revel.Config.StringDefault("db.dbname", "leanote"))
```

## 测试建议

1. **单元测试**：为每个 Service 的函数编写测试
2. **集成测试**：测试完整的业务流程
3. **性能测试**：对比 MongoDB 和 PostgreSQL 的性能
4. **数据迁移测试**：验证从 MongoDB 迁移的数据完整性

## 常见问题

### 1. ObjectId 到 UUID 的转换
MongoDB 的 ObjectId 是 12 字节的十六进制字符串，而 UUID 是 36 字符的标准格式。需要创建一个迁移工具将现有 ObjectId 转换为 UUID。

### 2. 数组和嵌套对象
MongoDB 原生支持数组和嵌套对象，PostgreSQL 需要使用数组类型或 JSONB。

### 3. 查询语法
MongoDB 使用 BSON 查询，PostgreSQL 使用 SQL。需要学习 PostgreSQL 的查询语法。

### 4. 事务处理
MongoDB 4.0 以后支持事务，但使用方式与 PostgreSQL 不同。

## 下一步

1. 按照优先级逐个迁移 Service 文件
2. 创建数据迁移脚本将现有数据从 MongoDB 导出到 PostgreSQL
3. 全面测试迁移后的功能
4. 优化性能，添加必要的索引
5. 更新文档

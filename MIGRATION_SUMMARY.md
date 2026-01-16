# MongoDB 到 PostgreSQL 迁移总结

## 已完成的工作

### 1. 项目依赖更新 ✅
- **文件**: `go.mod`
- **变更**:
  - 移除 `gopkg.in/mgo.v2` (MongoDB 驱动)
  - 添加 `github.com/lib/pq` (PostgreSQL 驱动)
  - 添加 `github.com/google/uuid` (UUID 生成)

### 2. 数据库连接管理 ✅
- **文件**: `app/db/Postgres.go`
- **功能**:
  - PostgreSQL 连接初始化和管理
  - 基础 CRUD 操作函数
  - 连接池配置
  - 健康检查

### 3. 数据模型更新 ✅
- **文件**: `app/info/*.go` (共 21 个文件)
- **变更**:
  - 所有 `bson.ObjectId` 替换为 `string` (UUID)
  - 更新 struct tags 从 `bson` 改为 `db`
  - 修改数据类型以匹配 PostgreSQL

**已更新的文件**:
- UserInfo.go
- NotebookInfo.go
- NoteInfo.go
- TagInfo.go
- AttachInfo.go
- GroupInfo.go
- Configinfo.go
- TokenInfo.go
- SessionInfo.go
- BlogInfo.go
- ThemeInfo.go
- NoteImage.go
- AlbumInfo.go
- FileInfo.go
- ShareNotebookNoteInfo.go
- EmailLogInfo.go
- SuggestionInfo.go
- ReportInfo.go

### 4. 数据库 Schema ✅
- **文件**: `database/schema.sql`
- **内容**:
  - 完整的 PostgreSQL 表结构
  - 所有必要的索引
  - 外键约束
  - UUID 作为主键
  - 支持数组和 JSONB 类型

**创建的表**:
- users
- notebooks
- notes
- note_contents
- note_content_histories
- tags
- tag_counts
- attachs
- share_notebooks
- share_notes
- groups
- group_users
- blogs
- blog_singles
- blog_likes
- blog_comments
- themes
- files
- albums
- note_images
- configs
- sessions
- tokens
- email_logs
- suggestions
- reports

### 5. Service 层示例 ✅
- **文件**: `app/service/UserService_pg.go`
- **功能**:
  - 用户 CRUD 操作
  - 展示 MongoDB 查询到 SQL 的转换模式
  - 作为其他 Service 迁移的参考

### 6. 配置文件更新 ✅
- **文件**: `conf/app.conf`
- **变更**:
  - 添加 PostgreSQL 配置选项
  - 保留 MongoDB 配置（已注释）
  - 支持灵活切换数据库

### 7. 数据迁移脚本 ✅
- **文件**: `scripts/migrate_mongo_to_pg.go`
- **功能**:
  - 从 MongoDB 读取数据
  - 转换 ObjectId 到 UUID
  - 写入 PostgreSQL
  - 支持主要数据类型的迁移

**支持迁移的数据**:
- Users (用户)
- Notebooks (笔记本)
- Notes (笔记)
- Note Contents (笔记内容)
- Tags (标签)
- Attachments (附件)

### 8. 迁移指南 ✅
- **文件**: `MIGRATION_GUIDE.md`
- **内容**:
  - 详细的迁移步骤
  - MongoDB 到 SQL 的转换示例
  - Service 迁移优先级
  - 常见问题和解决方案

## 技术决策

### 1. 主键类型
- **选择**: UUID (string 类型)
- **原因**:
  - 跨数据库兼容性好
  - 不依赖特定数据库特性
  - 便于分布式系统
  - 替换 MongoDB ObjectId 自然

### 2. 数据访问方式
- **选择**: 原生 SQL + lib/pq
- **原因**:
  - 学习曲线平缓
  - 性能可控
  - 不引入额外 ORM 复杂性
  - 便于调试和优化

### 3. 数据迁移策略
- **选择**: 保留 ObjectId 格式作为 UUID
- **原因**:
  - 保持原有数据可追溯性
  - 简化迁移逻辑
  - 减少 URL 变更影响

## 后续工作

### 1. Service 层迁移（需手动完成）
按照优先级迁移以下 Service 文件：

**高优先级**:
- [ ] NotebookService.go
- [ ] NoteService.go
- [ ] NoteContentHistoryService.go
- [ ] TagService.go
- [ ] ShareService.go
- [ ] AttachService.go
- [ ] AuthService.go

**中优先级**:
- [ ] FileService.go
- [ ] BlogService.go
- [ ] GroupService.go
- [ ] SessionService.go
- [ ] TokenService.go
- [ ] AlbumService.go
- [ ] NoteImageService.go
- [ ] ConfigService.go
- [ ] ThemeService.go
- [ ] EmailService.go
- [ ] SuggestionService.go

**低优先级**:
- [ ] PwdService.go
- [ ] TrashService.go
- [ ] UpgradeService.go

### 2. 测试
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能测试
- [ ] 数据迁移验证

### 3. 优化
- [ ] 查询性能优化
- [ ] 索引优化
- [ ] 连接池调优
- [ ] 事务管理优化

### 4. 部署
- [ ] 生产环境配置
- [ ] 数据库备份策略
- [ ] 监控和告警
- [ ] 回滚方案

## 迁移步骤

### 第一步：准备环境
```bash
# 1. 安装 PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# 2. 创建数据库和用户
sudo -u postgres psql
CREATE USER leanote WITH PASSWORD 'your_password';
CREATE DATABASE leanote OWNER leanote;
GRANT ALL PRIVILEGES ON DATABASE leanote TO leanote;
\q

# 3. 创建表结构
psql -U leanote -d leanote -f database/schema.sql

# 4. 更新配置文件
# 编辑 conf/app.conf，设置 PostgreSQL 连接信息
```

### 第二步：迁移数据
```bash
# 1. 编译迁移脚本
go build -o migrate_mongo_to_pg scripts/migrate_mongo_to_pg.go

# 2. 运行迁移
./migrate_mongo_to_pg

# 3. 验证数据
psql -U leanote -d leanote -c "SELECT COUNT(*) FROM users;"
```

### 第三步：切换应用
```bash
# 1. 安装新依赖
go mod tidy
go mod download

# 2. 更新初始化代码
# 在 app/service/init.go 中，将 db.Init() 改为 db.InitPG()

# 3. 编译并运行
revel run github.com/leanote/leanote
```

## 注意事项

### 1. 数据一致性
- 确保迁移过程中 MongoDB 不再有写入操作
- 验证迁移后的数据完整性
- 准备回滚方案

### 2. 性能考虑
- PostgreSQL 的查询优化与 MongoDB 不同
- 需要根据实际查询模式调整索引
- 大数据量迁移可能需要分批处理

### 3. 兼容性
- UUID 替换 ObjectId 后，URL 可能需要更新
- 前端代码可能需要适配 UUID 格式
- 现有的 API 可能需要调整

### 4. 安全性
- 更新数据库连接密码
- 配置 SSL 连接
- 定期备份数据库

## 常见问题

### Q1: 如何处理现有的 ObjectId？
A: 迁移脚本会将 ObjectId 转换为 UUID 格式。如果需要保持与旧系统的兼容性，可以在表中保留一个 legacy_object_id 字段。

### Q2: 数组字段如何处理？
A: PostgreSQL 原生支持数组类型。在 schema 中使用 `TEXT[]` 或 `INTEGER[]` 等类型。

### Q3: 如何处理嵌套对象？
A: 使用 PostgreSQL 的 JSONB 类型存储复杂对象，或者将嵌套对象拆分为关联表。

### Q4: 性能会有什么变化？
A: PostgreSQL 在复杂查询和事务处理上通常优于 MongoDB，但在简单的文档读写上可能稍慢。建议进行性能测试并优化查询。

### Q5: 如何回滚到 MongoDB？
A: 保留 MongoDB 数据库和配置，如需回滚只需切换回 MongoDB 配置并重启服务。

## 支持和反馈

如有问题，请参考以下资源：
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - 详细迁移指南
- [database/schema.sql](database/schema.sql) - 数据库结构
- [app/service/UserService_pg.go](app/service/UserService_pg.go) - Service 迁移示例

## 版本历史

- v1.0.0 - 初始版本，完成基础迁移框架
  - 数据库连接层
  - 数据模型更新
  - Schema 设计
  - 迁移工具

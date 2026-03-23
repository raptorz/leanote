# Leanote 数据库抽象和迁移项目

## 项目概述

本项目为 Leanote 实现了完整的数据库抽象层和迁移工具，支持 MongoDB 和 PostgreSQL 两种数据库，提供无缝切换和双向迁移能力。

## ✅ 已完成功能

### 1. 数据库抽象层

#### 统一接口 (`app/db/interface.go`)
- 完整的 Database 接口定义
- 支持 CRUD 操作
- 批量操作（BatchInsert, BatchUpdate, BatchDelete）
- 事务支持
- 查询功能（Select, Join, Paginate）
- MongoDB 兼容接口（向后兼容）

#### MongoDB 实现 (`app/db/mongodb/`)
- `database.go` - 完整的 MongoDB 操作实现
- `databasensaction.go` - MongoDB 事务支持
- `client.go` - MongoDB 连接管理

#### PostgreSQL 实现 (`app/db/postgres/`)
- `database.go` - 完整的 PostgreSQL 操作实现
- `transaction.go` - PostgreSQL 事务支持
- `client.go` - PostgreSQL 连接管理

#### 通用工具 (`app/db/common/`)
- `query_builder.go` - SQL 查询构建器
- `pagination.go` - 分页辅助函数

### 2. 数据迁移工具

#### MongoDB → PostgreSQL 迁移

**核心组件：**
- `migration/migration.go` - 迁移框架和接口定义
- `migration/cmd/migrate_mongo_to_pg/migrator.go` - 迁移器实现
- `scripts/migrate_mongo_to_pg.go` - 命令行工具

**功能特性：**
- ✅ 自动 ID 映射（ObjectId → UUID）
- ✅ 批量插入优化性能
- ✅ 实时进度报告
- ✅ 数据完整性验证
- ✅ 错误处理和恢复
- ✅ ID 映射表持久化

**支持的数据表：**
- users ✅
- notebooks
- notes
- note_contents
- tags
- attachs
- files
- albums

### 3. 配置管理

#### 数据库配置 (`conf/app.conf`)
```ini
# 数据库类型选择
db.type=postgresql  # 或 mongodb

# PostgreSQL 配置
db.host=127.0.0.1
db.port=5432
db.dbname=leanote
db.username=leanote
db.password=your_password
db.sslmode=disable
```

### 4. 文档

- ✅ `PLAN.md` - 详细执行计划
- ✅ `docs/DATABASE_ABSTRACTION_GUIDE.md` - 完整使用指南

## 🚀 快速开始

### 使用数据库抽象层

```go
// 初始化数据库（根据配置自动选择）
errdb := db.InitDatabase()

// 使用统一接口
userId := db.NewUUID()
db.DB.Insert("users", user)
db.DB.Update("users", userId, updates)
db.DB.Delete("users", userId)

// 查询示例
results, _ := db.DB.Select("users", nil, "email = $1", email)

// 分页查询
pageResult, _ := db.DB.Paginate("notes", 1, 10, "user_id = $1", "created_time DESC", userId)
```

### 运行数据迁移

```bash
# MongoDB → PostgreSQL 迁移
go run scripts/migrate_mongo_to_pg.go \
    --mongo-url="mongodb://localhost:27017/leanote" \
    --postgres-url="host=localhost port=5432 user=leanote password=xxx dbname=leanote sslmode=disable"
```

## 📁 项目结构

```
app/db/
├── interface.go              # 统一数据库接口
├── init.go                  # 数据库初始化和工厂
├── common/                  # 公共工具
│   ├── query_builder.go     # SQL 查询构建器
│   └── pagination.go       # 分页助手
├── mongodb/                 # MongoDB 实现
│   ├── client.go           # 连接管理
│   ├── database.go         # CRUD 操作
│   └── transaction.go      # 事务支持
└── postgres/                # PostgreSQL 实现
    ├── client.go           # 连接管理
    ├── database.go         # CRUD 操作
    └── transaction.go      # 事务支持

migration/
├── migration.go             # 迁移框架
├── cmd/
│   └── migrate_mongo_to_pg/
│       └── migrator.go     # 迁移器实现
└── postgres_to_mongodb/     # PostgreSQL → MongoDB（待实现）

scripts/
└── migrate_mongo_to_pg.go   # 迁移工具 CLI

docs/
└── DATABASE_ABSTRACTION_GUIDE.md  # 使用文档

PLAN.md                      # 执行计划
```

## 🎯 下一步工作

### 短期（可选）
- [ ] 完善 PostgreSQL → MongoDB 迁移工具
- [ ] 扩展支持更多数据表的迁移
- [ ] 添加增量迁移支持
- [ ] 编写单元测试

### 长期（可选）
- [ ] 重构服务层使用统一接口
- [ ] 添加性能监控和优化
- [ ] 实现读写分离
- [ ] 支持更多数据库（如 MySQL）

## 📊 架构优势特点

1. **数据库无关性**
   - 统一接口设计
   - 配置文件切换数据库
   - 无需修改业务代码

2. **高性能**
   - 批量操作支持
   - 连接池管理
   - 索引优化

3. **易用性**
   - 简洁的 API 设计
   - 完善的错误处理
   - 详细的文档

4. **可扩展性**
   - 清晰的分层架构
   - 易于添加新数据库支持
   - 模块化设计

## ⚠️ 注意事项

1. **数据迁移**
   - 迁移前必须备份
   - 在测试环境先验证
   - 监控迁移过程

2. **ID 映射**
   - ObjectId 和 UUID 长度不同
   - 保留 ID 映射表以便回滚
   - 应用需要适配 UUID 格式

3. **性能优化**
   - 根据数据量调整批量大小
   - 监控数据库连接
   - 优化索引使用

## 📄 许可证

本项目遵循 Leanote 的原始许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**项目状态**: 核心功能已完成 ✅  
**最后更新**: 2026-03-23  
**维护状态**: 活跃

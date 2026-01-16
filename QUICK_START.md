# 快速开始指南

## 5 分钟快速迁移

### 前提条件
- 已安装 Go 1.15+
- 已安装 PostgreSQL 12+
- 已安装 MongoDB（仅用于数据迁移）

### 步骤 1: 安装 PostgreSQL

```bash
# Ubuntu/Debian
sudo apt-get update && sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql

# 启动服务
sudo systemctl start postgresql  # Linux
brew services start postgresql    # macOS
```

### 步骤 2: 创建数据库

```bash
sudo -u postgres psql << EOF
CREATE USER leanote WITH PASSWORD 'leanote123';
CREATE DATABASE leanote OWNER leanote;
GRANT ALL PRIVILEGES ON DATABASE leanote TO leanote;
\q
EOF
```

### 步骤 3: 创建表结构

```bash
psql -U leanote -d leanote -f database/schema.sql
```

### 步骤 4: 配置应用

编辑 `conf/app.conf`，确保 PostgreSQL 配置正确：

```ini
db.host=127.0.0.1
db.port=5432
db.dbname=leanote
db.username=leanote
db.password=leanote123
```

### 步骤 5: 更新代码（可选）

如果你要从 MongoDB 迁移数据：

```bash
# 编辑 app/service/init.go
# 将 db.Init() 改为 db.InitPG()
```

### 步骤 6: 运行应用

```bash
# 开发模式
revel run github.com/leanote/leanote

# 或使用 PostgreSQL 模式
 revel run github.com/leanote/leanote --db=postgres
```

### 步骤 7: 访问应用

打开浏览器访问: http://localhost:9000

## 从 MongoDB 迁移数据

### 1. 确保 MongoDB 运行

```bash
mongod --dbpath /path/to/mongodb/data
```

### 2. 运行迁移脚本

```bash
cd scripts
go run migrate_mongo_to_pg.go
```

迁移脚本会自动：
- 连接到 MongoDB
- 读取所有用户、笔记本、笔记等数据
- 转换 ObjectId 到 UUID
- 写入 PostgreSQL

### 3. 验证迁移

```bash
psql -U leanote -d leanote << EOF
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Notebooks', COUNT(*) FROM notebooks
UNION ALL
SELECT 'Notes', COUNT(*) FROM notes;
EOF
```

## 验证功能

### 1. 测试用户注册

1. 访问 http://localhost:9000
2. 点击"注册"
3. 填写表单
4. 检查数据库是否有新用户

```bash
psql -U leanote -d leanote -c "SELECT username, email, created_time FROM users ORDER BY created_time DESC LIMIT 1;"
```

### 2. 测试笔记创建

1. 登录系统
2. 创建新笔记
3. 检查数据库

```bash
psql -U leanote -d leanote -c "SELECT title, created_time FROM notes ORDER BY created_time DESC LIMIT 1;"
```

### 3. 测试 CRUD 操作

```bash
# 查看所有用户
psql -U leanote -d leanote -c "SELECT username, email FROM users LIMIT 10;"

# 更新用户信息
psql -U leanote -d leanote -c "UPDATE users SET logo = 'new_logo.png' WHERE username = 'testuser';"

# 删除测试数据
psql -U leanote -d leanote -c "DELETE FROM users WHERE username = 'testuser';"
```

## 常见问题快速解决

### 问题: 连接被拒绝

```bash
# 检查 PostgreSQL 是否运行
sudo systemctl status postgresql

# 启动 PostgreSQL
sudo systemctl start postgresql
```

### 问题: 密码认证失败

```bash
# 重置密码
sudo -u postgres psql
ALTER USER leanote WITH PASSWORD 'new_password';
\q

# 更新配置文件中的密码
```

### 问题: 权限不足

```bash
# 授予权限
sudo -u postgres psql
GRANT ALL PRIVILEGES ON DATABASE leanote TO leanote;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO leanote;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO leanote;
\q
```

## 下一步

1. **完成 Service 迁移**: 按照 [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) 迁移剩余的 Service 文件

2. **性能优化**: 参考 [DEPLOYMENT.md](DEPLOYMENT.md) 优化数据库配置

3. **生产部署**: 配置 SSL、备份、监控等生产环境设置

4. **测试**: 进行完整的集成测试和性能测试

## 文件位置

- **数据库 Schema**: `database/schema.sql`
- **数据库连接**: `app/db/Postgres.go`
- **Service 示例**: `app/service/UserService_pg.go`
- **迁移脚本**: `scripts/migrate_mongo_to_pg.go`
- **配置文件**: `conf/app.conf`

## 获取帮助

- 详细迁移指南: [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)
- 部署指南: [DEPLOYMENT.md](DEPLOYMENT.md)
- 迁移总结: [MIGRATION_SUMMARY.md](MIGRATION_SUMMARY.md)

## 快速命令参考

```bash
# 数据库操作
psql -U leanote -d leanote                           # 连接数据库
psql -U leanote -d leanote -f database/schema.sql      # 导入 schema
pg_dump -U leanote -d leanote > backup.sql            # 备份
psql -U leanote -d leanote < backup.sql               # 恢复

# 应用操作
go mod download                                        # 下载依赖
revel run github.com/leanote/leanote                  # 运行应用
revel build github.com/leanote/leanote leanote         # 编译应用

# 迁移操作
go run scripts/migrate_mongo_to_pg.go                 # 运行迁移

# 监控操作
psql -U leanote -d leanote -c "SELECT count(*) FROM pg_stat_activity WHERE datname = 'leanote';"  # 查看连接数
```

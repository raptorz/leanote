# 安装和部署指南

## 环境要求

- Go 1.15+
- PostgreSQL 12+
- MongoDB (仅用于数据迁移)

## 安装步骤

### 1. 安装 PostgreSQL

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
```

#### macOS
```bash
brew install postgresql
```

#### Windows
下载并安装 [PostgreSQL for Windows](https://www.postgresql.org/download/windows/)

### 2. 配置 PostgreSQL

```bash
# 切换到 postgres 用户
sudo -u postgres psql

# 创建数据库和用户
CREATE USER pearlnote WITH PASSWORD 'your_secure_password';
CREATE DATABASE pearlnote OWNER pearlnote;
GRANT ALL PRIVILEGES ON DATABASE pearlnote TO pearlnote;
\q
```

### 3. 创建数据库表结构

```bash
# 导入 schema
psql -U pearlnote -d pearlnote -f database/schema.sql

# 验证表是否创建成功
psql -U pearlnote -d pearlnote -c "\dt"
```

### 4. 配置应用

编辑 `conf/app.conf`:

```ini
# 数据库配置
db.host=127.0.0.1
db.port=5432
db.dbname=pearlnote
db.username=pearlnote
db.password=your_secure_password

# 或使用连接字符串
db.url=host=127.0.0.1 port=5432 user=pearlnote password=your_secure_password dbname=pearlnote sslmode=disable
```

## 数据迁移（从 MongoDB）

### 1. 安装 MongoDB

#### Ubuntu/Debian
```bash
sudo apt-get install mongodb
```

#### macOS
```bash
brew install mongodb-community
```

### 2. 导出 MongoDB 数据

```bash
# 确保 MongoDB 正在运行
mongod --dbpath /path/to/mongodb/data

# 备份数据
mongodump --host localhost --port 27017 --db pearlnote --out mongodb_backup
```

### 3. 运行迁移脚本

```bash
# 编译迁移工具
go build -o migrate ./tools/migration

# 运行迁移（确保 MongoDB 和 PostgreSQL 都在运行）
./migrate -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/pearlnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=pearlnote password=pearlnote dbname=pearlnote sslmode=disable'

# 或者直接运行
go run ./tools/migration -direction mongo_to_pg
```

### 4. 验证迁移

```bash
# 连接到 PostgreSQL
psql -U pearlnote -d pearlnote

# 检查数据
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM notebooks;
SELECT COUNT(*) FROM notes;

\q
```

## 应用部署

### 1. 安装依赖

```bash
# 下载依赖
go mod download

# 如果有新的依赖
go mod tidy
```

### 2. 更新初始化代码

编辑 `app/service/init.go`:

```go
// 将原来的 MongoDB 初始化改为 PostgreSQL
// db.Init(revel.Config.StringDefault("db.url", ""), revel.Config.StringDefault("db.dbname", "pearlnote"))

// 改为
db.InitPG(revel.Config.StringDefault("db.url", ""), revel.Config.StringDefault("db.dbname", "pearlnote"))
```

### 3. 编译应用

```bash
# 开发模式
revel run github.com/pearlnote/pearlnote

# 生产模式编译
revel build github.com/pearlnote/pearlnote pearlnote
```

### 4. 运行应用

```bash
# 开发环境
revel run github.com/pearlnote/pearlnote

# 生产环境
./pearlnote/run.sh
```

## 验证部署

### 1. 检查应用日志

```bash
# 查看日志
tail -f logs/app.log

# 检查是否有数据库连接错误
grep "database" logs/app.log
```

### 2. 测试基本功能

1. 访问 `http://localhost:9000`
2. 注册新用户
3. 创建笔记
4. 测试 CRUD 操作

### 3. 监控数据库连接

```bash
# 查看当前连接
psql -U pearlnote -d pearlnote -c "SELECT count(*) FROM pg_stat_activity WHERE datname = 'pearlnote';"

# 查看慢查询
psql -U pearlnote -d pearlnote -c "SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;"
```

## 性能优化

### 1. 数据库配置

编辑 `postgresql.conf`:

```ini
# 内存配置
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB

# 连接配置
max_connections = 100

# 查询优化
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4MB
```

重启 PostgreSQL:
```bash
sudo systemctl restart postgresql
```

### 2. 应用配置

编辑 `app/db/Postgres.go`:

```go
DB.SetMaxOpenConns(50)
DB.SetMaxIdleConns(25)
DB.SetConnMaxLifetime(10 * time.Minute)
```

### 3. 创建额外索引

根据实际查询模式添加索引：

```sql
-- 示例：为常用查询添加索引
CREATE INDEX idx_notes_user_notebook ON notes(user_id, notebook_id);
CREATE INDEX idx_note_contents_updated_time ON note_contents(updated_time DESC);
```

## 备份和恢复

### 备份

```bash
# 完整备份
pg_dump -U pearlnote -d pearlnote > backup_$(date +%Y%m%d_%H%M%S).sql

# 仅备份数据
pg_dump -U pearlnote -d pearlnote --data-only > data_backup.sql

# 仅备份结构
pg_dump -U pearlnote -d pearlnote --schema-only > schema_backup.sql
```

### 恢复

```bash
# 恢复完整备份
psql -U pearlnote -d pearlnote < backup_20240116_120000.sql

# 恢复到新数据库
createdb -U pearlnote pearlnote_new
psql -U pearlnote -d pearlnote_new < backup_20240116_120000.sql
```

## 监控和日志

### 1. PostgreSQL 日志

编辑 `postgresql.conf`:

```ini
logging_collector = on
log_directory = 'pg_log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_statement = 'all'
log_duration = on
```

### 2. 应用监控

添加 Prometheus 监控:

```go
import "github.com/prometheus/client_golang/prometheus"

var dbConnections = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "db_connections",
    Help: "Number of active database connections",
})
```

## 故障排除

### 问题 1: 连接被拒绝

```
Error: connection refused
```

**解决方案**:
```bash
# 检查 PostgreSQL 是否运行
sudo systemctl status postgresql

# 检查端口
netstat -an | grep 5432

# 检查防火墙
sudo ufw status
```

### 问题 2: 认证失败

```
Error: password authentication failed
```

**解决方案**:
```bash
# 重置密码
sudo -u postgres psql
ALTER USER pearlnote WITH PASSWORD 'new_password';
\q

# 更新 conf/app.conf
```

### 问题 3: 性能问题

**解决方案**:
```sql
-- 查看慢查询
SELECT query, mean_exec_time, calls 
FROM pg_stat_statements 
ORDER BY mean_exec_time DESC 
LIMIT 10;

-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 分析表
ANALYZE users;
ANALYZE notebooks;
ANALYZE notes;
```

### 问题 4: 迁移失败

**解决方案**:
```bash
# 检查 MongoDB 连接
mongo --host localhost --port 27017 --eval "db.version()"

# 检查 PostgreSQL 连接
psql -U pearlnote -d pearlnote -c "SELECT version();"

# 查看迁移日志
./migrate 2>&1 | tee migration.log
```

## 安全建议

### 1. 数据库安全

- 使用强密码
- 限制远程访问
- 启用 SSL 连接
- 定期更新 PostgreSQL

### 2. 应用安全

- 使用环境变量存储敏感信息
- 启用 HTTPS
- 实施访问控制
- 定期备份数据

### 3. 网络安全

- 配置防火墙
- 使用 VPN 或 SSH 隧道
- 限制数据库端口访问
- 实施 IP 白名单

## 生产环境检查清单

- [ ] 使用强密码
- [ ] 配置 SSL/TLS
- [ ] 设置定期备份
- [ ] 配置监控和告警
- [ ] 优化数据库性能
- [ ] 配置负载均衡
- [ ] 设置日志轮转
- [ ] 实施灾难恢复计划
- [ ] 进行压力测试
- [ ] 文档更新

## 支持

如有问题，请参考：
- [MIGRATION_SUMMARY.md](MIGRATION_SUMMARY.md)
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)

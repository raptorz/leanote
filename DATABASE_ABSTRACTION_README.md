# Leanote 数据库抽象

数据库后端由 `db.type` 选择，可配置为 `mongodb` 或 `postgresql`。业务 service 不直接使用 `mgo.Collection`，统一通过 `app/db` 的 CRUD、条件查询、排序、分页和投影接口访问数据库。

## ID 约定

两个后端都使用 24 位 MongoDB ObjectId。应用模型继续使用 `bson.ObjectId`，PostgreSQL 适配层负责在 `bson.ObjectId` 与 `CHAR(24)` 之间转换。

这样可以保持现有 API、客户端同步协议和 MongoDB 数据兼容，也使双向迁移无需维护 UUID 映射。

## 查询兼容层

PostgreSQL 适配器支持当前业务使用的 BSON 条件和更新操作，包括：

- `$or`、`$and`；
- `$gt`、`$gte`、`$lt`、`$lte`、`$ne`；
- `$in`、`$nin`、`$all`、`$exists`、`$regex`；
- `$set`、`$inc`、`$push`、`$pull`、`$addToSet`；
- 排序、offset、limit 和字段投影。

未知操作符会生成恒假条件，不会退化为全表更新或删除。

## PostgreSQL schema

DDL 位于 `database/schema.sql`，表和列与 `app/info` 中的持久化模型保持一致。`docker-compose.postgres.yml` 会在新的 PostgreSQL 数据目录中自动执行该文件。

已有 PostgreSQL 数据库升级时，应先备份，再显式执行经过审核的 schema migration；容器初始化脚本只对空数据目录生效。

## 数据迁移

双向迁移统一使用：

```bash
go run ./tools/migration -direction mongo_to_pg ...
go run ./tools/migration -direction pg_to_mongo ...
```

完整参数、备份和校验流程见 [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)。

## 测试

普通测试不依赖本机数据库：

```bash
go test ./...
```

PostgreSQL CRUD 集成测试：

```bash
LEANOTE_INTEGRATION_POSTGRES_URL='host=127.0.0.1 port=5432 user=leanote password=leanote dbname=leanote sslmode=disable' \
  go test ./app/db -run TestPostgresCRUDContract -v
```

双向迁移往返测试：

```bash
LEANOTE_INTEGRATION_MONGO_URL='mongodb://127.0.0.1:27017/leanote' \
LEANOTE_INTEGRATION_POSTGRES_URL='host=127.0.0.1 port=5432 user=leanote password=leanote dbname=leanote sslmode=disable' \
  go test ./tools/migration -run TestRoundTripMigration -v
```

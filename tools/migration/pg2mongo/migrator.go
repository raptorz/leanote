package pg2mongo

import (
	"log"

	"app/db"
	"tools/migration/core"
)

type PG2MongoMigrator struct {
	config *core.MigrationConfig
	stats  *db.MigrationStats
}

func NewMigrator(config *core.MigrationConfig) core.Migrator {
	return &PG2MongoMigrator{
		config: config,
		stats:  &db.MigrationStats{},
	}
}

func (m *PG2MongoMigrator) Migrate() (*db.MigrationStats, error) {
	log.Println("开始 PostgreSQL 到 MongoDB 迁移")
	log.Printf("源数据库: %s\n", m.config.Target.Database)
	log.Printf("目标数据库: %s\n", m.config.Source.Database)
	log.Printf("试运行: %v\n", m.config.DryRun)

	m.stats.TotalTables = len(core.GetTableMigrationOrder())

	log.Println("反向迁移功能暂未实现，请使用 MongoDB 到 PostgreSQL 迁移")
	m.stats.SuccessTables = 0
	m.stats.FailedTables = m.stats.TotalTables

	return m.stats, nil
}

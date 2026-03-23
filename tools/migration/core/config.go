package core

import (
	"encoding/json"
	"io/ioutil"

	"app/db/common"
)

// MigrationConfig 迁移配置
type MigrationConfig struct {
	Direction   string // "mongo_to_pg" or "pg_to_mongo"
	Source      common.Config
	Target      common.Config
	BatchSize   int
	DryRun      bool
	MappingFile string
}

func LoadConfig(path string) (*MigrationConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config MigrationConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.BatchSize == 0 {
		config.BatchSize = 1000
	}
	if config.MappingFile == "" {
		config.MappingFile = "id_mappings.json"
	}

	return &config, nil
}

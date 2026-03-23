package common

import (
	"fmt"
)

// ParseConfig 解析数据库配置
type Config struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
	URL      string
}

func ParseConfigFromMap(configMap map[string]interface{}) (*Config, error) {
	config := &Config{}

	if typ, ok := configMap["type"].(string); ok {
		config.Type = typ
	}
	if host, ok := configMap["host"].(string); ok {
		config.Host = host
	}
	if port, ok := configMap["port"].(int); ok {
		config.Port = port
	}
	if port, ok := configMap["port"].(float64); ok {
		config.Port = int(port)
	}
	if username, ok := configMap["username"].(string); ok {
		config.Username = username
	}
	if password, ok := configMap["password"].(string); ok {
		config.Password = password
	}
	if database, ok := configMap["database"].(string); ok {
		config.Database = database
	}
	if dbname, ok := configMap["dbname"].(string); ok {
		config.Database = dbname
	}
	if sslmode, ok := configMap["sslmode"].(string); ok {
		config.SSLMode = sslmode
	}
	if url, ok := configMap["url"].(string); ok {
		config.URL = url
	}

	if config.Type == "" {
		if config.URL != "" {
			if len(config.URL) > 10 && config.URL[:10] == "mongodb://" {
				config.Type = "mongodb"
			} else {
				config.Type = "postgresql"
			}
		} else {
			config.Type = "postgresql" // 默认使用PostgreSQL
		}
	}

	return config, nil
}

// BuildURL 构建数据库连接URL
func BuildURL(config *Config) string {
	if config.URL != "" {
		return config.URL
	}

	switch config.Type {
	case "mongodb":
		if config.Username != "" && config.Password != "" {
			return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
				config.Username, config.Password, config.Host, config.Port, config.Database)
		}
		return fmt.Sprintf("mongodb://%s:%d/%s", config.Host, config.Port, config.Database)
	case "postgresql":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)
	default:
		return ""
	}
}

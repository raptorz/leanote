package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/db/common"
	"github.com/leanote/leanote/app/lea"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
)

type PostgresDatabase struct {
	db     *sql.DB
	config common.Config
	idGen  common.IDGenerator
}

func (p *PostgresDatabase) Initialize(config db.DatabaseConfig) error {
	var ok bool
	var url string

	if config.URL == "" {
		url, ok = revel.Config.String("db.url")
		if !ok {
			url, ok = revel.Config.String("db.urlEnv")
			if ok {
				lea.Log("get db conf from urlEnv: " + url)
			}
		} else {
			lea.Log("get db conf from db.url: " + url)
		}

		if ok {
			urls := strings.Split(url, "/")
			p.config.Database = urls[len(urls)-1]

			if strings.Contains(p.config.Database, "?") {
				urls = strings.Split(p.config.Database, "?")
				p.config.Database = urls[0]
			}
		}
	}
	if p.config.Database == "" {
		p.config.Database, _ = revel.Config.String("db.dbname")
	}

	if !ok {
		portStr, _ := revel.Config.String("db.port")
		if portStr != "" {
			port, _ := strconv.Atoi(portStr)
			p.config.Port = port
		} else {
			p.config.Port = 5432
		}
		p.config.Host, _ = revel.Config.String("db.host")
		p.config.Username, _ = revel.Config.String("db.username")
		p.config.Password, _ = revel.Config.String("db.password")

		url = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			p.config.Host, p.config.Port, p.config.Username, p.config.Password, p.config.Database)
	}

	lea.Log(url)

	var err error
	p.db, err = sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	p.db.SetMaxOpenConns(25)
	p.db.SetMaxIdleConns(25)
	p.db.SetConnMaxLifetime(5 * time.Minute)

	err = p.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	p.idGen = &common.UUIDGenerator{}

	lea.Log("Connected to PostgreSQL database successfully")
	return nil
}

func (p *PostgresDatabase) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDatabase) Ping() error {
	return p.db.Ping()
}

func (p *PostgresDatabase) IsConnected() bool {
	return p.db != nil && p.Ping() == nil
}

func (p *PostgresDatabase) CheckConnection() {
	err := p.db.Ping()
	if err != nil {
		lea.Log("Lost connection to database!")
		p.db.Close()
		err = p.db.Ping()
		if err == nil {
			lea.Log("Reconnect to database successful.")
		} else {
			lea.Log("Reconnect failed!!!! Warning")
		}
	}
}

func (p *PostgresDatabase) NewID() string {
	return p.idGen.Generate()
}

func (p *PostgresDatabase) IsValidID(id string) bool {
	return p.idGen.IsValid(id)
}

func (p *PostgresDatabase) QueryRow(query string, args ...interface{}) db.Row {
	return p.db.QueryRow(query, args...)
}

func (p *PostgresDatabase) Query(query string, args ...interface{}) (db.Rows, error) {
	return p.db.Query(query, args...)
}

func (p *PostgresDatabase) Exec(query string, args ...interface{}) (db.Result, error) {
	return p.db.Exec(query, args...)
}

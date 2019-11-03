package rorm

import (
	"database/sql"
	"errors"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func (re *Engine) SetDB(db *sql.DB) {
	re.db = db
}

func (re *Engine) GetDB() *sql.DB {
	return re.db
}

func (re *Engine) GetPreparedValues() []interface{} {
	return re.preparedValue
}

// // New - init new RORM Engine
// func New(dbDriver, connectionURL string, tbPrefix ...string) *Engine {
// 	re := &Engine{}
// 	if err := re.Connect(dbDriver, connectionURL, tbPrefix...); err != nil {
// 		return nil
// 	}
// 	return re
// }

// New - init new RORM Engine
func New(cfg *DbConfig) (*Engine, error) {
	re := &Engine{
		config:  cfg,
		options: &DbOptions{},
	}
	connStr, err := generateConnectionString(cfg)
	if err != nil {
		return nil, err
	}
	if err := re.Connect(cfg.Driver, connStr); err != nil {
		log.Println("Cannot Connect to DB: ", err.Error())
		return nil, err
	}
	// set default for table case format
	re.options.tbFormat = "snake"
	log.Println("Successful Connect to DB")
	return re, nil
}

func generateConnectionString(cfg *DbConfig) (connectionString string, err error) {
	if strings.TrimSpace(cfg.Driver) == "" || strings.TrimSpace(cfg.Host) == "" || strings.TrimSpace(cfg.Username) == "" || strings.TrimSpace(cfg.DbName) == "" {
		err = errors.New("Config is not set correctly")
		return
	}
	dbURL := &url.URL{
		Scheme: cfg.Driver,
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Host,
	}

	query := url.Values{}
	dbURL.Path = cfg.DbName
	switch cfg.Driver {
	case "postgres":
		query.Add("sslmode", "disable")
		query.Add("search_path", "public")
		if cfg.DbScheme != "" {
			query.Set("search_path", cfg.DbScheme)
		}
		dbURL.Host += ":5432"
		if cfg.Port != "" {
			dbURL.Host = cfg.Host + ":" + cfg.Port
		}
	case "mysql":
		if cfg.Protocol == "" {
			cfg.Protocol = "tcp"
		}
		dbURL.Host += ":3306"
		if cfg.Port != "" {
			dbURL.Host = cfg.Host + ":" + cfg.Port
		}
		replacedStr := "@" + cfg.Protocol + "($1:$2)/"
		rgx := regexp.MustCompile(`@([a-zA-Z0-9]+):([0-9]+)/`)
		res := rgx.ReplaceAllString(dbURL.String(), replacedStr)
		res = res[8:]
		connectionString = res
		return
	case "sqlserver":
		dbURL.Path = ""
		if cfg.DbInstance != "" {
			dbURL.Path = cfg.DbInstance
		}
		query.Add("database", cfg.DbName)
	}
	dbURL.RawQuery = query.Encode()
	connectionString = dbURL.String()
	return
}

// Connect - connect to db Driver
func (re *Engine) Connect(dbDriver, connectionURL string) error {
	var err error
	re.db, err = sql.Open(dbDriver, connectionURL)
	return err
}

func (re *Engine) clearField() {
	re.condition = ""
	re.column = ""
	re.orderBy = ""
	re.tableName = ""
	re.limit = ""
	re.join = ""
	re.isRaw = false
	re.isBulk = false
	re.preparedValue = nil
	re.multiPreparedValue = nil
	re.counter = 0
	re.bulkCounter = 0
}

func (re *Engine) SetTableOptions(tbCaseFormat, tbPrefix string) {
	re.options.tbFormat = tbCaseFormat
	re.options.tbPrefix = tbPrefix
}

func (re *Engine) adjustPreparedParam(old string) string {
	if strings.TrimSpace(re.config.Driver) != "mysql" {

		replacement := ""
		switch re.config.Driver {
		case "postgres":
			replacement = POSTGRES_PREPARED_PARAM
		case "sqlserver":
			replacement = MSSQL_PREPARED_PARAM
		default:
			replacement = ORACLE_PREPARED_PARAM
		}

		for {
			if strings.Index(old, "?") == -1 {
				break
			}
			re.counter++
			old = strings.Replace(old, "?", replacement+strconv.Itoa(re.counter), 1)
		}
	}
	return old
}

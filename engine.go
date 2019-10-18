package rorm

import (
	"database/sql"
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

// New2 - init new RORM Engine
func New(cfg *DbConfig) *Engine {
	re := &Engine{
		config:  cfg,
		options: &DbOptions{},
	}
	driver, connStr := generateConnectionString(cfg)
	if err := re.Connect(driver, connStr); err != nil {
		log.Println("Error When Connect to DB: ", err.Error())
		return nil
	}
	// set default for table case format
	re.options.tbFormat = "snake"
	log.Println("Successful Connect to DB")
	return re
}

func generateConnectionString(cfg *DbConfig) (dbDriver string, connectionString string) {
	dbDriver = cfg.Driver
	dbURL := &url.URL{
		Scheme: cfg.Driver,
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Host,
	}
	if cfg.Port != "" {
		dbURL.Host += ":" + cfg.Port
	} else {
		log.Fatalln("Missing Port in RORM Config")
	}

	query := url.Values{}
	dbURL.Path = cfg.DbName
	switch dbDriver {
	case "postgres":
		query.Add("sslmode", "disable")
		query.Add("search_path", "public")
		if cfg.DbScheme != "" {
			query.Set("search_path", cfg.DbScheme)
		}
	case "mysql":
		if cfg.Protocol == "" {
			cfg.Protocol = "tcp"
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
	re.config.Driver = dbDriver
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

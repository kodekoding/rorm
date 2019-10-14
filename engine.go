package rorm

import (
	"database/sql"
	"strconv"
	"strings"
)

// RormEngine - Raw Query ORM Engine structure
type RormEngine struct {
	DB             *sql.DB
	dbDriver       string
	result         map[string]string
	results        []map[string]string
	conditionValue []interface{}
	rawQuery       string
	column         string
	orderBy        string
	condition      string
	tableName      string
	tablePrefix    string
	limit          string
	join           string
	groupBy        string
}

const (
	MYSQL_PREPARED_PARAM    = "?"
	POSTGRES_PREPARED_PARAM = "$"
	ORACLE_PREPARED_PARAM   = ":"
)

func (re *RormEngine) SetDB(db *sql.DB) {
	re.DB = db
}

func (re *RormEngine) GetDB() *sql.DB {
	return re.DB
}

func (re *RormEngine) GetPreparedValues() []interface{} {
	return re.conditionValue
}

// InitRormEngine - init new RORM Engine
func New(dbDriver, connectionURL string, tbPrefix ...string) *RormEngine {
	re := &RormEngine{}
	re.Connect(dbDriver, connectionURL, tbPrefix...)
	return re
}

// Connect - connect to db Driver
func (re *RormEngine) Connect(dbDriver, connectionURL string, tbPrefix ...string) error {
	var err error
	re.DB, err = sql.Open(dbDriver, connectionURL)
	re.dbDriver = dbDriver
	if tbPrefix != nil {
		re.tablePrefix = tbPrefix[0]
	}
	return err
}

func (re *RormEngine) clearField() {
	re.condition = ""
	re.column = ""
	re.orderBy = ""
	re.tableName = ""
	re.limit = ""
	re.join = ""
}

func (re *RormEngine) adjustPreparedParam(old string) string {
	if strings.TrimSpace(re.dbDriver) != "mysql" {
		idx := 0
		replacement := ORACLE_PREPARED_PARAM
		if re.dbDriver == "postgres" {
			replacement = POSTGRES_PREPARED_PARAM
		}
		for {
			if strings.Index(old, "?") == -1 {
				break
			}
			idx++
			old = strings.Replace(old, "?", replacement+strconv.Itoa(idx), 1)
		}
	}
	return old
}

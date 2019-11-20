package rorm

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type (
	// Engine - Raw Query ORM Engine structure
	Engine struct {
		db               *sqlx.DB
		config           *DbConfig
		options          *DbOptions
		connectionString string
		operations
		builderOperations
	}

	// DbConfig - DB Connection struct
	DbConfig struct {
		Host     string
		Port     string
		Username string
		Password string
		DbName   string
		DbScheme string
		Driver   string
		Protocol string
		// DbInstance - for SQL Server
		DbInstance string
	}
	// DbOptions - options for DB structure
	DbOptions struct {
		tbPrefix  string
		tbPostfix string
		tbFormat  string //table Format (camel case/snake case)
		colFormat string //column format (camel case/snake case)
	}
	// Operations - list of property query string

	builderOperations struct {
		rawQueryBuilder  strings.Builder
		conditionBuilder strings.Builder
		columnBuilder    strings.Builder
		orderByBuilder   strings.Builder
		limitBuilder     strings.Builder
		joinBuilder      strings.Builder
		groupByBuilder   strings.Builder
	}
	operations struct {
		isRaw              bool
		syntaxQuote        string
		stmt               *sqlx.Stmt
		isBulk             bool
		isMultiRows        bool
		bulkOptimized      bool
		bulkCounter        int
		updatedCol         map[string]bool
		tmpStruct          interface{}
		counter            int
		rawQuery           string
		column             string
		orderBy            string
		condition          string
		preparedValue      []interface{}
		multiPreparedValue [][]interface{}
		tableName          string
		limit              string
		join               string
		groupBy            string
		having             string
	}
)

const (
	MYSQL_PREPARED_PARAM    = "?"
	POSTGRES_PREPARED_PARAM = "$"
	ORACLE_PREPARED_PARAM   = ":"
	MSSQL_PREPARED_PARAM    = "@"
)

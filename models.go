package rorm

import "github.com/jmoiron/sqlx"

type (
	// Engine - Raw Query ORM Engine structure
	Engine struct {
		db               *sqlx.DB
		config           *DbConfig
		options          *DbOptions
		result           map[string]string
		results          []map[string]string
		connectionString string
		operations
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
	operations struct {
		isRaw              bool
		isBulk             bool
		isMultiRows        bool
		bulkOptimized      bool
		bulkCounter        int
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

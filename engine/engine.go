package rorm

import "database/sql"

// RormEngine - Raw Query ORM Engine structure
type RormEngine struct {
	DB             *sql.DB
	result         map[string]string
	results        []map[string]string
	conditionValue []string
	rawQuery       string
	column         string
	orderBy        string
	condition      string
	tableName      string
	limit          string
}

// InitRormEngine - init new RORM Engine
func InitRormEngine() *RormEngine {
	return &RormEngine{}
}

// Connect - connect to db Driver
func (re *RormEngine) Connect(dbDriver, connectionURL string) error {
	var err error
	re.DB, err = sql.Open(dbDriver, connectionURL)
	return err
}

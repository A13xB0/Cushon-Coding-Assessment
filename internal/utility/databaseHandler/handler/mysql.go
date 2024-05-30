package databaseHandlers

import (
	"database/sql"
	"fmt"
)

type MySql struct {
	databaseConnection *sql.DB // The database connection
	MySqlConfig                // Composition of config
}

type MySqlConfig struct {
	Hostname string
	Username string
	Password string
	Database string
}

// Connect to the database
func (db *MySql) Connect(config interface{}) error {
	db.MySqlConfig = config.(MySqlConfig)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", db.Username, db.Password, db.Hostname, db.Database)
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	db.databaseConnection = conn
	return nil
}

// Disconnect from the database
func (db *MySql) Disconnect() {
	if db.databaseConnection != nil {
		db.databaseConnection.Close()
	}
}

// Reset connection from the database
func (db *MySql) ResetConnection() error {
	db.Disconnect()
	return db.Connect(db.MySqlConfig)
}

// Check connection to the database
func (db *MySql) CheckConnection() bool {
	err := db.databaseConnection.Ping()
	return err == nil
}

// Query the database and unmarshal results into format specified in the callback
func (db *MySql) Query(statement string, parameters interface{}, unmarhsalItemCallback func(interface{}) (items interface{}, err error)) (interface{}, error) {
	params := parameters.([]interface{})
	rows, err := db.databaseConnection.Query(statement, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	output, err := unmarhsalItemCallback(rows)
	return output, err
}

// Execute statement
func (db *MySql) Execute(statement string, parameters interface{}) error {
	params := parameters.([]interface{})
	if _, err := db.databaseConnection.Exec(statement, params...); err != nil {
		return err
	}
	return nil
}

// This is a database handler utility, I have used an interface method in the event the database type changes
// i.e. it is currently mysql, but if it because an oracle database, dynamo database, or a form of immutable ledger
// this is designed with the assumption that most databases will use an SQL statement and allows for various config
// structs for various configuration changes.
package databaseHandler

import (
	databaseHandlers "cushioninterview/internal/utility/databaseHandler/handler"
	"fmt"
)

type DatabaseType int

const (
	SQL DatabaseType = 0
)

type DatabaseHandler interface {

	// Connect to the database
	Connect(config any) error

	// Disconnect from the database
	Disconnect()

	// Reset connection from the database
	ResetConnection() error

	// Check connection to the database
	CheckConnection() bool

	// Query the database and unmarshal results into format specified in the callback
	Query(statement string, parameters interface{}, unmarhsalItemCallback func(interface{}) (items interface{}, err error)) (interface{}, error)

	// Execute statement
	Execute(statement string, parameters interface{}) error
}

func New(databaseType DatabaseType) (DBHandler DatabaseHandler, err error) {
	switch databaseType {
	case SQL:
		return &databaseHandlers.MySql{}, nil
	default:
		return nil, fmt.Errorf("")
	}
}

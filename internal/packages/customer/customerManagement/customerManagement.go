// This package allows for the management of the customer accounts, this presumes we are using SQL statements
package customerManagement

import (
	"cushioninterview/internal/utility/authenticate"
	"cushioninterview/internal/utility/databaseHandler"
	"database/sql"
	"fmt"
)

// SQL statements
const updateEmailStatement = "UPDATE Customers SET customer_email = ? WHERE customer_id = ?;"
const updatePasswordStatement = "UPDATE Customers SET customer_password = ? WHERE customer_id = ?;"
const selectPasswordStatement = "SELECT customer_password FROM Customers WHERE customer_id = ?;"

// This function updates the customers email in the database
func UpdateEmail(customerId int, email string, dbHandler databaseHandler.DatabaseHandler) error {
	parameters := []interface{}{email, customerId}
	dbHandler.Execute(updateEmailStatement, parameters)
	return nil
}

// This function updates the customers password in the database
func UpdatePassword(customerId int, oldPassword, newPassword string, dbHandler databaseHandler.DatabaseHandler) error {
	oldPasswordHashed := authenticate.HashPassword(oldPassword)
	//Check old password
	parameters := []interface{}{customerId}
	dbPassword, err := dbHandler.Query(selectPasswordStatement, parameters, unmarshalPassword)
	if err != nil {
		return err
	}
	if dbPassword != oldPasswordHashed {
		return fmt.Errorf("passwords do not match")
	}
	//Update password
	parameters = []interface{}{customerId, newPassword}
	dbHandler.Execute(updatePasswordStatement, parameters)
	return nil
}

// This allows for the unmarshalling of the password from the SQL rows
func unmarshalPassword(data interface{}) (interface{}, error) {
	dbrows := data.(*sql.Rows)
	var password string
	var count int
	for dbrows.Next() {
		if err := dbrows.Scan(&password); err != nil {
			return nil, err
		}
		count++
	}
	//Returned more than one result
	if count > 1 {
		return nil, fmt.Errorf("more than one result should not be returned")
	}
	if count < 1 {
		return nil, fmt.Errorf("account does not exist")
	}

	return password, nil
}

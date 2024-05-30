// This allows to manage and submit investments.
package investmentManagement

import (
	"cushioninterview/internal/models"
	"cushioninterview/internal/utility/databaseHandler"
	"database/sql"
	"fmt"
	"time"
)

// SQL Statements
const getCustomerInvestmentsStatement = "SELECT I.investment_id, F.fund_name, I.fund_id, I.amount_invested, I.date_invested FROM Investments I JOIN Funds F ON I.fund_id = F.fund_id WHERE I.customer_id = ?;"
const getFundsStatement = "Select * FROM Funds"
const submitInvestmentStatement = "INSERT INTO Investments (customer_id, fund_id, amount_invested, date_invested) VALUES (?, ?, ?, ?);"

// This function will return all the different funds available for the customer to be able to invest in.
func GetFunds(dbHandler databaseHandler.DatabaseHandler) (funds []string, err error) {
	rows, err := dbHandler.Query(getFundsStatement, nil, getFundsUnmarshal)
	if err != nil {
		return []string{}, err
	}
	return rows.([]string), nil
}

// This function unmarshals the SQL rows into a format
func getFundsUnmarshal(data interface{}) (interface{}, error) {
	dbrows := data.(*sql.Rows)
	var fundName string
	var fundSlice []string
	for dbrows.Next() {
		if err := dbrows.Scan(&fundName); err != nil {
			return nil, err
		}
		fundSlice = append(fundSlice, fundName)
	}
	return fundSlice, nil
}

// This function gets the customers current investments
func GetCustomerInvestments(customerId int, dbHandler databaseHandler.DatabaseHandler) (Investments []models.InvestmentRow, err error) {
	parameters := []interface{}{customerId}
	rows, err := dbHandler.Query(getCustomerInvestmentsStatement, parameters, getCustomerInvestmentsUnmarshal)
	if err != nil {
		return []models.InvestmentRow{}, nil
	}
	InvestmentRows := rows.([]models.InvestmentRow)
	return InvestmentRows, nil
}

// This function unmarshals the SQL rows into a format
func getCustomerInvestmentsUnmarshal(data interface{}) (interface{}, error) {
	dbrows := data.(*sql.Rows)
	var investmentId int
	var fundName string
	var fundId int
	var amountInvested int
	var dateInvested int64
	var investmentSlice []models.InvestmentRow
	for dbrows.Next() {
		if err := dbrows.Scan(&investmentId, &fundName, &fundId, &amountInvested, &dateInvested); err != nil {
			return nil, err
		}
		investmentSlice = append(investmentSlice,
			models.InvestmentRow{
				InvestmentId:   investmentId,
				FundName:       fundName,
				FundId:         fundId,
				AmountInvested: float64(amountInvested),
				DateInvested:   time.Unix(dateInvested, 0),
			},
		)
	}
	return investmentSlice, nil
}

// This function submits investments for the customer to the SQL database.
func SubmitInvestment(customerId, fundId int, amountInvested float64, dbHandler databaseHandler.DatabaseHandler) (err error) {
	parameters := []interface{}{customerId, fundId, amountInvested, time.Now().Unix()}
	if err := dbHandler.Execute(submitInvestmentStatement, parameters); err != nil {
		return fmt.Errorf("could not execute investment for %v, %v", customerId, err)
	}
	return nil
}

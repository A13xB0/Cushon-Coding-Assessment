// These are models for our HTTP post requests of the investment REST APIs
package models

import "time"

type InvestmentHTTPGetFunds struct {
	AccessToken string `json:"access_token"`
}

type InvestmentHTTPGetCustInvestments struct {
	AccessToken string `json:"access_token"`
}

type InvestmentHTTPSubmitInvestment struct {
	FundId         int     `json:"fund_id"`
	AmountInvested float64 `json:"amount_invested"`
	AccessToken    string  `json:"access_token"`
}

type InvestmentHTTPReturnGetFunds struct {
	Funds []string
}

type InvestmentHTTPReturnGetCustInvestments struct {
	Investments []InvestmentRow
}

type InvestmentRow struct {
	InvestmentId   int
	FundName       string
	FundId         int
	AmountInvested float64
	DateInvested   time.Time
}

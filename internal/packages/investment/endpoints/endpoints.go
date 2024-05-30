package endpoints

import (
	"cushioninterview/internal/models"
	"cushioninterview/internal/packages/investment/investmentManagement"
	"cushioninterview/internal/utility/authenticate"
	"cushioninterview/internal/utility/databaseHandler"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServiceEndpoints struct {
	AuthHandler authenticate.Authenticate
	Log         *zap.SugaredLogger
	DbHandler   databaseHandler.DatabaseHandler
}

func (se *ServiceEndpoints) GetReady(c *gin.Context) {
	// Once a backend database is created this will check connection
	// is still active returning a negative response if not
	if se.DbHandler.CheckConnection() {
		c.Writer.WriteHeader(204)
		return
	}
	c.Writer.WriteHeader(503) //Internal Servier Error
}

// This endpoint shows the service is read returning no data
func (se *ServiceEndpoints) GetLive(c *gin.Context) {
	c.Writer.WriteHeader(204)
}

func (se *ServiceEndpoints) GetFunds(c *gin.Context) {
	var requestBody models.InvestmentHTTPGetFunds
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad format"}) // obscure error as real error can be logged out
		return
	}
	//Authentication
	_, err := se.AuthHandler.ValidateUser(requestBody.AccessToken)
	if err != nil {
		c.Writer.WriteHeader(401)
		se.Log.Errorf("Could not authenticate user", zap.Error(err))
		return
	}

	fundRows, err := investmentManagement.GetFunds(se.DbHandler)
	if err != nil {
		c.Writer.WriteHeader(500)
		se.Log.Errorf("%v", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.InvestmentHTTPReturnGetFunds{
		Funds: fundRows,
	})
}

func (se *ServiceEndpoints) GetCustomerInvestments(c *gin.Context) {
	var requestBody models.InvestmentHTTPGetCustInvestments
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad format"}) // obscure error as real error can be logged out
		return
	}
	//Authentication
	userPayload, err := se.AuthHandler.ValidateUser(requestBody.AccessToken)
	if err != nil {
		c.Writer.WriteHeader(401)
		se.Log.Errorf("Could not authenticate user", zap.Error(err))
		return
	}

	investmentRows, err := investmentManagement.GetCustomerInvestments(userPayload.CustomerId, se.DbHandler)
	if err != nil {
		c.Writer.WriteHeader(500)
		se.Log.Errorf("%v", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.InvestmentHTTPReturnGetCustInvestments{
		Investments: investmentRows,
	})
}

func (se *ServiceEndpoints) SubmitInvestment(c *gin.Context) {
	var requestBody models.InvestmentHTTPSubmitInvestment
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad format"}) // obscure error as real error can be logged out
		return
	}
	//Authentication
	userPayload, err := se.AuthHandler.ValidateUser(requestBody.AccessToken)
	if err != nil {
		c.Writer.WriteHeader(401)
		se.Log.Errorf("Could not authenticate user", zap.Error(err))
		return
	}

	if err := investmentManagement.SubmitInvestment(userPayload.CustomerId, requestBody.FundId, requestBody.AmountInvested, se.DbHandler); err != nil {
		c.Writer.WriteHeader(500)
		se.Log.Errorf("%v", zap.Error(err))
		return
	}
	c.Writer.WriteHeader(200)
	return
}

package endpoints

import (
	"cushioninterview/internal/models"
	"cushioninterview/internal/packages/customer/customerManagement"
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

// This endpoint updates the customers email address
func (se *ServiceEndpoints) UpdateEmail(c *gin.Context) {
	var requestBody models.CustHTTPUpdateEmail
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
	//Update Email
	if err := customerManagement.UpdateEmail(userPayload.CustomerId, requestBody.Email, se.DbHandler); err != nil {

	}
	c.Writer.WriteHeader(http.StatusOK)
}

// This endpoint updates the customers email password
func (se *ServiceEndpoints) UpdatePassword(c *gin.Context) {
	var requestBody models.CustHTTPUpdatePassword
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
	customerManagement.UpdatePassword(userPayload.CustomerId, requestBody.OldPassword, requestBody.NewPassword, se.DbHandler)
	c.Writer.WriteHeader(http.StatusOK)
}

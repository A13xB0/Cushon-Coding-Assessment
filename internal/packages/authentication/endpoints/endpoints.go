package endpoints

import (
	"cushioninterview/internal/models"
	"cushioninterview/internal/utility/authenticate"
	"cushioninterview/internal/utility/databaseHandler"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Custom routes struct obj
type ServiceEndpoints struct {
	AuthHandler authenticate.Authenticate
	DbHandler   databaseHandler.DatabaseHandler
	Log         *zap.SugaredLogger
}

// This endpoint shows the service is read returning no data
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

// This endpoint authenticates the user using JWT
func (se *ServiceEndpoints) Authenticate(c *gin.Context) {
	var requestBody models.AuthHTTPPostPayload
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad format"}) // obscure error as real error can be logged out
		se.Log.Errorf("bad json format, %v", requestBody.Username, zap.Error(err))
		return
	}
	bearerToken, err := se.AuthHandler.AuthenticateUser(requestBody)
	if err != nil {
		c.Writer.WriteHeader(401) // Unauthorised is used due to failed login attempt so the server does not know who the client is
		se.Log.Errorf("could not log in user %v, %v", requestBody.Username, zap.Error(err))
		return
	}
	// Return json from struct
	returnPayload := models.AuthHTTPReturnPayload{
		AccessToken: bearerToken,
		TokenType:   "Bearer",
	}
	c.JSON(http.StatusOK, returnPayload)
}

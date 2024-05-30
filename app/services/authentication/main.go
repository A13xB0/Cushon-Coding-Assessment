// This service allows for stateless authentication for the client between services via REST API. This service ultilises SQL and JWT. The client will first request authentication
// with the SQL username and passowrd, a JWT token with the account inforamtion will be created and this will be stored in the browser for future API calls to other services.
// Other services can validate this token using the authenticate utility package to check the authentication of the bearer token and account information.
package main

import (
	AuthConfiguration "cushioninterview/internal/packages/authentication/configuration"
	"cushioninterview/internal/packages/authentication/endpoints"
	"cushioninterview/internal/utility/authenticate"
	"cushioninterview/internal/utility/databaseHandler"
	databaseHandlers "cushioninterview/internal/utility/databaseHandler/handler"
	"cushioninterview/internal/utility/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type authService struct {
	AuthConfiguration.AuthConfig
	log         *zap.SugaredLogger
	stopChannel chan os.Signal
	router      *gin.Engine
	server      *http.Server
}

// This function sets up the service ready to start
func createService(config AuthConfiguration.AuthConfig, log *zap.SugaredLogger) (service authService, err error) {
	// Database setup and connection
	dbHandler, err := databaseHandler.New(databaseHandler.SQL) //Create database handler
	if err != nil {
		return authService{}, fmt.Errorf("could not create database handler, %v", err)
	}
	dbConfig := databaseHandlers.MySqlConfig{
		Hostname: config.SqlHostname,
		Username: config.SqlUsername,
		Password: config.SqlPassword,
		Database: config.SqlDatabase,
	}
	dbHandler.Connect(dbConfig)
	// Authentication Package
	authHandler := authenticate.NewAuthenticator(config.MAC, dbHandler, 15*time.Minute)
	// Endpoints
	ep := endpoints.ServiceEndpoints{
		Log:         log,
		AuthHandler: authHandler,
		DbHandler:   dbHandler,
	}
	// Gin Setup
	gin.SetMode(gin.ReleaseMode) // set to production
	router := gin.New()          // empty engine
	router.Use(gin.Recovery())   // allow app to recover from sigv
	//Setup Routes
	router.GET("/ready", ep.GetReady)
	router.GET("/live", ep.GetLive)
	//API V1 routes
	v1 := router.Group("/api/v1") //API Version Endpoint Group
	v1.POST("/authenticate", ep.Authenticate)

	return authService{
		AuthConfig:  config,
		log:         log,
		stopChannel: make(chan os.Signal),
		router:      router,
	}, nil
}

// This function starts the service
func (service *authService) startService() {
	// start service
	service.service()

	// wait for stop signal
	<-service.stopChannel

	// gracefully shutdown service
	service.stopService()
}

// This function stops the service gracefully
func (service *authService) stopService() {
	if err := service.server.Close(); err != nil {
		//service.logger.Error("could not close http server", err)
		service.log.Error("could not close http server", zap.Error(err))
	}
}

// This function is the service and its processing
func (service *authService) service() {
	// Set up signal for receiving shutdown
	signal.Notify(service.stopChannel, syscall.SIGTERM, syscall.SIGINT)
	//Make a http server
	service.server = &http.Server{
		Addr:              fmt.Sprintf("%v:%v", service.WebServiceHost, service.WebServicePort),
		Handler:           service.router,
		ReadHeaderTimeout: 30 * time.Second,
	}
	go func() {
		if err := service.server.ListenAndServe(); err != nil {
			service.log.Error("http server error", zap.Error(err))
			service.stopChannel <- nil
		} else {
			service.log.Infof("Authentication Service listening on %v:%v")
		}
	}()
}

func main() {
	config := AuthConfiguration.GetConfig() //Get configuration from flags
	log, err := logger.New(config.LoggerServiceName, config.LoggerLevel)
	if err != nil {
		fmt.Println("could not start logger, fatal")
		os.Exit(1)
	}
	service, err := createService(config, log)
	if err != nil {
		log.Fatal("could not create service", zap.Error(err))
	}
	service.startService()
}

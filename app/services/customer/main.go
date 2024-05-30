// This service is to allow customer management via a REST API, in this service the intention is to have controls which allow for staff to create accounts (if roles were implemented,
// which they aren't) and for the customer to be able to update and manage their account and settings on their account.
package main

import (
	CustomerConfiguration "cushioninterview/internal/packages/customer/configuration"
	"cushioninterview/internal/packages/customer/endpoints"
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

type customerService struct {
	CustomerConfiguration.CustomerConfig
	log         *zap.SugaredLogger
	stopChannel chan os.Signal
	router      *gin.Engine
	server      *http.Server
}

// This function sets up the service ready to start
func createService(config CustomerConfiguration.CustomerConfig, log *zap.SugaredLogger) (service customerService, err error) {
	// Database setup and connection
	dbHandler, err := databaseHandler.New(databaseHandler.SQL) //Create database handler
	if err != nil {
		return customerService{}, fmt.Errorf("could not create database handler, %v", err)
	}
	dbConfig := databaseHandlers.MySqlConfig{
		Hostname: config.SqlHostname,
		Username: config.SqlUsername,
		Password: config.SqlPassword,
		Database: config.SqlDatabase,
	}
	dbHandler.Connect(dbConfig)
	// Authentication Package
	authHandler := authenticate.NewValidator(config.MAC)
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
	v1.POST("/updateEmail", ep.UpdateEmail)
	v1.POST("/updatePassword", ep.UpdatePassword)

	return customerService{
		CustomerConfig: config,
		log:            log,
		stopChannel:    make(chan os.Signal),
		router:         router,
	}, nil
}

// This function starts the service
func (service *customerService) startService() {
	// start service
	service.service()

	// wait for stop signal
	<-service.stopChannel

	// gracefully shutdown service
	service.stopService()
}

// This function stops the service gracefully
func (service *customerService) stopService() {
	if err := service.server.Close(); err != nil {
		//service.logger.Error("could not close http server", err)
		service.log.Error("could not close http server", zap.Error(err))
	}
}

// This function is the service and its processing
func (service *customerService) service() {
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

// The service starts here
func main() {
	config := CustomerConfiguration.GetConfig() //Get configuration from flags
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

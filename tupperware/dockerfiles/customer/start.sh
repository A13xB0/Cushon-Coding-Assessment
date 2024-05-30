#!/bin/sh

# Set default values for environment variables
: ${SQL_HOSTNAME:=127.0.0.1}
: ${SQL_PORT:=1443}
: ${SQL_USERNAME:=username}
: ${SQL_PASSWORD:=password}
: ${SQL_DATABASE:=CushionDB}
: ${WEB_SERVER_HOST:=0.0.0.0}
: ${WEB_SERVER_PORT:=8080}
: ${MAC:=K1ng5l3y}
: ${LOGGER_SERVICE_NAME:=authenticator-1}
: ${LOGGER_LEVEL:=info}

# Run the Go application with the specified flags
./authentication \
  -sqlhostname "$SQL_HOSTNAME" \
  -sqlport "$SQL_PORT" \
  -sqlusername "$SQL_USERNAME" \
  -sqlpassword "$SQL_PASSWORD" \
  -sqldatabase "$SQL_DATABASE" \
  -webserverhost "$WEB_SERVER_HOST" \
  -webserverport "$WEB_SERVER_PORT" \
  -mac "$MAC" \
  -servicename "$LOGGER_SERVICE_NAME" \
  -loggerlevel "$LOGGER_LEVEL"
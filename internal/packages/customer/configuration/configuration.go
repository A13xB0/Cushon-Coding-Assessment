// This is the configuration package for the customer service. This is managed by flags but if I had more time and this was a production project I would use VIPER so we can take in
// env variables instead of having to convert from env variable to flags via a start.sh in the container.
package CustomerConfiguration

import "flag"

//Configuration struct
type CustomerConfig struct {
	//SQL Config
	SqlHostname string //SQL Hostname e.g. 127.0.0.1
	SqlPort     int    //SQL Port e.g. 1443
	SqlUsername string //SQL Username e.g. username
	SqlPassword string //SQL Password e.g. pa$$w0rd
	SqlDatabase string //SQL Database e.g. CushionDB
	//Host Config
	WebServiceHost string //WebService Host e.g. 0.0.0.0
	WebServicePort int    //WebService Port e.g. 8443
	//Service Config
	MAC string //Message Authentication Code e.g. K1ng5l3y
	//Logger Config
	LoggerServiceName string //LoggerServiceName e.g. Auth-Service-1
	LoggerLevel       string //LoggerLevel e.g. info, debug
}

// This function parses the flags and places them into a usable struct for the service
func GetConfig() CustomerConfig {
	//Set Flag Pointers
	sqlHostname := flag.String("sqlhostname", "127.0.0.1", "SQL Hostname")
	sqlPort := flag.Int("sqlport", 1443, "SQL Port")
	sqlUsername := flag.String("sqlusername", "username", "SQL Username")
	sqlPassword := flag.String("sqlpassword", "password", "SQL Password")
	sqlDatabase := flag.String("sqldatabase", "CushionDB", "SQL Database")
	webServerHost := flag.String("webserverhost", "0.0.0.0", "Web Hostname")
	webServicePort := flag.Int("webserverhost", 8081, "Web Hostname")
	mac := flag.String("mac", "K1ng5l3y", "Message Authentication Code, used for authentication")
	loggerServiceName := flag.String("servicename", "customer-1", "The service name used for the logger")
	loggerLevel := flag.String("loggerlevel", "info", "Logger Level info|debug")
	//Parse Flags
	flag.Parse()
	//Place into return struct
	return CustomerConfig{
		SqlHostname:       *sqlHostname,
		SqlPort:           *sqlPort,
		SqlUsername:       *sqlUsername,
		SqlPassword:       *sqlPassword,
		SqlDatabase:       *sqlDatabase,
		MAC:               *mac,
		WebServiceHost:    *webServerHost,
		WebServicePort:    *webServicePort,
		LoggerServiceName: *loggerServiceName,
		LoggerLevel:       *loggerLevel,
	}
}

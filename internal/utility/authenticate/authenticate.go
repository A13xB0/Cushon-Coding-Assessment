// This utility allows services to authenticate and validate JWT bearer tokens for authentication into our services
package authenticate

import (
	"cushioninterview/internal/models"
	"cushioninterview/internal/utility/databaseHandler"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/blake2b"
)

// This is the struct to set variables in the authenticate utility
type Authenticate struct {
	hmac      string                          //Required for both authenticator and validator
	timeout   time.Duration                   //Required for authenticator
	dbHandler databaseHandler.DatabaseHandler //Required for authenticator
}

// SQL statements
const GetCustomerAccount = "SELECT customer_id, customer_name, customer_email, customer_type FROM Customer WHERE customer_email = ? AND customer_password = ?;"

// This function creates a new authenticator which has the ability to sign new bearer tokens
func NewAuthenticator(hmac string, dbHandler databaseHandler.DatabaseHandler, timeout time.Duration) Authenticate {
	return Authenticate{
		hmac:      hmac,
		timeout:   timeout,
		dbHandler: dbHandler,
	}
}

// This function creates a new validator which only allows to validate existing bearer tokens
func NewValidator(hmac string) Authenticate {
	return Authenticate{
		hmac: hmac,
	}
}

// This function will authenticate a new user using a HTTP Post payload from models
func (Auth *Authenticate) AuthenticateUser(AuthPayload models.AuthHTTPPostPayload) (newBearerToken string, err error) {
	if Auth.timeout == 0*time.Second || Auth.dbHandler == nil {
		return "", fmt.Errorf("no timeout or database handler provided, this may have been initialised as a validator")
	}
	// Get database customer info from database and validate exists
	HashedPassword := HashPassword(AuthPayload.Password)
	parameters := []interface{}{AuthPayload.Username, HashedPassword}
	accountRow, err := Auth.dbHandler.Query(GetCustomerAccount, parameters, unmarshalCustomerAccount)
	if err != nil {
		return "", err
	}
	jwtPayload := accountRow.(models.AuthJwtPayload)
	// Create bearer token
	newBearerToken, err = Auth.jwtSignPayload(jwtPayload)
	if err != nil {
		return "", err
	}
	return newBearerToken, err
}

func unmarshalCustomerAccount(data interface{}) (interface{}, error) {
	dbrows := data.(*sql.Rows)
	newJwtPayload := models.AuthJwtPayload{}
	var count int
	for dbrows.Next() {
		if err := dbrows.Scan(&newJwtPayload.CustomerId, &newJwtPayload.Name, &newJwtPayload.Email); err != nil {
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

	return newJwtPayload, nil
}

// This functions signs the JWT payload
func (Auth *Authenticate) jwtSignPayload(payload models.AuthJwtPayload) (bearerToken string, err error) {
	payload.Expiry = time.Now().Add(Auth.timeout).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customerId": payload.CustomerId,
		"email":      payload.Email,
		"name":       payload.Name,
		"expiry":     payload.Expiry,
		"type":       payload.Type,
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(Auth.hmac))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// This function validates a bearer token and returns its payload
func (Auth *Authenticate) ValidateUser(tokenString string) (payload models.AuthJwtPayload, err error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Auth.hmac), nil
	})

	if err != nil {
		return models.AuthJwtPayload{}, err
	}

	// Verify the token claims and extract the payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return models.AuthJwtPayload{}, fmt.Errorf("invalid token")
	}

	// Create a new AuthJwtPayload struct and populate it with the extracted values
	payload = models.AuthJwtPayload{
		CustomerId: claims["customerId"].(int),
		Email:      claims["email"].(string),
		Name:       claims["name"].(string),
		Expiry:     claims["expiry"].(int64),
		Type:       claims["type"].(string),
	}

	//Check if payload has expired
	if payload.Expiry > time.Now().Unix() {
		return models.AuthJwtPayload{}, fmt.Errorf("token has expired")
	}

	return payload, nil
}

// This hashes a passowrd using blake2b
// For demonstration only, this does not protect against certain attacks
func HashPassword(data string) string {
	hash := blake2b.Sum512([]byte(data))              //Create hash object and ingest data
	return base64.StdEncoding.EncodeToString(hash[:]) //return base64 string of sum of hash object
}

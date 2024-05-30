// These are models for our jwt authentication to be used cross services
package models

// This is the authentication payload for within the jwt token
type AuthJwtPayload struct {
	CustomerId int    `json:"customerId"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Expiry     int64  `json:"expiry"`
	Type       string `json:"type"`
}

// This is the HTTP Post payload to send to the authentication server
type AuthHTTPPostPayload struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

// Authentication HTTP Return payload as defined by OAuth2.0
type AuthHTTPReturnPayload struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

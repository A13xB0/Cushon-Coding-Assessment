// These are models for our HTTP post requests of the customer REST APIs
package models

// This is a model for updating an email in the customer account
type CustHTTPUpdateEmail struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

// This is a model for updating the password in the customer account
type CustHTTPUpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	AccessToken string `json:"access_token"`
}

package domain

// Client stores oauth client info
type Client struct {
	ID       int64  `json:"id"`        // ID is the internal id of the client
	Name     string `json:"name"`      // Name is the client display name
	PublicId string `json:"public_id"` // PublicId is the client public id
	Secret   string `json:"secret"`    // Secret is required for generate an access token to create accounts
}

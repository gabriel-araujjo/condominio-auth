package domain

import (
	"net/url"
)

// Email store the email address and a flag indicating whether the email was verified
type Email struct {
	Email    string `json:"email"`    // The email
	Verified bool   `json:"verified"` // Whether the email was verified
}

// Phone store the phone number and a flag indicating whether the phone was verified
type Phone struct {
	Phone    string `json:"phone"`    // The phone
	Verified bool   `json:"verified"` // Whether the phone was verified
}

// User stores basic user info
type User struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	CPF          string   `json:"cpf"`
	FbID         string   `json:"fb_id"`
	Avatar       *url.URL `json:"avatar"`
	Phones       []Phone  `json:"phones"`
	Emails       []Email  `json:"emails"`
	PasswordHash string   `json:"-"`
}

// PrimaryEmail returns the user's primary email
func (u *User) PrimaryEmail() *Email {
	for _, email := range u.Emails {
		return &email
	}
	return nil
}

// PrimaryPhone returns the user's primary phone
func (u *User) PrimaryPhone() *Phone {
	for _, phone := range u.Phones {
		return &phone
	}
	return nil
}


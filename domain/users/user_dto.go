package users

import (
	"github.com/Abacode7/bookstore_utils-go/v2/rest_error"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type Users []User

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

func (user *User) Validate() rest_error.RestErr {
	if user.Email == "" {
		return rest_error.NewBadRequestError("invalid email address")
	}
	if user.Password == "" {
		return rest_error.NewBadRequestError("invalid password")
	}
	return nil
}

package users

import "github.com/Abacode7/bookstore_users-api/utils/errors"

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

func (user *User) Validate() *errors.RestErr {
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	if user.Password == "" {
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}

package users

import (
	"github.com/Abacode7/bookstore_utils-go/v2/rest_error"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:password`
}

func (ulr *UserLoginRequest) Validate() rest_error.RestErr {
	if ulr.Email == "" {
		return rest_error.NewBadRequestError("invalid email address")
	}
	if ulr.Password == "" {
		return rest_error.NewBadRequestError("invalid password")
	}
	return nil
}

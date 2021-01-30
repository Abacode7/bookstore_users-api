package users

import "github.com/Abacode7/bookstore_users-api/utils/errors"

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:password`
}

func (ulr *UserLoginRequest) Validate() *errors.RestErr {
	if ulr.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	if ulr.Password == "" {
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}

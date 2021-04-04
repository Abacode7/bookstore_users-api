package users

import (
	"encoding/json"
	"github.com/Abacode7/bookstore_utils-go/v2/logger"
	"github.com/Abacode7/bookstore_utils-go/v2/rest_error"
)

type PublicUser struct {
	//Id          int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	//Email       string `json:"email"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

func (user *User) Marshall(isPublic bool) (interface{}, rest_error.RestErr) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		logger.Error("error while marshalling user data", err)
		return nil, rest_error.NewInternalServerError("error while marshalling user data")
	}
	if isPublic {
		var publicUser PublicUser
		if err := json.Unmarshal(jsonData, &publicUser); err != nil {
			logger.Error("error while formatting user data", err)
			return nil, rest_error.NewInternalServerError("error while formatting user data")
		}
		return publicUser, nil
	}
	var privateUser PrivateUser
	if err := json.Unmarshal(jsonData, &privateUser); err != nil {
		logger.Error("error while formatting user data", err)
		return nil, rest_error.NewInternalServerError("error while formatting user data")
	}
	return privateUser, nil
}

func (users Users) Marshall(isPublic bool) (interface{}, rest_error.RestErr) {
	jsonData, err := json.Marshal(users)
	if err != nil {
		logger.Error("error while marshalling user data", err)
		return nil, rest_error.NewInternalServerError("error while marshalling user data")
	}
	if isPublic {
		var publicUsers []PublicUser
		if err := json.Unmarshal(jsonData, &publicUsers); err != nil {
			logger.Error("error while formatting user data", err)
			return nil, rest_error.NewInternalServerError("error while formatting user data")
		}
		return publicUsers, nil
	}
	var privateUsers []PrivateUser
	if err := json.Unmarshal(jsonData, &privateUsers); err != nil {
		logger.Error("error while formatting user data", err)
		return nil, rest_error.NewInternalServerError("error while formatting user data")
	}
	return privateUsers, nil
}

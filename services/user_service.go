package services

import (
	"github.com/Abacode7/bookstore_users-api/domain/users"
	"github.com/Abacode7/bookstore_users-api/utils/crypto_utils"
	"github.com/Abacode7/bookstore_users-api/utils/date_utils"
	"github.com/Abacode7/bookstore_utils-go/v2/logger"
	"github.com/Abacode7/bookstore_utils-go/v2/rest_error"
)

type IUserService interface {
	CreateUser(users.User) (*users.User, rest_error.RestErr)
	GetUser(int64) (*users.User, rest_error.RestErr)
	SearchUser(string) (users.Users, rest_error.RestErr)
	UpdateUser(bool, users.User) (*users.User, rest_error.RestErr)
	DeleteUser(int64) rest_error.RestErr
	LoginUser(users.UserLoginRequest) (*users.User, rest_error.RestErr)
}

type userService struct {
	userDao users.IUserDao
}

/// NewUserService is userService's constructor
func NewUserService(userDao users.IUserDao) IUserService {
	return &userService{userDao: userDao}
}

func (us *userService) CreateUser(user users.User) (*users.User, rest_error.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	var err error
	user.Password, err = crypto_utils.GetHash(user.Password)
	if err != nil {
		logger.Error("error generating password hash", err)
		restErr := rest_error.NewBadRequestError("invalid user password")
		return nil, restErr
	}
	user.DateCreated = date_utils.GetDbFormattedTime()
	user.Status = users.StatusActive

	newUser, daoErr := us.userDao.Save(user)
	if daoErr != nil {
		return nil, daoErr
	}
	return newUser, nil
}

func (us *userService) GetUser(userID int64) (*users.User, rest_error.RestErr) {
	return us.userDao.Get(userID)
}

func (us *userService) SearchUser(status string) (users.Users, rest_error.RestErr) {
	return us.userDao.FindByStatus(status)
}

func (us *userService) UpdateUser(isTotalUpdate bool, user users.User) (*users.User, rest_error.RestErr) {
	oldUser, getErr := us.userDao.Get(user.Id)
	if getErr != nil {
		return nil, getErr
	}
	// For fields email, password, status and date_created, if values
	// aren't provided, they retain their old values.
	if user.Password == "" {
		user.Password = oldUser.Password
	} else {
		var err error
		user.Password, err = crypto_utils.GetHash(user.Password)
		if err != nil {
			logger.Error("error generating password hash", err)
			restErr := rest_error.NewBadRequestError("invalid user password")
			return nil, restErr
		}
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}
	if user.Status == "" {
		user.Status = oldUser.Status
	}
	user.DateCreated = oldUser.DateCreated

	// For total update if values aren't provided for fields first_name
	// and last_name they take and empty default value
	if !isTotalUpdate {
		if user.FirstName == "" {
			user.FirstName = oldUser.FirstName
		}
		if user.LastName == "" {
			user.LastName = oldUser.LastName
		}
	}
	return us.userDao.Update(user)
}

func (us *userService) DeleteUser(userId int64) rest_error.RestErr {
	return us.userDao.Delete(userId)
}

func (us *userService) LoginUser(request users.UserLoginRequest) (*users.User, rest_error.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	user, err := us.userDao.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if err := crypto_utils.CompareHashAndPassword(user.Password, request.Password); err != nil {
		logger.Error("passwords do not match", err)
		return nil, rest_error.NewBadRequestError("wrong user password")
	}
	return user, nil
}

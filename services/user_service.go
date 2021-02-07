package services

import (
	"github.com/Abacode7/bookstore_users-api/domain/users"
	"github.com/Abacode7/bookstore_users-api/utils/crypto_utils"
	"github.com/Abacode7/bookstore_users-api/utils/date_utils"
	"github.com/Abacode7/bookstore_users-api/utils/errors"
	"github.com/Abacode7/bookstore_users-api/utils/logger"
)

type IUserService interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	SearchUser(string) (users.Users, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	LoginUser(users.UserLoginRequest) (*users.User, *errors.RestErr)
}

type userService struct {
	userDao users.IUserDao
}

/// NewUserService is userService's constructor
func NewUserService(userDao users.IUserDao) IUserService {
	return &userService{userDao: userDao}
}

func (us *userService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	var err error
	user.Password, err = crypto_utils.GetHash(user.Password)
	if err != nil {
		logger.Error("error generating password hash", err)
		restErr := errors.NewBadRequestError("invalid user password")
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

func (us *userService) GetUser(userID int64) (*users.User, *errors.RestErr) {
	return us.userDao.Get(userID)
}

func (us *userService) SearchUser(status string) (users.Users, *errors.RestErr) {
	return us.userDao.FindByStatus(status)
}

func (us *userService) UpdateUser(isTotalUpdate bool, user users.User) (*users.User, *errors.RestErr) {
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
			restErr := errors.NewBadRequestError("invalid user password")
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

func (us *userService) DeleteUser(userId int64) *errors.RestErr {
	return us.userDao.Delete(userId)
}

func (us *userService) LoginUser(request users.UserLoginRequest) (*users.User, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	user, err := us.userDao.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if err := crypto_utils.CompareHashAndPassword(user.Password, request.Password); err != nil {
		logger.Error("passwords do not match", err)
		return nil, errors.NewBadRequestError("wrong user password")
	}
	return user, nil
}

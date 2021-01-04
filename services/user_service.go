package services

import (
	"github.com/Abacode7/bookstore_users-api/domain/users"
	"github.com/Abacode7/bookstore_users-api/utils/date_utils"
	"github.com/Abacode7/bookstore_users-api/utils/errors"
)

type IUserService interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	SearchUser(string) ([]users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
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

	user.DateCreated = date_utils.GetDbFormattedTime()
	if user.Status == "" {
		user.Status = "inactive"
	}

	newUser, daoErr := us.userDao.Save(user)
	if daoErr != nil {
		return nil, daoErr
	}
	return newUser, nil
}

func (us *userService) GetUser(userID int64) (*users.User, *errors.RestErr) {
	user, err := us.userDao.Get(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) SearchUser(status string) ([]users.User, *errors.RestErr) {
	users, err := us.userDao.FindByStatus(status)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *userService) UpdateUser(isTotalUpdate bool, user users.User) (*users.User, *errors.RestErr) {

	oldUser, getErr := us.userDao.Get(user.Id)
	if getErr != nil {
		return nil, getErr
	}

	if !isTotalUpdate {
		if user.FirstName == "" {
			user.FirstName = oldUser.FirstName
		}
		if user.LastName == "" {
			user.LastName = oldUser.LastName
		}
		if user.Email == "" {
			user.Email = oldUser.Email
		}
		if user.Status == "" {
			user.Status = oldUser.Status
		}
		if user.Password == "" {
			user.Password = oldUser.Password
		}
	}
	user.DateCreated = oldUser.DateCreated

	result, updateErr := us.userDao.Update(user)
	if updateErr != nil {
		return nil, updateErr
	}
	return result, nil
}

func (us *userService) DeleteUser(userId int64) *errors.RestErr {
	if err := us.userDao.Delete(userId); err != nil {
		return err
	}
	return nil
}

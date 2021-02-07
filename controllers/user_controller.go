package controllers

import (
	"github.com/Abacode7/bookstore_users-api/domain/users"
	"github.com/Abacode7/bookstore_users-api/services"
	"github.com/Abacode7/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type IUserController interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	SearchUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

type userController struct {
	userService services.IUserService
}

/// NewUserController is userController's constructor
func NewUserController(us services.IUserService) *userController {
	return &userController{us}
}

func (uc *userController) CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	resultUser, serviceErr := uc.userService.CreateUser(user)
	if serviceErr != nil {
		c.JSON(serviceErr.Status, serviceErr)
		return
	}
	result, marshErr := resultUser.Marshall(true)
	if marshErr != nil {
		c.JSON(marshErr.Status, marshErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *userController) GetUser(c *gin.Context) {
	id := c.Param("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid parameter")
		c.JSON(restErr.Status, restErr)
		return
	}
	resultUser, serviceErr := uc.userService.GetUser(userID)
	if serviceErr != nil {
		c.JSON(serviceErr.Status, serviceErr)
		return
	}
	result, marshErr := resultUser.Marshall(true)
	if marshErr != nil {
		c.JSON(marshErr.Status, marshErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *userController) SearchUser(c *gin.Context) {
	value := strings.TrimSpace(c.Query("status"))
	users, err := uc.userService.SearchUser(value)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	result, marshErr := users.Marshall(c.GetHeader("x-public") == "true")
	if marshErr != nil {
		c.JSON(marshErr.Status, marshErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *userController) UpdateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		jsonErr := errors.NewBadRequestError("invalid json body")
		c.JSON(jsonErr.Status, jsonErr)
		return
	}
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		paramErr := errors.NewBadRequestError("invalid request parameter")
		c.JSON(paramErr.Status, paramErr)
		return
	}
	user.Id = userId

	var isTotalUpdate bool
	if c.Request.Method == "PUT" {
		isTotalUpdate = true
	} else {
		isTotalUpdate = false
	}
	resultUser, sevErr := uc.userService.UpdateUser(isTotalUpdate, user)
	if sevErr != nil {
		c.JSON(sevErr.Status, sevErr)
		return
	}
	result, marshErr := resultUser.Marshall(true)
	if marshErr != nil {
		c.JSON(marshErr.Status, marshErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (uc *userController) DeleteUser(c *gin.Context) {
	userId, strErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if strErr != nil {
		err := errors.NewBadRequestError("invalid request parameter")
		c.JSON(err.Status, err)
		return
	}
	if err := uc.userService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (uc *userController) LoginUser(c *gin.Context) {
	var ulr users.UserLoginRequest
	if err := c.ShouldBindJSON(&ulr); err != nil {
		restErr := errors.NewBadRequestError("invalid requests body")
		c.JSON(restErr.Status, restErr)
		return
	}
	resultUser, err := uc.userService.LoginUser(ulr)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	result, marshErr := resultUser.Marshall(c.GetHeader("x-public") == "true")
	if marshErr != nil {
		c.JSON(marshErr.Status, marshErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

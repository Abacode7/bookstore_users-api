package app

import (
	"github.com/Abacode7/bookstore_users-api/controllers"
)

func mapUrl(userCtlr controllers.IUserController) {
	router.GET("/ping", controllers.Ping)

	router.POST("/users", userCtlr.CreateUser)
	router.GET("/users/:user_id", userCtlr.GetUser)
	router.PUT("/users/:user_id", userCtlr.UpdateUser)
	router.PATCH("/users/:user_id", userCtlr.UpdateUser)
	router.DELETE("/users/:user_id", userCtlr.DeleteUser)
	router.POST("/users/login", userCtlr.LoginUser)

	router.GET("/internal/users/search", userCtlr.SearchUser)
}

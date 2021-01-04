package app

import (
	"fmt"
	"github.com/Abacode7/bookstore_users-api/controllers"
	"github.com/Abacode7/bookstore_users-api/datasources/mysql"
	"github.com/Abacode7/bookstore_users-api/domain/users"
	"github.com/Abacode7/bookstore_users-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var router = gin.Default()

func StartApplication() {

	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")

	/// Generates connections string and opens the connection
	/// using the provided data source
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := mysql.Init(dataSource)
	if err != nil {
		log.Fatal(err)
	}

	/// Factory and DI: Initializes all applications layers
	userDao := users.NewUserDao(db)
	userService := services.NewUserService(userDao)
	userController := controllers.NewUserController(userService)

	/// Maps urls to controllers
	mapUrl(userController)

	/// Starts the server
	router.Run(":8080")
}

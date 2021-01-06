package app

import (
	"fmt"
	"github.com/Abacode7/bookstore_users-api/controllers"
	"github.com/Abacode7/bookstore_users-api/datasources/mysql"
	"github.com/Abacode7/bookstore_users-api/domain/users"
	"github.com/Abacode7/bookstore_users-api/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var router = gin.Default()

func StartApplication() {
	/// Load .env config file into the os
	err := godotenv.Load()
	if err != nil{
		log.Fatalln(err)
	}

	/// Get environment variables from os
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	/// Generates connections string and opens the connection
	/// using the provided data source
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, sqlErr := mysql.Init(dataSource)
	if sqlErr != nil {
		log.Fatalln(sqlErr)
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

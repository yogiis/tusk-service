package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogiis/tusk-service/config"
	"github.com/yogiis/tusk-service/controller"
	"github.com/yogiis/tusk-service/models"
)

func main() {
	//Database
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnAccount(db)

	//Controller
	userController := controller.UserController{DB: db}

	//Router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome Tusk API sevice")
	})

	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)

	router.Static("/media", "./media")
	router.Run("192.168.1.23:8080")
}

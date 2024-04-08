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
	taskController := controller.TaskController{DB: db}

	//Router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome Tusk API sevice")
	})

	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.DELETE("/users/:id", userController.Delete)
	router.GET("/users/employee", userController.Getemployee)

	router.POST("/tasks", taskController.Create)
	router.DELETE("/tasks/:id", taskController.Delete)
	router.PATCH("/tasks/:id/submit", taskController.Submit)

	router.Static("/media", "./media")
	router.Run("192.168.1.21:8080")
}

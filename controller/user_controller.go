package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogiis/tusk-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (u *UserController) Login(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	password := user.Password

	err = u.DB.Where("email=?", user.Email).Take(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is Wrong"})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is Wrong"})
		return
	}

	c.JSON(http.StatusOK, user)

}

func (u *UserController) CreateAccount(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	emailExist := u.DB.Where("email=?", user.Email).First(&user).RowsAffected != 0
	if emailExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exist"})
		return
	}

	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	user.Password = string(hashedPasswordBytes)
	user.Role = "Employee"

	err = u.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := u.DB.Delete(&models.User{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Deleted")
}

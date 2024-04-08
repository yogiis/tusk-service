package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yogiis/tusk-service/models"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

func (t *TaskController) Create(c *gin.Context) {
	task := models.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = t.DB.Create(&task).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (t *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")
	task := models.Task{}

	if err := t.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	err := t.DB.Delete(&models.Task{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task.Attachment != "" {
		os.Remove("attachments/" + task.Attachment)
	}

	c.JSON(http.StatusOK, "Deleted")
}

func (t *TaskController) Submit(c *gin.Context) {
	task := models.Task{}
	id := c.Param("id")
	submitDate := c.PostForm("submitDate")
	file, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := t.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	//remove old attachment
	attachment := task.Attachment
	fileInfo, _ := os.Stat("attachments/" + attachment)
	if fileInfo != nil {
		os.Remove("attachments/" + attachment)
	}

	//create new attachment
	attachment = file.Filename
	err = c.SaveUploadedFile(file, "attachments/"+attachment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = t.DB.Where("id=?", id).Updates(models.Task{
		Status:     "Review",
		SubmitDate: submitDate,
		Attachment: attachment,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task.Attachment != "" {
		os.Remove("attachments/" + task.Attachment)
	}

	c.JSON(http.StatusOK, "Submite to Review")
}

func (u *TaskController) Getemployee(c *gin.Context) {
	tasks := models.Task{}

	err := u.DB.Select("id,name").Where("role=?", "Employee").Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

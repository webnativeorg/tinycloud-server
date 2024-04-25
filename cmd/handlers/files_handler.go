package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/webnativeorg/tinycloud-server/cmd/database"
	"github.com/webnativeorg/tinycloud-server/cmd/services"
	"github.com/webnativeorg/tinycloud-server/cmd/share"
)

func UploadFilesHandler(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	base_path := c.PostForm("path")
	user_interface, exist := c.Get("user")
	user := user_interface.(share.UserContext)
	if !exist {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	responseChanel := make(chan []*database.File)
	go func() {
		responseChanel <- services.UploadFiles(files, user, base_path, c)
	}()

	result := <-responseChanel
	if result == nil {
		c.JSON(500, gin.H{"message": "Error saving files"})
		return
	}
	c.JSON(200, gin.H{"message": "Files saved successfully", "files": result})
}

package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webnativeorg/tinycloud-server/cmd/database"
	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"github.com/webnativeorg/tinycloud-server/cmd/share"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadFile(file *multipart.FileHeader, path string, base_path string, user_id primitive.ObjectID) (*database.File, error) {

	file_record := database.File{
		Name:      file.Filename,
		UserId:    user_id,
		UserPath:  base_path,
		Path:      path,
		IsDir:     false,
		Size:      file.Size,
		CreatedAt: time.Now().Unix(),
	}
	err := database.CreateFile(file_record)
	if err != nil {
		return nil, err
	}
	return &file_record, nil
}

func UploadFiles(files []*multipart.FileHeader, user share.UserContext, base_path string, c *gin.Context) []*database.File {
	var files_response []*database.File
	for _, file := range files {
		path := environment.DATA_DIR + "/" + user.Id.Hex() + "/" + file.Filename

		file_record, err := UploadFile(file, path, base_path, user.Id)
		if err != nil {
			fmt.Println("Error saving file: ", err)
			return nil
		}
		files_response = append(files_response, file_record)
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			fmt.Println("Error saving file: ", err)
			return nil
		}
	}
	return files_response
}

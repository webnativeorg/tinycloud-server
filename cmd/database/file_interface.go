package database

import (
	"context"

	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	UserId    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	UserPath  string             `bson:"user_path,omitempty" json:"user_path"`
	Path      string             `bson:"path,omitempty" json:"path"`
	IsDir     bool               `bson:"is_dir" json:"is_dir"`
	Size      int64              `bson:"size" json:"size"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}

var FILES_COLLECTION_NAME string = environment.MONGO_COLLECTION_PREFIX + "files"

func CreateFile(file File) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	_, err = collection.InsertOne(context.TODO(), file)
	if err != nil {
		return err
	}
	return nil
}

func CreateFileBulk(files []File) error {
	client, err := GetMongoClient()
	bulkFiles := make([]interface{}, len(files))

	if err != nil {
		return err
	}

	for i, file := range files {
		bulkFiles[i] = file
	}

	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	_, err = collection.InsertMany(context.TODO(), bulkFiles)
	if err != nil {
		return err
	}
	return nil
}

func GetFilesByUserId(userId primitive.ObjectID) ([]File, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	cursor, err := collection.Find(context.TODO(), File{UserId: userId})
	if err != nil {
		return nil, err
	}
	var files []File
	if err = cursor.All(context.Background(), &files); err != nil {
		return nil, err
	}
	return files, nil
}
func GetFileById(id primitive.ObjectID) (File, error) {
	var file File
	client, err := GetMongoClient()
	if err != nil {
		return file, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	err = collection.FindOne(context.TODO(), File{ID: id}).Decode(&file)
	if err != nil {
		return file, err
	}
	return file, nil
}
func GetFileByPath(path string) (File, error) {
	var file File
	client, err := GetMongoClient()
	if err != nil {
		return file, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	err = collection.FindOne(context.TODO(), File{Path: path}).Decode(&file)
	if err != nil {
		return file, err
	}
	return file, nil
}
func UpdateFileById(id primitive.ObjectID, file File) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	_, err = collection.UpdateOne(context.TODO(), File{ID: id}, file)
	if err != nil {
		return err
	}
	return nil
}
func DeleteFileById(id primitive.ObjectID) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(FILES_COLLECTION_NAME)
	_, err = collection.DeleteOne(context.TODO(), File{ID: id})
	if err != nil {
		return err
	}
	return nil
}

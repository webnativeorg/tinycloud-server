package database

import (
	"context"

	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Raw user struct
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	LastName  string             `bson:"last_name,omitempty" json:"last_name"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Password  string             `bson:"password,omitempty" json:"password"`
	Status    int                `bson:"status" json:"status"` // 0: active, 1: not active
	Biography string             `bson:"biography,omitempty" json:"biography"`
	Birthday  string             `bson:"birthday,omitempty" json:"birthday"`
	Avatar    string             `bson:"avatar,omitempty" json:"avatar"`
	Phone     string             `bson:"phone,omitempty" json:"phone"`
	IsAdmin   bool               `bson:"is_admin" json:"is_admin"`
}
type UserResponse struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	LastName  string             `bson:"last_name,omitempty" json:"last_name"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Status    int                `bson:"status,omitempty" json:"status"` // 0: active, 1: not active
	Biography string             `bson:"biography,omitempty" json:"biography"`
	Birthday  string             `bson:"birthday,omitempty" json:"birthday"`
	Avatar    string             `bson:"avatar,omitempty" json:"avatar"`
	Phone     string             `bson:"phone,omitempty" json:"phone"`
	IsAdmin   bool               `bson:"is_admin,omitempty" json:"is_admin"`
}

var USER_COLLECTION_NAME string = environment.MONGO_COLLECTION_PREFIX + "users"

func CreateUser(user User) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func CreateUserBulk(users []User) error {
	client, err := GetMongoClient()
	bulkUsers := make([]interface{}, len(users))

	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	_, err = collection.InsertMany(context.TODO(), bulkUsers)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	client, err := GetMongoClient()
	if err != nil {
		return user, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	err = collection.FindOne(context.TODO(), User{Email: email}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
func GetAllUsers(skip, limit int64) ([]UserResponse, error) {
	filter := bson.D{{}}
	users := []UserResponse{}
	client, err := GetMongoClient()
	if err != nil {
		return users, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)
	findOptions.SetProjection(bson.M{"password": 0})

	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return users, err
	}
	for cur.Next(context.TODO()) {
		var user UserResponse
		err := cur.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer cur.Close(context.TODO())
	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

func UpdateUserById(id string, user User) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserById(id string) error {
	client, err := GetMongoClient()
	if err != nil {
		return err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func ExistsUserByEmail(email string) (bool, error) {
	client, err := GetMongoClient()
	if err != nil {
		return false, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	count, err := collection.CountDocuments(context.TODO(), bson.M{"email": email})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func CountUsers() (int64, error) {
	client, err := GetMongoClient()
	if err != nil {
		return 0, err
	}
	collection := client.Database(environment.MONGO_DATABASE).Collection(USER_COLLECTION_NAME)
	count, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return 0, err
	}
	return count, nil
}

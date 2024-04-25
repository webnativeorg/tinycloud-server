package share

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserContext struct {
	Id      primitive.ObjectID `json:"id"`
	Email   string             `json:"email"`
	Name    string             `json:"password"`
	IsAdmin bool               `json:"is_admin"`
}

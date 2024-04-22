package services

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/webnativeorg/tinycloud-server/cmd/database"
	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

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
type LoginInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	KeepAlive bool   `json:"keep_alive"`
}
type RegisterInput struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if input.Password == "" {
		c.JSON(400, gin.H{"message": "Password is required", "code": "password_required"})
		return
	}
	if input.Email == "" {
		c.JSON(400, gin.H{"message": "Email is required", "code": "email_required"})
		return
	}
	exist, err := database.ExistsUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error checking email", "code": "email_check_error"})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{"message": "email already exists", "code": "email_exists"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error hashing password", "code": "hash_error"})
		return
	}
	user := database.User{
		ID:       primitive.NewObjectID(),
		Name:     input.Name,
		LastName: input.LastName,
		Email:    input.Email,
		Password: string(hash),
		Status:   0,
		IsAdmin:  false,
	}
	err = database.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error creating user", "code": "user_creation_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Login(c *gin.Context) {
	var loginInput LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if loginInput.Password == "" {
		c.JSON(400, gin.H{"message": "Password is required", "code": "password_required"})
		return
	}
	if loginInput.Email == "" {
		c.JSON(400, gin.H{"message": "Email is required", "code": "email_required"})
		return
	}
	dbUser, err := database.GetUserByEmail(loginInput.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials", "code": "invalid_credentials"})
		return
	}
	if !ValidatePassword(loginInput.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials", "code": "invalid_credentials"})
		return
	}
	var token *jwt.Token
	if loginInput.KeepAlive {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       dbUser.ID.Hex(),
			"email":    dbUser.Email,
			"is_admin": dbUser.IsAdmin,
			"name":     dbUser.Name,
			"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       dbUser.ID.Hex(),
			"email":    dbUser.Email,
			"is_admin": dbUser.IsAdmin,
			"name":     dbUser.Name,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	}

	tokenString, err := token.SignedString([]byte(environment.JWT_SECRET))
	if err != nil {
		fmt.Println("Error generating token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error generating token", "code": "token_generation_error"})
		return
	}
	// Remove password from response
	userResponse := UserResponse{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
		Status:    dbUser.Status,
		Biography: dbUser.Biography,
		Birthday:  dbUser.Birthday,
		Avatar:    dbUser.Avatar,
		Phone:     dbUser.Phone,
		IsAdmin:   dbUser.IsAdmin,
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": tokenString, "user": userResponse})
}

func GetUsers(c *gin.Context) {
	var limit int64 = 10
	var skip int64 = 0

	if c.Query("$limit") != "" {
		limitInt, err := strconv.Atoi(c.Query("$limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid limit value", "code": "invalid_limit"})
			return
		}
		limit = int64(limitInt)
	}
	if c.Query("$skip") != "" {
		skipInt, err := strconv.Atoi(c.Query("$skip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid skip value", "code": "invalid_skip"})
			return
		}
		skip = int64(skipInt)
	}
	users, err := database.GetAllUsers(
		skip,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting users", "code": "get_users_error"})
		return
	}
	total, err := database.CountUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting users count", "code": "get_users_count_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "limit": limit, "skip": skip, "total": total})
}

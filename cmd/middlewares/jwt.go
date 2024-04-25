package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"github.com/webnativeorg/tinycloud-server/cmd/share"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != token.Method {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(environment.JWT_SECRET), nil
		})
		if err != nil {
			fmt.Println("Error parsing token", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id, _ := primitive.ObjectIDFromHex(claims["id"].(string))
			c.Set("user", share.UserContext{
				Email:   claims["email"].(string),
				Id:      id,
				Name:    claims["name"].(string),
				IsAdmin: claims["is_admin"].(bool),
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	_ "online-store/core/domain/entities"
	"online-store/core/domain/repositories"
	"online-store/package/config"
	"online-store/package/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

// RequireAuthMiddleware is a middleware that checks if the request is authenticated.
func RequireAuthMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request has an Authorization header
		if c.GetHeader("Authorization") == "" {
			response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}
		// config
		cfg := config.GetConfig()

		// get the token from the Authorization header bearer
		tokenString := c.GetHeader("Authorization")[7:]
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check if the signing method is correct
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key used to sign the token
			return []byte(cfg.JWTSecret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["sub"], claims["exp"])

			// check exp time
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				response.RespondError(c, http.StatusUnauthorized, "Token expired", nil)
				c.Abort()
				return
			}

			// check user id sub
			userID := claims["sub"].(string)
			userRepo := repositories.NewCustomerRepository()
			user, err := userRepo.GetByID(c, db, userID)
			if err != nil {
				response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
				c.Abort()
				return
			}

			c.Set("user", user)




		} else {
			fmt.Println(err)
			response.RespondError(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
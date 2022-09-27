package middleware

import (
	"net/http"
	"strings"

	"danain/auth"
	"danain/helper"
	"danain/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		/* ambil nilai header Authorization: Bearer tokentokentoken */
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		/* dari header Authorization, kita ambil tokennya aja */
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		/* validasi token melalui auth service */
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		/* ambil nilai user_id */
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		/* casting nilai user_id dari payload menjadi int */
		userId := int(claim["user_id"].(float64))

		/* ambil user berdasarkan user_id dari payload jwt melalui userService */
		user, err := userService.GetUserByID(userId)
		if err != nil {
			if err != nil {
				response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}

		/* set global context yang isinya data user */
		c.Set("currentUser", user)
	}
}

package main

import (
	"bwastartup/auth"
	"bwastartup/campaigns"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	userAccount := helper.GoDotEnvVariable("USER")
	password := helper.GoDotEnvVariable("PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local", userAccount, password)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Database connected")

	userRepository := user.NewRepository(db)
	campaignRepository := campaigns.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewJWTService()
	campaignService := campaigns.NewService(campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	r := gin.Default()

	// servers static assets
	r.Static("/images", "./images")

	// api versioning
	api := r.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", authMiddleware(authService, userService), campaignHandler.GetCampaignsHandler)

	r.Run(":8080")
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
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

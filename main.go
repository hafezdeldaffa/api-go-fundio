package main

import (
	"bwastartup/auth"
	"bwastartup/auth/middleware"
	"bwastartup/campaigns"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewJWTService()
	campaignService := campaigns.NewService(campaignRepository)
	transactionService := transaction.NewService(transactionRepository, campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	r := gin.Default()

	// servers static assets
	r.Static("/images", "./images")

	// api versioning
	api := r.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.POST("/campaigns", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaignHandler)
	api.POST("/campaign-images", middleware.AuthMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns", campaignHandler.GetCampaignsHandler)
	api.GET("/campaigns/:id", campaignHandler.GetDetailCampaignHandler)
	api.GET("/campaigns/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)

	api.PUT("/campaigns/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)

	r.Run(":8080")
}

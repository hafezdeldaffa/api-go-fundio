package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:Hafezd123@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Database connected")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	// api versioning
	api := r.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)

	r.Run(":8080")
}

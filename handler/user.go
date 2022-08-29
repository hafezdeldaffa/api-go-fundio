package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	/*
		- dapet input dari user
		- hander : mapping input dari user ke struct RegisterUserInput
		- service : melakukan mapping dari struct RegisterUserInput ke struct User,
					lalu ngirim struct User ke repository
		- repository : dapet data struct user buat di save ke db
	*/

	// tangkap input daru  user
	// map input dari user ke struct RegisterUserInput
	// struct diatas di passing sebagai parameter service

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to bind the json input", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Failed to register account", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userJson := user.FormatUser(newUser, "tokentokentokentoken")

	response := helper.APIResponse("Account has been created", http.StatusOK, "success", userJson)

	c.JSON(http.StatusOK, response)
}

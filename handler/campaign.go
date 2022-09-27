package handler

import (
	"danain/campaigns"
	"danain/helper"
	"danain/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter
// handler ke service
// service yg menentukan repository mana yg dipanggil
// repository akses ke db: Find All Campaign & Find Campaign by user id

type campaignHandler struct {
	service campaigns.Service
}

func NewCampaignHandler(service campaigns.Service) *campaignHandler {
	return &campaignHandler{service}
}

// /api/v1/campaigns
func (h *campaignHandler) GetCampaignsHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaign, err := h.service.GetCampaigns(userID)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to get the campaigns", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignJSON := campaigns.FormatCampaigns(campaign)

	response := helper.APIResponse("Successfully get the campaigns", http.StatusOK, "success", campaignJSON)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetDetailCampaignHandler(c *gin.Context) {
	/*
		api/v1/2
		handler: mapping id  yg di url ke struct input => service, call formatter
		service: input struct => menangkap id di url, manggil repo
		repository: get campaign by id ke db
	*/

	var input campaigns.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to get the detail campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignById(input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to get the detail campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetailJSON := campaigns.FormatCampaignDetail(campaignDetail)

	response := helper.APIResponse("Successfully get the detail campaign", http.StatusOK, "success", campaignDetailJSON)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaignHandler(c *gin.Context) {
	var input campaigns.CreateCampaignInput
	currentUser := c.MustGet("currentUser").(user.User)

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.User = currentUser

	campaign, err := h.service.CreateCampaign(input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	responseJSON := campaigns.FormatCampaign(campaign)

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", responseJSON)
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaigns.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaigns.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		errorMessage := gin.H{
			"errors": err.Error(),
		}

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dataJSON := campaigns.FormatCampaign(updatedCampaign)

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", dataJSON)
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaigns.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
			"errors":      err.Error(),
		}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
			"errors":      err.Error(),
		}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	/* ambil current usernya dari global context */
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	input.User = currentUser

	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
			"errors":      err.Error(),
		}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
			"errors":      err.Error(),
		}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}

	response := helper.APIResponse("Success to upload campaign image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

package handler

import (
	"bwastartup/campaigns"
	"bwastartup/helper"
	"bwastartup/user"
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

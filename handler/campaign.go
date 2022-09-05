package handler

import (
	"bwastartup/campaigns"
	"bwastartup/helper"
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

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get the campaigns", http.StatusNotFound, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Successfully get the campaigns", http.StatusOK, "success", campaigns)
	c.JSON(http.StatusOK, response)
}

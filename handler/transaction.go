package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{
			"errors": errors,
		}

		response := helper.APIResponse("Failed to get the campaign transactions", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		errorMessage := gin.H{
			"errors": err.Error(),
		}

		response := helper.APIResponse("Failed to get the campaign transactions", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	responseJSON := transaction.FormatCampaignTransactions(transactions)
	response := helper.APIResponse("Success to get the campaign transactions", http.StatusOK, "success", responseJSON)
	c.JSON(http.StatusOK, response)
}

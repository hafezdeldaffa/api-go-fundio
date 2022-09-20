package handler

import (
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service        transaction.Service
	paymentService payment.Service
}

func NewTransactionHandler(service transaction.Service, paymentService payment.Service) *transactionHandler {
	return &transactionHandler{service, paymentService}
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

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	transactions, err := h.service.GetTransactionsByUserID(currentUser.ID)
	if err != nil {
		errorMessage := gin.H{
			"errors": err.Error(),
		}

		response := helper.APIResponse("Failed to get the list of user's transactions", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	responseJSON := transaction.FormatUserTransactions(transactions)
	response := helper.APIResponse("Success to get the list of user's transactions", http.StatusOK, "success", responseJSON)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{
			"errors": err,
		}

		response := helper.APIResponse("Failed to create the transaction", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		errorMessage := gin.H{
			"errors": err,
		}

		response := helper.APIResponse("Failed to create the transaction", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	responseJSON := transaction.FormatTransaction(newTransaction)

	response := helper.APIResponse("Success to create the transaction", http.StatusOK, "success", responseJSON)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{
			"errors": err,
		}

		response := helper.APIResponse("Failed to process the notification", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.ProcessPayment(input)
	if err != nil {
		errorMessage := gin.H{
			"errors": err,
		}

		response := helper.APIResponse("Failed to process the notification", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, input)
}

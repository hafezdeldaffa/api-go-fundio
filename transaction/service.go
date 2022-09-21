package transaction

import (
	"bwastartup/campaigns"
	"bwastartup/payment"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionsByUserID(UserID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
}

type service struct {
	repository         Repository
	campaignRepository campaigns.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaigns.Repository, payment payment.Service) *service {
	return &service{repository, campaignRepository, payment}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindCampaignById(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not authorized to get the campaign transaction")
	}

	transactions, err := s.repository.FindTransactionsByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserID(UserID int) ([]Transaction, error) {
	transactions, err := s.repository.FindTransactionByUserID(UserID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}

	rand.Seed(time.Now().UTC().UnixNano())

	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = fmt.Sprintf("TRX-%d-%d-%d", rand.Int(), input.CampaignID, input.User.ID)

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentUrl, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentUrl

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.FindByID(transaction_id)
	if err != nil {
		return err
	}

	if input.TransactionStatus == "capture" {
		if input.PaymentType == "credit_card" {
			if input.FraudStatus == "challenge" {
				transaction.Status = "pending"
			} else {
				transaction.Status = "paid"
			}
		}
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "pending" {
		transaction.Status = "pending"
	} else if input.TransactionStatus == "deny" {
		transaction.Status = "cancelled"
	} else if input.TransactionStatus == "expire" {
		transaction.Status = "cancelled"
	} else if input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindCampaignById(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}

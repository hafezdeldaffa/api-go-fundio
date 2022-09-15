package transaction

import (
	"time"
)

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionsFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}

	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt

	return formatter
}

func FormatCampaignTransactions(transaction []Transaction) []CampaignTransactionFormatter {
	if len(transaction) == 0 {
		return []CampaignTransactionFormatter{}
	}

	var transactionsFormatter []CampaignTransactionFormatter

	for _, eachTransaction := range transaction {
		formatter := FormatCampaignTransaction(eachTransaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

func FormatUserTransaction(transaction Transaction) UserTransactionsFormatter {
	formatter := UserTransactionsFormatter{}

	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := CampaignFormatter{}

	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageURL = ""

	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = campaignFormatter

	return formatter
}

func FormatUserTransactions(transaction []Transaction) []UserTransactionsFormatter {
	if len(transaction) == 0 {
		return []UserTransactionsFormatter{}
	}

	var userTransactions []UserTransactionsFormatter

	for _, eachTransaction := range transaction {
		formatter := FormatUserTransaction(eachTransaction)
		userTransactions = append(userTransactions, formatter)
	}

	return userTransactions
}

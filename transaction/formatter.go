package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
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

package transaction

import "gorm.io/gorm"

type Repository interface {
	FindTransactionsByCampaignID(id int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactionsByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Where("campaign_id = ?", campaignID).Preload("User").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

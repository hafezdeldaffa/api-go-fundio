package transaction

import "gorm.io/gorm"

type Repository interface {
	FindTransactionsByCampaignID(id int) ([]Transaction, error)
	FindTransactionByUserID(id int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactionsByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) FindTransactionByUserID(userID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id", userID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

package transaction

import (
	"bwastartup/campaigns"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaigns.Repository
}

func NewService(repository Repository, campaignRepository campaigns.Repository) *service {
	return &service{repository, campaignRepository}
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

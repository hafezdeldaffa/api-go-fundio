package transaction

import (
	"bwastartup/campaigns"
	"bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int                `json:"id"`
	CampaignID int                `json:"campaign_id"`
	UserID     int                `json:"user_id"`
	Amount     int                `json:"amount"`
	Status     string             `json:"status"`
	Code       string             `json:"code"`
	PaymentURL string             `json:"payment_url"`
	User       user.User          `json:"user"`
	Campaign   campaigns.Campaign `json:"campaign"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

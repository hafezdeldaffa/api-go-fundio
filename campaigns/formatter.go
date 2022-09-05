package campaigns

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}

	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Title = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.ImageURL = ""
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	if len(campaigns) == 0 {
		return []CampaignFormatter{}
	}

	var campaignsFormatter []CampaignFormatter

	for _, each := range campaigns {
		campaignFormatter := FormatCampaign(each)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

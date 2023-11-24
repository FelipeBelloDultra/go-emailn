package endpoints

import "github.com/FelipeBelloDultra/emailn/internal/domain/campaign"

type Handler struct {
	CampaignService campaign.Service
}

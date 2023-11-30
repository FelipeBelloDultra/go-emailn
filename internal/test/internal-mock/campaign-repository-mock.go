package internalmock

import (
	"github.com/FelipeBelloDultra/emailn/internal/domain/campaign"
	"github.com/stretchr/testify/mock"
)

type CampaignRepositoryMock struct {
	mock.Mock
}

func (r *CampaignRepositoryMock) Save(campaign *campaign.Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *CampaignRepositoryMock) Get() ([]campaign.Campaign, error) {
	// args := r.Called(campaign)
	return nil, nil
}

func (r *CampaignRepositoryMock) GetBy(id string) (*campaign.Campaign, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.Campaign), nil
}

func (r *CampaignRepositoryMock) Delete(campaign *campaign.Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *CampaignRepositoryMock) Update(campaign *campaign.Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *CampaignRepositoryMock) GetCampaignsToBeSent() ([]campaign.Campaign, error) {
	args := r.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]campaign.Campaign), nil
}

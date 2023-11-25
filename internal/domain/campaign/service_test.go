package campaign

import (
	"errors"
	"testing"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Save(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *RepositoryMock) Get() ([]Campaign, error) {
	// args := r.Called(campaign)
	return nil, nil
}

func (r *RepositoryMock) GetBy(id string) (*Campaign, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Campaign), nil
}

var (
	newCampaign = contract.NewCampaign{
		Name:    "Test X",
		Content: "<body><h1>Hello, world!</h1></body>",
		Emails:  []string{"test@example.com"},
	}
	service = ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)

	repositoryMock := new(RepositoryMock)
	repositoryMock.On("Save", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	id, err := service.Create(newCampaign)

	assert.NotEmpty(id)
	assert.Nil(err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)

	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	repositoryMock := new(RepositoryMock)

	repositoryMock.On("Save", mock.MatchedBy(func(campaign *Campaign) bool {
		if campaign.Name != newCampaign.Name {
			return false
		}
		if campaign.Content != newCampaign.Content {
			return false
		}
		if len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)
	service.Repository = repositoryMock

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)

	repositoryMock.On("Save", mock.Anything).Return(errors.New("error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(RepositoryMock)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)
	service.Repository = repositoryMock

	campaignReturned, _ := service.GetBy(campaign.ID)

	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.Status, campaignReturned.Status)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(RepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	service.Repository = repositoryMock

	_, err := service.GetBy(campaign.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

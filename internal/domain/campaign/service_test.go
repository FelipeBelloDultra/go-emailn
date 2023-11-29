package campaign_test

import (
	"errors"
	"testing"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	"github.com/FelipeBelloDultra/emailn/internal/domain/campaign"
	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
	internalmock "github.com/FelipeBelloDultra/emailn/internal/test/internal-mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test X",
		Content:   "<body><h1>Hello, world!</h1></body>",
		Emails:    []string{"test@example.com"},
		CreatedBy: "test@example.com",
	}
	service = campaign.ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)

	repositoryMock := new(internalmock.CampaignRepositoryMock)
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
	repositoryMock := new(internalmock.CampaignRepositoryMock)

	repositoryMock.On("Save", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
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
	repositoryMock := new(internalmock.CampaignRepositoryMock)

	repositoryMock.On("Save", mock.Anything).Return(errors.New("error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)
	service.Repository = repositoryMock

	campaignReturned, _ := service.GetBy(campaign.ID)

	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.Status, campaignReturned.Status)
	assert.Equal(campaign.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	service.Repository = repositoryMock

	_, err := service.GetBy(campaign.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnErrorWhenSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)
	invalidCampaignId := "invalid"
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	service.Repository = repositoryMock

	err := service.Delete(invalidCampaignId)

	assert.Equal(err, gorm.ErrRecordNotFound)
}

func Test_Delete_ReturnInvalidStatus(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "abc-123", Status: campaign.Started}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal(err.Error(), "Campaign status invalid")
}

func Test_Delete_ReturnInternalErrorWhenDeleteHasProblem(t *testing.T) {
	assert := assert.New(t)
	createdCampaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(createdCampaign, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return createdCampaign == campaign
	})).Return(errors.New("error on delete campaign"))
	service.Repository = repositoryMock

	err := service.Delete(createdCampaign.ID)

	assert.Equal(err.Error(), internalerrors.ErrInternal.Error())
}

func Test_Delete_DeleteCampaignSuccessfully(t *testing.T) {
	assert := assert.New(t)
	createdCampaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(createdCampaign, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return createdCampaign == campaign
	})).Return(nil)
	service.Repository = repositoryMock

	err := service.Delete(createdCampaign.ID)

	assert.Nil(err)
}

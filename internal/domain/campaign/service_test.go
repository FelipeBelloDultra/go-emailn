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
	campaignPending *campaign.Campaign
	campaignStarted *campaign.Campaign
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImp{}
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock

	campaignPending, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "abc-123", Status: campaign.Started}
}

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	setUp()

	repositoryMock.On("Save", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotEmpty(id)
	assert.Nil(err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	setUp()

	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	setUp()
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

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	setUp()

	repositoryMock.On("Save", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	setUp()

	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)

	campaignReturned, _ := service.GetBy(campaignPending.ID)

	assert.Equal(campaignPending.ID, campaignReturned.ID)
	assert.Equal(campaignPending.Name, campaignReturned.Name)
	assert.Equal(campaignPending.Content, campaignReturned.Content)
	assert.Equal(campaignPending.Status, campaignReturned.Status)
	assert.Equal(campaignPending.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)
	setUp()

	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy(campaignPending.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnErrorWhenSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)
	setUp()
	invalidCampaignId := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete(invalidCampaignId)

	assert.Equal(err, gorm.ErrRecordNotFound)
}

func Test_Delete_ReturnInvalidStatus(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(err.Error(), "Campaign status invalid")
}

func Test_Delete_ReturnInternalErrorWhenDeleteHasProblem(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error on delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(err.Error(), internalerrors.ErrInternal.Error())
}

func Test_Delete_DeleteCampaignSuccessfully(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(err)
}

func Test_Start_ReturnErrorWhenSomethingWrongExists(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid")

	assert.Equal(err, gorm.ErrRecordNotFound)
}

func Test_Start_ReturnInvalidStatus(t *testing.T) {
	assert := assert.New(t)
	setUp()

	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Start(campaignStarted.ID)

	assert.Equal(err.Error(), "Campaign status invalid")
}

func Test_Start_ShouldSendMail(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	emailWasSent := false
	service.SendMail = func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignPending.ID {
			emailWasSent = true
		}

		return nil
	}

	service.Start(campaignPending.ID)

	assert.True(emailWasSent)
}

func Test_Start_ShouldErrorWhenSendEmailFail(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	service.SendMail = func(campaign *campaign.Campaign) error {
		return errors.New("error to send email")
	}

	err := service.Start(campaignPending.ID)

	assert.Equal(internalerrors.ErrInternal, err)
}

func Test_Start_WhenSendEmailMustChangeStatusToDone(t *testing.T) {
	assert := assert.New(t)
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID &&
			campaignToUpdate.Status == campaign.Done
	})).Return(nil)
	service.SendMail = func(campaign *campaign.Campaign) error {
		return nil
	}

	service.Start(campaignPending.ID)

	assert.Equal(campaign.Done, campaignPending.Status)
}

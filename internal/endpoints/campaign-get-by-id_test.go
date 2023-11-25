package endpoints

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	"github.com/FelipeBelloDultra/emailn/internal/domain/campaign"
	internalmock "github.com/FelipeBelloDultra/emailn/internal/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignGetById_ShouldReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaign := contract.CampaignResponse{
		ID:      "123",
		Name:    "Test name",
		Content: "Hi everyone",
		Status:  campaign.Pending,
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("GetBy", mock.Anything).Return(&campaign, nil)
	handler := Handler{CampaignService: service}

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	response, status, _ := handler.CampaignGetById(res, req)

	assert.Equal(200, status)
	assert.Equal(campaign.ID, response.(*contract.CampaignResponse).ID)
	assert.Equal(campaign.Name, response.(*contract.CampaignResponse).Name)
	assert.Equal(campaign.Content, response.(*contract.CampaignResponse).Content)
	assert.Equal(campaign.Status, response.(*contract.CampaignResponse).Status)
}

func Test_CampaignGetById_ShouldReturnErrorWhenSomethinWrong(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errExpected := errors.New("something went wrong")
	service.On("GetBy", mock.Anything).Return(nil, errExpected)
	handler := Handler{CampaignService: service}

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	_, _, errReturned := handler.CampaignGetById(res, req)

	assert.Equal(errReturned.Error(), errExpected.Error())
}

package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	internalmock "github.com/FelipeBelloDultra/emailn/internal/test/mock"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignPost_ShouldSaveNewCampaign(t *testing.T) {
	assert := assert.New(t)
	body := contract.NewCampaign{
		Name:    "Service",
		Content: "Hi everyone",
		Emails:  []string{"user@example.com"},
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.MatchedBy(func(req contract.NewCampaign) bool {
		return req.Name == body.Name && req.Content == body.Content
	})).Return("123x", nil)
	handler := Handler{CampaignService: service}

	var buff bytes.Buffer
	json.NewEncoder(&buff).Encode(body)

	req, _ := http.NewRequest("POST", "/", &buff)
	res := httptest.NewRecorder()

	_, status, err := handler.CampaignPost(res, req)

	assert.Equal(201, status)
	assert.Nil(err)
}

func Test_CampaignPost_ShouldInformErrorWhenExists(t *testing.T) {
	assert := assert.New(t)
	body := contract.NewCampaign{}
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := Handler{CampaignService: service}

	var buff bytes.Buffer
	json.NewEncoder(&buff).Encode(body)

	req, _ := http.NewRequest("POST", "/", &buff)
	res := httptest.NewRecorder()

	_, _, err := handler.CampaignPost(res, req)

	assert.NotNil(err)
}
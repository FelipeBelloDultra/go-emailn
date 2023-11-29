package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	internalmock "github.com/FelipeBelloDultra/emailn/internal/test/internal-mock"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup(body contract.NewCampaign, createdByExpected string) (*http.Request, *httptest.ResponseRecorder) {
	var buff bytes.Buffer
	json.NewEncoder(&buff).Encode(body)

	req, _ := http.NewRequest("POST", "/", &buff)
	ctx := context.WithValue(req.Context(), "email", createdByExpected)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()

	return req, res
}

func Test_CampaignPost_ShouldSaveNewCampaign(t *testing.T) {
	assert := assert.New(t)
	createdByExpected := "test@example.com"
	body := contract.NewCampaign{
		Name:    "Service",
		Content: "Hi everyone",
		Emails:  []string{"user@example.com"},
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.MatchedBy(func(req contract.NewCampaign) bool {
		return req.Name == body.Name &&
			req.Content == body.Content &&
			req.CreatedBy == createdByExpected
	})).Return("123x", nil)
	handler := Handler{CampaignService: service}

	req, res := setup(body, createdByExpected)

	_, status, err := handler.CampaignPost(res, req)

	assert.Equal(201, status)
	assert.Nil(err)
}

func Test_CampaignPost_ShouldInformErrorWhenExists(t *testing.T) {
	assert := assert.New(t)
	createdByExpected := "test@example.com"
	body := contract.NewCampaign{}
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := Handler{CampaignService: service}

	req, res := setup(body, createdByExpected)

	_, _, err := handler.CampaignPost(res, req)

	assert.NotNil(err)
}

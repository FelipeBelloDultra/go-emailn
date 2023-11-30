package endpoints

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalmock "github.com/FelipeBelloDultra/emailn/internal/test/internal-mock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Start_200(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	campaignId := "123-abc"
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == campaignId
	})).Return(nil)
	handler := Handler{CampaignService: service}

	req, _ := http.NewRequest("PATCH", "/", nil)

	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add("id", campaignId)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiContext))

	res := httptest.NewRecorder()

	_, status, err := handler.CampaignStart(res, req)

	assert.Equal(200, status)
	assert.Nil(err)
}

func Test_Start_Err(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	expectedError := errors.New("something wrong")
	service.On("Start", mock.Anything).Return(expectedError)
	handler := Handler{CampaignService: service}

	req, _ := http.NewRequest("PATCH", "/", nil)
	res := httptest.NewRecorder()

	_, _, err := handler.CampaignStart(res, req)

	assert.Equal(err, expectedError)
}

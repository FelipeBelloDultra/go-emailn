package endpoints

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Start_200(t *testing.T) {
	setup()
	assert := assert.New(t)
	campaignId := "123-abc"
	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == campaignId
	})).Return(nil)

	req, res := newHttpTest("PATCH", "/", nil)

	req = addParameter(req, "id", campaignId)

	_, status, err := handler.CampaignStart(res, req)

	assert.Equal(200, status)
	assert.Nil(err)
}

func Test_Start_Err(t *testing.T) {
	setup()
	assert := assert.New(t)
	expectedError := errors.New("something wrong")
	service.On("Start", mock.Anything).Return(expectedError)

	req, res := newHttpTest("PATCH", "/", nil)

	_, _, err := handler.CampaignStart(res, req)

	assert.Equal(err, expectedError)
}

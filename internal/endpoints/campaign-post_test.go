package endpoints

import (
	"fmt"
	"testing"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	createdByExpected = "test@example.com"
	body              = contract.NewCampaign{
		Name:    "Service",
		Content: "Hi everyone",
		Emails:  []string{"user@example.com"},
	}
)

func Test_CampaignPost_201(t *testing.T) {
	setup()
	assert := assert.New(t)
	service.On("Create", mock.MatchedBy(func(req contract.NewCampaign) bool {
		return req.Name == body.Name &&
			req.Content == body.Content &&
			req.CreatedBy == createdByExpected
	})).Return("123x", nil)

	req, res := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, status, err := handler.CampaignPost(res, req)

	assert.Equal(201, status)
	assert.Nil(err)
}

func Test_CampaignPost_Err(t *testing.T) {
	setup()
	assert := assert.New(t)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))

	req, res := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, _, err := handler.CampaignPost(res, req)

	assert.NotNil(err)
}

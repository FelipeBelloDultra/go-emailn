package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name      = "Campaign X"
	content   = "<body><h1>Hello, world!</h1></body>"
	contacts  = []string{"john@doe.com"}
	createdBy = "test@test.com.br"
	fake      = faker.New()
)

func Test_Campaign_NewCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
	assert.Equal(campaign.CreatedBy, createdBy)
}

func Test_Campaign_IDIsNotEmpty(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.NotEmpty(campaign.ID)
}

func Test_Campaign_MustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Greater(campaign.CreatedOn, now)
}

func Test_Campaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign("", content, contacts, createdBy)

	assert.Equal("name is required with min 5", error.Error())
}

func Test_Campaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(fake.Lorem().Text(25), content, contacts, createdBy)

	assert.Equal("name is required with max 24", error.Error())
}

func Test_Campaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, "", contacts, createdBy)

	assert.Equal("content is required with min 5", error.Error())
}

func Test_Campaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, fake.Lorem().Text(1040), contacts, createdBy)

	assert.Equal("content is required with max 1024", error.Error())
}

func Test_Campaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, content, nil, createdBy)

	assert.Equal("contacts is required with min 1", error.Error())
}

func Test_Campaign_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, content, []string{"invalid-email"}, createdBy)

	assert.Equal("email is invalid", error.Error())
}

func Test_Campaign_MustStatusStartWithPending(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(Pending, campaign.Status)
}

func Test_Campaign_MustValidateEmail(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, contacts, "")

	assert.Equal("createdby is invalid", err.Error())
}

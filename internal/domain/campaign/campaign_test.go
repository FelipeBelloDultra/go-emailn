package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name     = "Campaign X"
	content  = "<body><h1>Hello, world!</h1></body>"
	contacts = []string{"john@doe.com"}
	fake     = faker.New()
)

func Test_Campaign_NewCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
}

func Test_Campaign_IDIsNotEmpty(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.NotEmpty(campaign.ID)
}

func Test_Campaign_MustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.Greater(campaign.CreatedOn, now)
}

func Test_Campaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign("", content, contacts)

	assert.Equal("name is required with min 5", error.Error())
}

func Test_Campaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(fake.Lorem().Text(25), content, contacts)

	assert.Equal("name is required with max 24", error.Error())
}

func Test_Campaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, "", contacts)

	assert.Equal("content is required with min 5", error.Error())
}

func Test_Campaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, fake.Lorem().Text(1040), contacts)

	assert.Equal("content is required with max 1024", error.Error())
}

func Test_Campaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, content, nil)

	assert.Equal("contacts is required with min 1", error.Error())
}

func Test_Campaign_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, content, []string{"invalid-email"})

	assert.Equal("email is invalid", error.Error())
}

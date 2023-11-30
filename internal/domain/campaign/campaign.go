package campaign

import (
	"time"

	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
	"github.com/google/uuid"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Email      string `validate:"email" gorm:"size:100"`
	CampaignID string `gorm:"size:50"`
}

const (
	Pending  = "Pending"
	Canceled = "Canceled"
	Deleted  = "Deleted"
	Started  = "Started"
	Fail     = "Fail"
	Done     = "Done"
)

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50;not null"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100;not null"`
	UpdatedOn time.Time `validate:"required"`
	CreatedOn time.Time `validate:"required" gorm:"not null"`
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024;not null"`
	Contacts  []Contact `validate:"min=1,dive"`
	Status    string    `gorm:"size:20;not null"`
	CreatedBy string    `validate:"email" gorm:"size:50;not null"`
}

func (c *Campaign) Cancel() {
	c.UpdatedOn = time.Now()
	c.Status = Canceled
}

func (c *Campaign) Delete() {
	c.UpdatedOn = time.Now()
	c.Status = Deleted
}

func (c *Campaign) Done() {
	c.UpdatedOn = time.Now()
	c.Status = Done
}

func (c *Campaign) Fail() {
	c.UpdatedOn = time.Now()
	c.Status = Fail
}

func (c *Campaign) Started() {
	c.UpdatedOn = time.Now()
	c.Status = Started
}

func NewCampaign(name string, content string, emails []string, createdBy string) (*Campaign, error) {
	contacts := make([]Contact, len(emails))
	for index, value := range emails {
		contacts[index].Email = value
		contacts[index].ID = uuid.New().String()
	}

	campaign := &Campaign{
		ID:        uuid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(),
		Contacts:  contacts,
		Status:    Pending,
		CreatedBy: createdBy,
	}

	err := internalerrors.ValidateStruct(campaign)

	if err == nil {
		return campaign, nil
	}

	return nil, err
}

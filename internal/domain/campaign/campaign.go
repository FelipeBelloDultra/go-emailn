package campaign

import (
	"time"

	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
	"github.com/google/uuid"
)

type Contact struct {
	Email string `validate:"email"`
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
	ID        string    `validate:"required"`
	Name      string    `validate:"min=5,max=24"`
	CreatedOn time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024"`
	Contacts  []Contact `validate:"min=1,dive"`
	Status    string
}

func NewCampaign(name string, content string, emails []string) (*Campaign, error) {
	contacts := make([]Contact, len(emails))
	for index, value := range emails {
		contacts[index].Email = value
	}

	campaign := &Campaign{
		ID:        uuid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(),
		Contacts:  contacts,
		Status:    Pending,
	}

	err := internalerrors.ValidateStruct(campaign)

	if err == nil {
		return campaign, nil
	}

	return nil, err
}

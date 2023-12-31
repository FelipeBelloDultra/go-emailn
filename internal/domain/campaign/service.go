package campaign

import (
	"errors"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
)

type ServiceImp struct {
	Repository Repository
	SendMail   func(campaign *Campaign) error
}

type Service interface {
	Create(newcampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaignResponse, error)
	Delete(id string) error
	Start(id string) error
	SendEmailAndUpdateStatus(campaign *Campaign)
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	if err != nil {
		return "", err
	}

	err = s.Repository.Save(campaign)
	if err != nil {
		return "", internalerrors.ErrInternal
	}

	return campaign.ID, nil
}

func (s *ServiceImp) GetBy(id string) (*contract.CampaignResponse, error) {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return nil, internalerrors.ProcessErrorToReturn(err)
	}

	return &contract.CampaignResponse{
		ID:                   campaign.ID,
		Name:                 campaign.Name,
		Content:              campaign.Content,
		Status:               campaign.Status,
		AmountOfEmailsToSend: len(campaign.Contacts),
		CreatedBy:            campaign.CreatedBy,
	}, nil
}

func (s *ServiceImp) Delete(id string) error {
	campaign, err := s.Repository.GetBy(id)
	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	campaign.Delete()

	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil
}

func (s *ServiceImp) SendEmailAndUpdateStatus(campaign *Campaign) {
	err := s.SendMail(campaign)

	if err != nil {
		campaign.Done()
	} else {
		campaign.Fail()
	}
	s.Repository.Update(campaign)
}

func (s *ServiceImp) Start(id string) error {
	campaign, err := s.Repository.GetBy(id)
	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	go s.SendEmailAndUpdateStatus(campaign)

	campaign.Started()

	err = s.Repository.Update(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil
}

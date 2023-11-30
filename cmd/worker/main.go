package main

import (
	"log"
	"time"

	"github.com/FelipeBelloDultra/emailn/internal/domain/campaign"
	"github.com/FelipeBelloDultra/emailn/internal/infra/database"
	"github.com/FelipeBelloDultra/emailn/internal/infra/mail"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db := database.NewDb()
	repository := database.CampaignRepository{Db: db}
	campaignService := campaign.ServiceImp{
		Repository: &repository,
		SendMail:   mail.SendMail,
	}

	for {
		campaigns, err := repository.GetCampaignsToBeSent()

		if err != nil {
			println(err.Error())
		}

		println("amount of campaigns: ", len(campaigns))

		for _, campaign := range campaigns {
			campaignService.SendEmailAndUpdateStatus(&campaign)
			println("Campaign sent: ", campaign.ID)
		}

		time.Sleep(30 * time.Second)
	}
}

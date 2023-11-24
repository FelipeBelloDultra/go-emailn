package main

import (
	"errors"
	"net/http"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	"github.com/FelipeBelloDultra/emailn/internal/domain/campaign"
	"github.com/FelipeBelloDultra/emailn/internal/infra/database"
	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	service := campaign.Service{
		Repository: &database.CampaignRepository{},
	}

	r.Post("/campaigns", func(w http.ResponseWriter, r *http.Request) {
		var request contract.NewCampaign
		render.DecodeJSON(r.Body, &request)

		id, err := service.Create(request)

		if err != nil {
			if errors.Is(err, internalerrors.ErrInternal) {
				render.Status(r, 500)
				render.JSON(w, r, map[string]string{"error": err.Error()})
				return
			}

			render.Status(r, 400)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		render.Status(r, 201)
		render.JSON(w, r, map[string]string{"id": id})
	})

	http.ListenAndServe(":3000", r)
}

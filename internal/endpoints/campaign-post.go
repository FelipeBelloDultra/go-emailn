package endpoints

import (
	"errors"
	"net/http"

	"github.com/FelipeBelloDultra/emailn/internal/contract"
	internalerrors "github.com/FelipeBelloDultra/emailn/internal/internal-errors"
	"github.com/go-chi/render"
)

func (h *Handler) CampaignPost(w http.ResponseWriter, r *http.Request) {
	var request contract.NewCampaign
	render.DecodeJSON(r.Body, &request)

	id, err := h.CampaignService.Create(request)

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
}

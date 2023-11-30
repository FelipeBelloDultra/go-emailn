package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	internalmock "github.com/FelipeBelloDultra/emailn/internal/test/internal-mock"
	"github.com/go-chi/chi/v5"
)

var (
	service *internalmock.CampaignServiceMock
	handler = Handler{}
)

func setup() {
	service = new(internalmock.CampaignServiceMock)
	handler.CampaignService = service
}

func newHttpTest(method string, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var buff bytes.Buffer
	if body != nil {
		json.NewEncoder(&buff).Encode(body)

	}

	req, _ := http.NewRequest(method, url, &buff)
	res := httptest.NewRecorder()

	return req, res

}

func addParameter(req *http.Request, keyParam string, valueParam string) *http.Request {
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add(keyParam, valueParam)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiContext))
}

func addContext(req *http.Request, keyParam string, valueParam string) *http.Request {
	ctx := context.WithValue(req.Context(), keyParam, valueParam)
	req = req.WithContext(ctx)

	return req
}

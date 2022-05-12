package internalhttp

import (
	"encoding/json"
	"net/http"

	"github.com/usmartpro/banner-rotation/internal/app"
)

type ServerHandlers struct {
	app *app.App
}

func NewServerHandlers(a *app.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func (s *ServerHandlers) AddBannerToSlot(w http.ResponseWriter, r *http.Request) {
	var bS BannerSlot
	if err := json.NewDecoder(r.Body).Decode(&bS); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.app.AddBannerToSlot(bS.BannerID, bS.SlotID); err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ServerHandlers) DeleteBannerFromSlot(w http.ResponseWriter, r *http.Request) {
	var bS BannerSlot
	var err error
	if err = json.NewDecoder(r.Body).Decode(&bS); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = s.app.DeleteBannerFromSlot(bS.BannerID, bS.SlotID); err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *ServerHandlers) ClickBanner(w http.ResponseWriter, r *http.Request) {
	var bSS BannerSlotSocialGroup
	var err error
	if err = json.NewDecoder(r.Body).Decode(&bSS); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err = s.app.ClickBanner(bSS.BannerID, bSS.SlotID, bSS.SocialGroupID); err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *ServerHandlers) GetBanner(w http.ResponseWriter, r *http.Request) {
	var sSG SlotSocialGroup
	var err error

	if err = json.NewDecoder(r.Body).Decode(&sSG); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	var idBanner int
	idBanner, err = s.app.GetBanner(sSG.SlotID, sSG.SocialGroupID)

	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	result := BannerResponse{BannerID: idBanner}

	responseData, err := json.Marshal(result)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func ResponseError(w http.ResponseWriter, code int, err error) {
	data, err := json.Marshal(Error{
		false,
		err.Error(),
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to marshall error"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

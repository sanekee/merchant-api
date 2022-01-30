package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/handler/auth"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type MerchantRepo interface {
	GetAll(model.Pagination) ([]*model.Merchant, error)
	Get(string) (*model.Merchant, error)
	Insert(*model.NewMerchant) (*model.Merchant, error)
	Update(string, *model.UpdateMerchant) (*model.Merchant, error)
	Delete(string) error
}

type MerchantHandler struct {
	path         string
	merchantRepo MerchantRepo
}

func NewMerchantHandler(path string, repo MerchantRepo) *MerchantHandler {
	return &MerchantHandler{
		path:         path,
		merchantRepo: repo,
	}
}
func (m *MerchantHandler) RegisterMux(router *mux.Router) {
	r := router.PathPrefix(m.path).Subrouter().StrictSlash(true)
	r.Path("").Methods(http.MethodGet).HandlerFunc(auth.JWTAuth(m.getAll))
	r.Path("").Methods(http.MethodPost).HandlerFunc(auth.JWTAuth(m.create))
	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(auth.JWTAuth(m.get))
	r.Path("/{id}").Methods(http.MethodPut).HandlerFunc(auth.JWTAuth(m.update))
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(auth.JWTAuth(m.delete))

}

func (m *MerchantHandler) getAll(w http.ResponseWriter, r *http.Request) {
	opt, err := getPaginationFromReq(r)
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: err.Error()})
		return
	}

	merchants, err := m.merchantRepo.GetAll(opt)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.CommonResponseStatusError, Message: err.Error()})
		return
	}
	ResponseJSON(w, http.StatusOK, merchants)
}

func (m *MerchantHandler) create(w http.ResponseWriter, r *http.Request) {
	var newMerchant model.NewMerchant
	if err := json.NewDecoder(r.Body).Decode(&newMerchant); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	created, err := m.merchantRepo.Insert(&newMerchant)
	if err != nil {
		code := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrDuplicate):
			code = http.StatusConflict
		}
		ResponseJSON(w, code, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error creating merchant"})
		return
	}
	ResponseJSON(w, http.StatusCreated, created)
}

func (m *MerchantHandler) get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	merchant, err := m.merchantRepo.Get(id)
	if err != nil {
		code := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrNoResults):
			code = http.StatusNotFound
		case errors.Is(err, model.ErrRequest):
			code = http.StatusBadRequest
		}
		ResponseJSON(w, code, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error getting merchant"})
		return
	}
	ResponseJSON(w, http.StatusOK, merchant)
}

func (m *MerchantHandler) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var mc model.UpdateMerchant
	if err := json.NewDecoder(r.Body).Decode(&mc); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	updated, err := m.merchantRepo.Update(id, &mc)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrNoResults):
			statusCode = http.StatusNotFound
		case errors.Is(err, model.ErrDuplicate):
			statusCode = http.StatusConflict
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "error updating merchant"})
		return
	}
	ResponseJSON(w, http.StatusOK, updated)
}

func (m *MerchantHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := m.merchantRepo.Delete(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrNoResults):
			statusCode = http.StatusNotFound
		case errors.Is(err, model.ErrRequest):
			statusCode = http.StatusBadRequest
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error deleting merchant"})
		return
	}
	ResponseJSON(w, http.StatusOK, model.CommonResponse{Status: model.CommonResponseStatusOk, Message: "merchant deleted"})
}

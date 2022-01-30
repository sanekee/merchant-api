package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type MerchantRepo interface {
	GetAll(model.Pagination) ([]*model.Merchant, error)
	Get(string) (*model.Merchant, error)
	Insert(*model.Merchant) (*model.Merchant, error)
	Update(*model.Merchant) (*model.Merchant, error)
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
	r.Path("").Methods(http.MethodGet).HandlerFunc(m.getAll)
	r.Path("").Methods(http.MethodPost).HandlerFunc(m.create)
	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(m.get)
	r.Path("/{id}").Methods(http.MethodPut).HandlerFunc(m.update)
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(m.delete)

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
	merchant := &model.Merchant{
		Id:   uuid.NewString(),
		Code: newMerchant.Code,
	}
	created, err := m.merchantRepo.Insert(merchant)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error creating merchant"})
		return
	}
	ResponseJSON(w, http.StatusCreated, created)
}

func (m *MerchantHandler) get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	merchant, err := m.merchantRepo.Get(id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, model.ErrNoResults) {
			code = http.StatusNotFound
		}
		ResponseJSON(w, code, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error getting merchant"})
		return
	}
	ResponseJSON(w, http.StatusOK, merchant)
}

func (m *MerchantHandler) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var merchant model.Merchant
	if err := json.NewDecoder(r.Body).Decode(&merchant); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	if merchant.Id != id {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request, merchant id is different"})
		return
	}
	updated, err := m.merchantRepo.Update(&merchant)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, model.ErrNoResults) {
			statusCode = http.StatusNotFound
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
		if errors.Is(err, model.ErrNoResults) {
			statusCode = http.StatusNotFound
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error deleting merchant"})
		return
	}
	ResponseJSON(w, http.StatusOK, model.CommonResponse{Status: model.CommonResponseStatusOk, Message: "merchant deleted"})
}

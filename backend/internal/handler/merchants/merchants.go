package merchants

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sanekee/limi"
	"github.com/sanekee/merchant-api/internal/handler"
	"github.com/sanekee/merchant-api/internal/model"
)

type MerchantRepo interface {
	GetAll(model.Pagination) ([]*model.Merchant, error)
	Get(string) (*model.Merchant, error)
	Insert(*model.NewMerchant) (*model.Merchant, error)
	Update(string, *model.UpdateMerchant) (*model.Merchant, error)
	Delete(string) error
}

type Merchants struct {
	merchantRepo MerchantRepo
}

func NewMerchants(repo MerchantRepo) *Merchants {
	return &Merchants{
		merchantRepo: repo,
	}
}

func (m *Merchants) Get(w http.ResponseWriter, r *http.Request) {
	opt, err := handler.GetPaginationFromReq(r)
	if err != nil {
		handler.ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.Error, Message: err.Error()})
		return
	}

	merchants, err := m.merchantRepo.GetAll(opt)
	if err != nil {
		handler.ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.Error, Message: err.Error()})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, merchants)
}

func (m *Merchants) Post(w http.ResponseWriter, r *http.Request) {
	var newMerchant model.NewMerchant
	if err := json.NewDecoder(r.Body).Decode(&newMerchant); err != nil {
		handler.ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.Error, Message: "invalid request"})
		return
	}

	created, err := m.merchantRepo.Insert(&newMerchant)
	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, model.ErrDuplicate) {
			code = http.StatusConflict
		}

		handler.ResponseJSON(w, code, model.CommonResponse{Status: model.Error, Message: "Error creating merchant"})
		return
	}

	handler.ResponseJSON(w, http.StatusCreated, created)
}

type Merchant struct {
	limi         struct{} `path:"{id}"`
	merchantRepo MerchantRepo
}

func NewMerchant(repo MerchantRepo) Merchant {
	return Merchant{
		merchantRepo: repo,
	}
}

func (m Merchant) Get(w http.ResponseWriter, r *http.Request) {
	id := limi.GetURLParam(r.Context(), "id")

	merchant, err := m.merchantRepo.Get(id)
	if err != nil {
		code := http.StatusInternalServerError

		switch true {
		case errors.Is(err, model.ErrNoResults):
			code = http.StatusNotFound
		case errors.Is(err, model.ErrRequest):
			code = http.StatusBadRequest
		}

		handler.ResponseJSON(w, code, model.CommonResponse{Status: model.Error, Message: "Error getting merchant"})
		return
	}
	handler.ResponseJSON(w, http.StatusOK, merchant)
}

func (m Merchant) Put(w http.ResponseWriter, r *http.Request) {
	id := limi.GetURLParam(r.Context(), "id")

	var mc model.UpdateMerchant
	if err := json.NewDecoder(r.Body).Decode(&mc); err != nil {
		handler.ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.Error, Message: "invalid request"})
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

		handler.ResponseJSON(w, statusCode, model.CommonResponse{Status: model.Error, Message: "error updating merchant"})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, updated)
}

func (m Merchant) Delete(w http.ResponseWriter, r *http.Request) {
	id := limi.GetURLParam(r.Context(), "id")

	err := m.merchantRepo.Delete(id)
	if err != nil {
		statusCode := http.StatusInternalServerError

		switch true {
		case errors.Is(err, model.ErrNoResults):
			statusCode = http.StatusNotFound
		case errors.Is(err, model.ErrRequest):
			statusCode = http.StatusBadRequest
		}

		handler.ResponseJSON(w, statusCode, model.CommonResponse{Status: model.Error, Message: "Error deleting merchant"})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, model.CommonResponse{Status: model.Ok, Message: "merchant deleted"})
}

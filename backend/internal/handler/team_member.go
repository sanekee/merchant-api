package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/handler/auth"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type TeamMemberRepo interface {
	GetByMerchantID(teamMemberID string, opts model.Pagination) ([]*model.TeamMember, error)
	Get(string) (*model.TeamMember, error)
	Insert(*model.NewTeamMember) (*model.TeamMember, error)
	Update(string, *model.UpdateTeamMember) (*model.TeamMember, error)
	Delete(string) error
}

type TeamMemberHandler struct {
	path           string
	teamMemberRepo TeamMemberRepo
}

func NewTeamMemberHandler(path string, repo TeamMemberRepo) *TeamMemberHandler {
	return &TeamMemberHandler{
		path:           path,
		teamMemberRepo: repo,
	}
}
func (t *TeamMemberHandler) RegisterMux(router *mux.Router) {
	r := router.PathPrefix(t.path).Subrouter().StrictSlash(true)
	r.Path("/merchant/{merchant_id}").Methods(http.MethodGet).HandlerFunc(auth.JWTAuth(t.listByMerchantID))
	r.Path("").Methods(http.MethodPost).HandlerFunc(auth.JWTAuth(t.create))
	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(auth.JWTAuth(t.get))
	r.Path("/{id}").Methods(http.MethodPut).HandlerFunc(auth.JWTAuth(t.update))
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(auth.JWTAuth(t.delete))

}

func (t *TeamMemberHandler) listByMerchantID(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["merchant_id"]
	opt, err := getPaginationFromReq(r)
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: err.Error()})
		return
	}

	teamMembers, err := t.teamMemberRepo.GetByMerchantID(merchantID, opt)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.CommonResponseStatusError, Message: err.Error()})
		return
	}
	ResponseJSON(w, http.StatusOK, teamMembers)
}

func (t *TeamMemberHandler) create(w http.ResponseWriter, r *http.Request) {
	var newTeamMember model.NewTeamMember
	if err := json.NewDecoder(r.Body).Decode(&newTeamMember); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	created, err := t.teamMemberRepo.Insert(&newTeamMember)
	if err != nil {
		code := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrRequest):
			code = http.StatusBadRequest
		case errors.Is(err, model.ErrDuplicate):
			code = http.StatusConflict
		}
		ResponseJSON(w, code, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error creating team member"})
		return
	}
	ResponseJSON(w, http.StatusCreated, created)
}

func (t *TeamMemberHandler) get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	teamMember, err := t.teamMemberRepo.Get(id)
	if err != nil {
		code := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrNoResults):
			code = http.StatusNotFound
		}
		ResponseJSON(w, code, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error getting team member"})
		return
	}
	ResponseJSON(w, http.StatusOK, teamMember)
}

func (t *TeamMemberHandler) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var utm model.UpdateTeamMember
	if err := json.NewDecoder(r.Body).Decode(&utm); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	updated, err := t.teamMemberRepo.Update(id, &utm)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrNoResults):
			statusCode = http.StatusNotFound
		case errors.Is(err, model.ErrDuplicate):
			statusCode = http.StatusConflict
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "error updating teamMember"})
		return
	}
	ResponseJSON(w, http.StatusOK, updated)
}

func (t *TeamMemberHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := t.teamMemberRepo.Delete(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch true {
		case errors.Is(err, model.ErrNoResults):
			statusCode = http.StatusNotFound
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error deleting team member"})
		return
	}
	ResponseJSON(w, http.StatusOK, model.CommonResponse{Status: model.CommonResponseStatusOk, Message: "Team member deleted"})
}

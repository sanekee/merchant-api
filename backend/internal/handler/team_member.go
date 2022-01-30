package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

type TeamMemberRepo interface {
	GetByMerchantID(teamMemberID string, opts model.Pagination) ([]*model.TeamMember, error)
	Get(string) (*model.TeamMember, error)
	Insert(*model.TeamMember) (*model.TeamMember, error)
	Update(*model.TeamMember) (*model.TeamMember, error)
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
func (m *TeamMemberHandler) RegisterMux(router *mux.Router) {
	r := router.PathPrefix(m.path).Subrouter().StrictSlash(true)
	r.Path("/merchant/{merchant_id}").Methods(http.MethodGet).HandlerFunc(m.listByMerchantID)
	r.Path("").Methods(http.MethodPost).HandlerFunc(m.create)
	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(m.get)
	r.Path("/{id}").Methods(http.MethodPut).HandlerFunc(m.update)
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(m.delete)

}

func (m *TeamMemberHandler) listByMerchantID(w http.ResponseWriter, r *http.Request) {
	merchantID := mux.Vars(r)["merchant_id"]
	opt, err := getPaginationFromReq(r)
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: err.Error()})
		return
	}

	teamMembers, err := m.teamMemberRepo.GetByMerchantID(merchantID, opt)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.CommonResponseStatusError, Message: err.Error()})
		return
	}
	ResponseJSON(w, http.StatusOK, teamMembers)
}

func (m *TeamMemberHandler) create(w http.ResponseWriter, r *http.Request) {
	var newTeamMember model.NewTeamMember
	if err := json.NewDecoder(r.Body).Decode(&newTeamMember); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	teamMember := &model.TeamMember{
		Id:         uuid.NewString(),
		Email:      newTeamMember.Email,
		MerchantId: newTeamMember.MerchantId,
	}
	created, err := m.teamMemberRepo.Insert(teamMember)
	if err != nil {
		ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error creating team member"})
		return
	}
	ResponseJSON(w, http.StatusCreated, created)
}

func (m *TeamMemberHandler) get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	teamMember, err := m.teamMemberRepo.Get(id)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, model.ErrNoResults) {
			code = http.StatusNotFound
		}
		ResponseJSON(w, code, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error getting team member"})
		return
	}
	ResponseJSON(w, http.StatusOK, teamMember)
}

func (m *TeamMemberHandler) update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var teamMember model.TeamMember
	if err := json.NewDecoder(r.Body).Decode(&teamMember); err != nil {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request"})
		return
	}
	if teamMember.Id != id {
		ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "invalid request, teamMember id is different"})
		return
	}
	updated, err := m.teamMemberRepo.Update(&teamMember)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, model.ErrNoResults) {
			statusCode = http.StatusNotFound
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "error updating teamMember"})
		return
	}
	ResponseJSON(w, http.StatusOK, updated)
}

func (m *TeamMemberHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := m.teamMemberRepo.Delete(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, model.ErrNoResults) {
			statusCode = http.StatusNotFound
		}
		ResponseJSON(w, statusCode, model.CommonResponse{Status: model.CommonResponseStatusError, Message: "Error deleting team member"})
		return
	}
	ResponseJSON(w, http.StatusOK, model.CommonResponse{Status: model.CommonResponseStatusOk, Message: "Team member deleted"})
}
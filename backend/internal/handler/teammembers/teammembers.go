package teammembers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sanekee/limi"
	"github.com/sanekee/merchant-api/internal/handler"
	"github.com/sanekee/merchant-api/internal/model"
)

type TeamMemberRepo interface {
	GetByMerchantID(teamMemberID string, opts model.Pagination) ([]*model.TeamMember, error)
	Get(string) (*model.TeamMember, error)
	Insert(*model.NewTeamMember) (*model.TeamMember, error)
	Update(string, *model.UpdateTeamMember) (*model.TeamMember, error)
	Delete(string) error
}

type TeamMembers struct {
	teamMemberRepo TeamMemberRepo
}

func NewTeamMembers(repo TeamMemberRepo) TeamMembers {
	return TeamMembers{
		teamMemberRepo: repo,
	}
}

func (t TeamMembers) Post(w http.ResponseWriter, r *http.Request) {
	var newTeamMember model.NewTeamMember
	if err := json.NewDecoder(r.Body).Decode(&newTeamMember); err != nil {
		handler.ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.Error, Message: "invalid request"})
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

		handler.ResponseJSON(w, code, model.CommonResponse{Status: model.Error, Message: "Error creating team member"})
		return
	}

	handler.ResponseJSON(w, http.StatusCreated, created)
}

type TeamMembersByMerchant struct {
	limi           struct{} `path:"merchant/{merchant_id}"`
	teamMemberRepo TeamMemberRepo
}

func NewTeamMembersByMerchant(repo TeamMemberRepo) TeamMembersByMerchant {
	return TeamMembersByMerchant{
		teamMemberRepo: repo,
	}
}

func (t TeamMembersByMerchant) Get(w http.ResponseWriter, r *http.Request) {
	merchantID := limi.GetURLParam(r.Context(), "merchant_id")

	opt, err := handler.GetPaginationFromReq(r)
	if err != nil {
		handler.ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.Error, Message: err.Error()})
		return
	}

	teamMembers, err := t.teamMemberRepo.GetByMerchantID(merchantID, opt)
	if err != nil {
		handler.ResponseJSON(w, http.StatusInternalServerError, model.CommonResponse{Status: model.Error, Message: err.Error()})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, teamMembers)
}

type TeamMember struct {
	limi           struct{} `path:"{id}"`
	teamMemberRepo TeamMemberRepo
}

func NewTeamMember(repo TeamMemberRepo) TeamMember {
	return TeamMember{
		teamMemberRepo: repo,
	}
}

func (t TeamMember) Get(w http.ResponseWriter, r *http.Request) {
	id := limi.GetURLParam(r.Context(), "id")

	teamMember, err := t.teamMemberRepo.Get(id)
	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, model.ErrNoResults) {
			code = http.StatusNotFound
		}

		handler.ResponseJSON(w, code, model.CommonResponse{Status: model.Error, Message: "Error getting team member"})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, teamMember)
}

func (t *TeamMember) Put(w http.ResponseWriter, r *http.Request) {
	id := limi.GetURLParam(r.Context(), "id")

	var utm model.UpdateTeamMember
	if err := json.NewDecoder(r.Body).Decode(&utm); err != nil {
		handler.ResponseJSON(w, http.StatusBadRequest, model.CommonResponse{Status: model.Error, Message: "invalid request"})
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

		handler.ResponseJSON(w, statusCode, model.CommonResponse{Status: model.Error, Message: "error updating teamMember"})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, updated)
}

func (t *TeamMember) Delete(w http.ResponseWriter, r *http.Request) {
	id := limi.GetURLParam(r.Context(), "id")

	err := t.teamMemberRepo.Delete(id)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, model.ErrNoResults) {
			statusCode = http.StatusNotFound
		}

		handler.ResponseJSON(w, statusCode, model.CommonResponse{Status: model.Error, Message: "Error deleting team member"})
		return
	}

	handler.ResponseJSON(w, http.StatusOK, model.CommonResponse{Status: model.Ok, Message: "Team member deleted"})
}

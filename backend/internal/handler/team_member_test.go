package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/mock"
	"github.com/sanekee/merchant-api/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTeamMemberHandler_update(t *testing.T) {
	tests := []struct {
		name     string
		repo     TeamMemberRepo
		id       string
		body     *model.TeamMember
		validate validateFunc
	}{
		{
			name: "Success",
			repo: mock.NewTeamMemberRepo(nil, []*model.TeamMember{{Id: "ID1", Email: "e@t.est", MerchantId: "M1"}}),
			id:   "ID1",
			body: &model.TeamMember{Id: "ID1", Email: "updated@t.est", MerchantId: "M1"},
			validate: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				var tm *model.TeamMember
				if err := json.NewDecoder(resp.Body).Decode(&tm); err != nil {
					t.Fatal("error unmarshalling response")
				}
				assert.NotNil(t, tm)
				if tm != nil {
					assert.Equal(t, "ID1", tm.Id, "id is not the same")
					assert.Equal(t, "updated@t.est", tm.Email, "updated email is not the same")
				}
			},
		},
		{
			name: "Not Exists",
			repo: mock.NewTeamMemberRepo(model.ErrNoResults, nil),
			id:   "ID1",
			body: &model.TeamMember{Id: "ID1", Email: "updated@t.est", MerchantId: "M1"},
			validate: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
				var crsp *model.CommonResponse
				if err := json.NewDecoder(resp.Body).Decode(&crsp); err != nil {
					t.Fatal("error unmarshalling response")
				}
				assert.NotNil(t, resp, "expected CommonResponse")
			},
		},
		{
			name: "Repo Error",
			repo: mock.NewTeamMemberRepo(errors.New("repo error"), nil),
			id:   "ID1",
			body: &model.TeamMember{Id: "ID1", Email: "updated@t.est", MerchantId: "M1"},
			validate: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
				var crsp *model.CommonResponse
				if err := json.NewDecoder(resp.Body).Decode(&crsp); err != nil {
					t.Fatal("error unmarshalling response")
				}
				assert.NotNil(t, resp, "expected CommonResponse")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewTeamMemberHandler("test", tt.repo)
			teamMemberHTTPTest(t, "PUT", map[string]string{"id": tt.id}, h.update, tt.body, tt.validate)
		})
	}
}

func teamMemberHTTPTest(t *testing.T, method string, urlVars map[string]string, fn http.HandlerFunc, body *model.TeamMember, vfn validateFunc) {
	var bReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bReader = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, "http://test", bReader)
	if urlVars != nil {
		req = mux.SetURLVars(req, urlVars)
	}
	w := httptest.NewRecorder()
	fn(w, req)
	vfn(t, w.Result())
}

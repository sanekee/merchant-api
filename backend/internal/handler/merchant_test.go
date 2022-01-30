package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/mock"
	"github.com/sanekee/merchant-api/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

type validateFunc func(*testing.T, *http.Response)

func TestMerchantHandler_get(t *testing.T) {
	tests := []struct {
		name     string
		repo     MerchantRepo
		id       string
		validate validateFunc
	}{
		{
			name: "Exists",
			repo: mock.NewMerchantRepo(nil, []*model.Merchant{{Id: "ID1", Code: "Code1"}}),
			id:   "ID1",
			validate: func(t *testing.T, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				var mc *model.Merchant
				if err := json.NewDecoder(resp.Body).Decode(&mc); err != nil {
					t.Fatal("error unmarshalling response")
				}
				assert.NotNil(t, mc)
				if mc != nil {
					assert.Equal(t, "ID1", mc.Id, "id is not the same")
					assert.Equal(t, "Code1", mc.Code, "code is not the same")
				}
			},
		},
		{
			name: "Not Exists",
			repo: mock.NewMerchantRepo(model.ErrNoResults, nil),
			id:   "ID1",
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
			repo: mock.NewMerchantRepo(errors.New("repo error"), nil),
			id:   "ID1",
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
			m := NewMerchantHandler("test", tt.repo)
			merchantHTTPTest(t, "GET", map[string]string{"id": tt.id}, m.get, tt.validate)
		})
	}
}

func merchantHTTPTest(t *testing.T, method string, urlVars map[string]string, fn http.HandlerFunc, vfn validateFunc) {
	req := httptest.NewRequest(method, "http://test", nil)
	if urlVars != nil {
		req = mux.SetURLVars(req, urlVars)
	}
	w := httptest.NewRecorder()
	fn(w, req)
	vfn(t, w.Result())
}

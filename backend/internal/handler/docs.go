package handler

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

type DocsHandler struct {
	httpPath string
	yamlPath string
}

func NewDocsHandler(httpPath string, yamlPath string) *DocsHandler {
	return &DocsHandler{
		httpPath: httpPath,
		yamlPath: yamlPath,
	}
}
func (s *DocsHandler) RegisterMux(router *mux.Router) {
	r := router.PathPrefix(s.httpPath).Subrouter()
	r.Path("").Methods(http.MethodGet).HandlerFunc(s.get)

}

func (s *DocsHandler) get(w http.ResponseWriter, r *http.Request) {
	opts := middleware.SwaggerUIOpts{SpecURL: s.yamlPath, Path: s.httpPath}
	sh := middleware.SwaggerUI(opts, nil)
	sh.ServeHTTP(w, r)
	return
}

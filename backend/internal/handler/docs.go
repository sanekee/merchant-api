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
	r.Path("/rapi").Methods(http.MethodGet).HandlerFunc(s.rapi)
	r.Path("/redoc").Methods(http.MethodGet).HandlerFunc(s.redoc)

}

func (s *DocsHandler) get(w http.ResponseWriter, r *http.Request) {
	opts := middleware.SwaggerUIOpts{SpecURL: s.yamlPath, Path: s.httpPath}
	sh := middleware.SwaggerUI(opts, nil)
	sh.ServeHTTP(w, r)
	return
}

func (s *DocsHandler) rapi(w http.ResponseWriter, r *http.Request) {
	opts := middleware.RapiDocOpts{SpecURL: s.yamlPath, Path: s.httpPath + "/rapi"}
	sh := middleware.RapiDoc(opts, nil)
	sh.ServeHTTP(w, r)
	return
}

func (s *DocsHandler) redoc(w http.ResponseWriter, r *http.Request) {
	opts := middleware.RedocOpts{SpecURL: s.yamlPath, Path: s.httpPath + "/redoc"}
	sh := middleware.Redoc(opts, nil)
	sh.ServeHTTP(w, r)
	return
}

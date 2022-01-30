package handler

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

type DocsHandler struct {
	httpPath      string
	yamlPath      string
	swaggerUIOpts middleware.SwaggerUIOpts
	rapidodcOpts  middleware.RapiDocOpts
	redocOpts     middleware.RedocOpts
}

func NewDocsHandler(httpPath string, yamlPath string) *DocsHandler {
	swaggerUIOpts := middleware.SwaggerUIOpts{SpecURL: yamlPath, Path: httpPath}
	rapidocOpts := middleware.RapiDocOpts{SpecURL: yamlPath, Path: httpPath + "/rapi"}
	redocOpts := middleware.RedocOpts{SpecURL: yamlPath, Path: httpPath + "/redoc"}
	return &DocsHandler{
		httpPath:      httpPath,
		yamlPath:      yamlPath,
		swaggerUIOpts: swaggerUIOpts,
		rapidodcOpts:  rapidocOpts,
		redocOpts:     redocOpts,
	}
}
func (s *DocsHandler) RegisterMux(router *mux.Router) {
	r := router.PathPrefix(s.httpPath).Subrouter()
	r.Path("").Methods(http.MethodGet).HandlerFunc(s.get)
	r.Path("/rapi").Methods(http.MethodGet).HandlerFunc(s.rapi)
	r.Path("/redoc").Methods(http.MethodGet).HandlerFunc(s.redoc)

}

func (s *DocsHandler) get(w http.ResponseWriter, r *http.Request) {
	sh := middleware.SwaggerUI(s.swaggerUIOpts, nil)
	sh.ServeHTTP(w, r)
	return
}

func (s *DocsHandler) rapi(w http.ResponseWriter, r *http.Request) {
	sh := middleware.RapiDoc(s.rapidodcOpts, nil)
	sh.ServeHTTP(w, r)
	return
}

func (s *DocsHandler) redoc(w http.ResponseWriter, r *http.Request) {
	sh := middleware.Redoc(s.redocOpts, nil)
	sh.ServeHTTP(w, r)
	return
}

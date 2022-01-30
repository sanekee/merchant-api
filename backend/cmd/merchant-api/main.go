package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/handler"
	"github.com/sanekee/merchant-api/backend/internal/log"
	go_up "github.com/ufoscout/go-up"
)

type MuxHandler interface {
	RegisterMux(r *mux.Router)
}

var (
	port    int
	docPath string
)

func init() {
	up, err := go_up.NewGoUp().
		AddReader(go_up.NewEnvReader("", false, false)). // Loads environment variables
		Build()
	if err != nil {
		panic("Error initializing environment reader")
	}
	port = up.GetIntOrDefault("APP_PORT", 8123)
	docPath = up.GetStringOrDefault("APP_SPEC", "./spec")
}

func main() {

	handlers := []MuxHandler{
		handler.NewDocsHandler("/docs", "/openapi.yaml"),
	}

	r := mux.NewRouter()
	for _, h := range handlers {
		h.RegisterMux(r)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(docPath)))

	listenPort := fmt.Sprintf(":%d", port)
	log.Info("Start listening %s", listenPort)
	http.ListenAndServe(listenPort, r)
}
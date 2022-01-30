package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanekee/merchant-api/backend/internal/db"
	"github.com/sanekee/merchant-api/backend/internal/handler"
	"github.com/sanekee/merchant-api/backend/internal/log"
	"github.com/sanekee/merchant-api/backend/internal/mock"
	"github.com/sanekee/merchant-api/backend/internal/repo"
	go_up "github.com/ufoscout/go-up"
)

type MuxHandler interface {
	RegisterMux(r *mux.Router)
}

var (
	port    int
	docPath string
	mockDB  bool
	pgUser  string
	pgPass  string
	pgHost  string
	pgDB    string
	pgPort  int
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
	mockDB = up.GetBoolOrDefault("APP_MOCKDB", true)

	pgUser = up.GetStringOrDefault("POSTGRES_USER", "mapi")
	pgPass = up.GetStringOrDefault("POSTGRES_PASSWORD", "mypostgres")
	pgHost = up.GetStringOrDefault("POSTGRES_HOST", "127.0.0.1")
	pgPort = up.GetIntOrDefault("POSTGRES_PORT", 5432)
	pgDB = up.GetStringOrDefault("POSTGRES_DB", "mapi")
}

func main() {

	var merchantRepo handler.MerchantRepo
	var teamMemberRepo handler.TeamMemberRepo

	if mockDB {
		merchantRepo = mock.NewMerchantRepo(nil, mock.GenerateMerchants(1000))
		teamMemberRepo = mock.NewTeamMemberRepo(nil, mock.GenerateTeamMembers(1000, "test-0"))
	} else {
		db := db.NewPGDB(pgHost, pgDB, pgPort, pgUser, pgPass)
		merchantRepo = repo.NewMerchantRepo(db)
		teamMemberRepo = mock.NewTeamMemberRepo(nil, mock.GenerateTeamMembers(1000, "test-0"))
		// teamMemberRepo = repo.NewTeamMemberRepo(db)
	}
	handlers := []MuxHandler{
		handler.NewDocsHandler("/docs", "/openapi.yaml"),
		handler.NewMerchantHandler("/merchant", merchantRepo),
		handler.NewTeamMemberHandler("/teammember", teamMemberRepo),
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

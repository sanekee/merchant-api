package main

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sanekee/limi"
	limimw "github.com/sanekee/limi/middleware"
	"github.com/sanekee/merchant-api/internal/handler/merchants"
	"github.com/sanekee/merchant-api/internal/handler/teammembers"
	"github.com/sanekee/merchant-api/internal/log"
	mw "github.com/sanekee/merchant-api/internal/middleware"
	"github.com/sanekee/merchant-api/internal/repo"
	go_up "github.com/ufoscout/go-up"
)

var (
	port    int
	docPath string
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

	pgUser = up.GetStringOrDefault("POSTGRES_USER", "mapi")
	pgPass = up.GetStringOrDefault("POSTGRES_PASSWORD", "mypostgres")
	pgHost = up.GetStringOrDefault("POSTGRES_HOST", "127.0.0.1")
	pgPort = up.GetIntOrDefault("POSTGRES_PORT", 5432)
	pgDB = up.GetStringOrDefault("POSTGRES_DB", "mapi")
}

func main() {

	db := repo.NewPGDB(pgHost, pgDB, pgPort, pgUser, pgPass)
	merchantRepo := repo.NewMerchantRepo(db)
	teamMemberRepo := repo.NewTeamMemberRepo(db)

	// api router with JWT auth
	r1, err := limi.NewRouter(
		"/api",
		limi.WithMiddlewares(
			limimw.Log(log.Logger{}),
			mw.CORS,
			mw.JWTAuth,
		),
	)
	if err != nil {
		panic(err)
	}

	if err := r1.AddHandlers([]limi.Handler{
		merchants.NewMerchants(merchantRepo),
		merchants.NewMerchant(merchantRepo),
		teammembers.NewTeamMembers(teamMemberRepo),
		teammembers.NewTeamMembersByMerchant(teamMemberRepo),
		teammembers.NewTeamMember(teamMemberRepo),
	}); err != nil {
		panic(err)
	}

	// a general router
	r2, err := limi.NewRouter(
		"/",
		limi.WithMiddlewares(
			limimw.Log(log.Logger{}),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := r2.AddHTTPHandler("/spec/", http.FileServer(http.Dir(docPath))); err != nil {
		panic(err)
	}

	rapiddoc := middleware.RapiDoc(middleware.RapiDocOpts{SpecURL: "/spec/openapi.yaml", Path: "/docs"}, nil)
	if err := r2.AddHTTPHandler("/docs", rapiddoc); err != nil {
		panic(err)
	}

	listenPort := fmt.Sprintf(":%d", port)
	log.Info("Start listening %s", listenPort)

	mux := limi.NewMux(r1, r2)
	http.ListenAndServe(listenPort, mux)
}

package db

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/backend/internal/log"
)

func NewPGDB(pgHost string, pgDB string, pgPort int, pgUser string, pgPass string) *pg.DB {
	opt := &pg.Options{
		User:     pgUser,
		Password: pgPass,
		Database: pgDB,
		Addr:     fmt.Sprintf("%s:%d", pgHost, pgPort),
	}
	log.Info("connecting to pgdb %s:%d %s %s",
		pgHost, pgPort, pgDB, pgUser)
	pg := pg.Connect(opt)
	return pg

}

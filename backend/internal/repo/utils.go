package repo

import (
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/backend/internal/log"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

func toAppError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pg.ErrNoRows) {
		return model.ErrNoResults
	}
	log.Error("Repo Error %s", err.Error())
	return model.ErrServer
}

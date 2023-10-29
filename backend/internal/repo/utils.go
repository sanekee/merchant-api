package repo

import (
	"errors"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/sanekee/merchant-api/internal/log"
	"github.com/sanekee/merchant-api/internal/model"
)

func toAppError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pg.ErrNoRows) {
		return model.ErrNoResults
	}

	errStr := err.Error()
	log.Error("Repo Error %s", errStr)

	switch true {
	case strings.Contains(errStr, "#23502"), strings.Contains(errStr, "#23503"):
		return model.ErrRequest
	case strings.Contains(errStr, "#23505"):
		return model.ErrDuplicate
	}

	return model.ErrServer
}

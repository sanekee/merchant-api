package handler

import (
	"net/http"
	"strconv"

	"github.com/sanekee/merchant-api/backend/internal/log"
	"github.com/sanekee/merchant-api/backend/internal/model"
)

func getPaginationFromReq(r *http.Request) (model.Pagination, error) {
	opt := model.Pagination{}
	params := r.URL.Query()
	limitStr := params.Get("limit")
	offsetStr := params.Get("offset")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return opt, model.ErrRequest
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		return opt, model.ErrRequest
	}
	opt.Limit = int(limit)
	opt.Offset = int(offset)
	log.Debug("Option: limit %d, offset %d", opt.Limit, opt.Offset)
	return opt, nil
}

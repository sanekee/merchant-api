package handler

import (
	"net/http"
	"strconv"

	"github.com/sanekee/merchant-api/backend/internal/model"
)

func getPaginationFromReq(r *http.Request) (model.Pagination, error) {
	opt := model.Pagination{Limit: 10, Offset: 0}
	params := r.URL.Query()
	limitStr := params.Get("limit")
	offsetStr := params.Get("offset")
	if len(limitStr) > 0 {
		limit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			return opt, model.ErrRequest
		}
		opt.Limit = int(limit)
	}
	if len(offsetStr) > 0 {
		offset, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			return opt, model.ErrRequest
		}
		opt.Offset = int(offset)
	}
	return opt, nil
}

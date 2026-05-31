package handlers

import (
	"net/http"
	"strconv"
)

func parsePage(r *http.Request, defaultLimit int) (page, limit, offset int) {
	page = 1
	limit = defaultLimit
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	offset = (page - 1) * limit
	return
}

func pageResponse(data any, total, page, limit int) map[string]any {
	return map[string]any{
		"data":  data,
		"total": total,
		"page":  page,
		"limit": limit,
	}
}

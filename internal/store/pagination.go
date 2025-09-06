package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Offset int    `json:"offset" validate:"min=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {

	qs := r.URL.Query()
	limit, offset, sort := qs.Get("limit"), qs.Get("offset"), qs.Get("sort")

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil

		}
		fq.Limit = l
	}

	if offset != "" {
		f, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = f

	}

	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}

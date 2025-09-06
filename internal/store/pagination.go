package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"min=1,max=100"`
	Offset int      `json:"offset" validate:"min=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// pagination
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil

		}
		fq.Limit = l
	}

	if offset := qs.Get("offset"); offset != "" {
		f, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = f

	}

	// sort
	if sort := qs.Get("sort"); sort != "" {
		fq.Sort = sort
	}

	// filters
	if tags := qs.Get("tags"); tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	if search := qs.Get("search"); search != "" {
		fq.Search = search
	}

	if since := qs.Get("since"); since != "" {
		fq.Since = parseTime(since)
	}

	if until := qs.Get("until"); until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(value string) string {
	t, err := time.Parse(time.DateTime, value)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}

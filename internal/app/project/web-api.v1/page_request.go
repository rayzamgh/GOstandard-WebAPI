package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/standard-go/project/internal/app/project"
)

var (
	pageRequestCtxKey = &ctxKey{"PageRequest"}
)

func PageRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pageRequest := &project.PageRequest{
			Filters:  make([]project.Filter, 0),
			Sorts:    make([]project.Sort, 0),
			Page:     1,
			PerPage:  10,
			Paginate: 0,
			Search:   "",
		}
		query := r.URL.Query()

		if query.Get("page") != "" {
			pageRequest.Page, _ = strconv.ParseInt(query.Get("page"), 0, 64)
		}

		if query.Get("per_page") != "" {
			pageRequest.PerPage, _ = strconv.ParseInt(query.Get("per_page"), 0, 64)
		}

		if query.Get("paginate") != "" {
			pageRequest.Paginate = 1
		}

		if query.Get("search") != "" {
			pageRequest.Search = query.Get("search")
		}

		if query.Get("sort") != "" {
			srts := strings.Split(query.Get("sort"), "|")
			if len(srts) >= 2 {
				pageRequest.Sorts = append(pageRequest.Sorts, project.Sort{Option: srts[0], Value: srts[1]})
			}
		}

		//Fetching filters data
		if len(query["filter[]"]) > 0 {
			for _, v := range query["filter[]"] {
				filter := project.Filter{}
				json.Unmarshal([]byte(v), &filter)
				pageRequest.Filters = append(pageRequest.Filters, filter)
			}
		}

		//Fetching sorts data
		if len(query["sorts[]"]) > 0 {
			for _, v := range query["sorts[]"] {
				sort := project.Sort{}
				json.Unmarshal([]byte(v), &sort)
				pageRequest.Sorts = append(pageRequest.Sorts, sort)
			}
		}

		ctx := context.WithValue(r.Context(), pageRequestCtxKey, pageRequest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

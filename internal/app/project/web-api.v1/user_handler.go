package api

import (
	"bytes"
	"encoding/json"
	// "errors"
	"math"
	"net/http"
	// "fmt"

	"github.com/go-chi/chi"

	"gitlab.com/standard-go/project/internal/app/project"
	"gitlab.com/standard-go/project/internal/app/config/env"
	"gitlab.com/standard-go/project/internal/app/responses"
)

func IndexUser(w http.ResponseWriter, r *http.Request) {
	pageRequest := r.Context().Value(pageRequestCtxKey).(*project.PageRequest)

	fetch, count, err := srv.FetchIndexUser(pageRequest)
	if err != nil {
		printError(err, w)
		return
	}

	totalPagesInt64 := int64(math.Ceil(float64(count) / float64(pageRequest.PerPage)))
	buffer := new(bytes.Buffer)
	isFirstParam := true

	for k, v := range r.URL.Query() {
		if k != "page" {
			if !isFirstParam {
				buffer.WriteString("&")
			} else {
				isFirstParam = false
			}
			buffer.WriteString(k + "=" + v[0])
		}
	}

	path := env.Get("APP_HOST") + "/api/v1/claims?" + buffer.String()
	nextPageUrl, prevPageUrl := "", ""

	if pageRequest.Page >= 1 && pageRequest.Page <= totalPagesInt64 {
		nextPageUrl = checkPage(pageRequest.Page, totalPagesInt64, 1, isFirstParam, path)
		prevPageUrl = checkPage(pageRequest.Page, 1, -1, isFirstParam, path)
	}

	res := setResponse(pageRequest, fetch, totalPagesInt64, count, path, nextPageUrl, prevPageUrl)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")

	fetch, err := srv.FetchShowUser(userId)
	if err != nil {
		printError(err, w)
		return
	}

	res := responses.NewResponse(fetch)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func StoreUser(w http.ResponseWriter, r *http.Request) {
	// auth := r.Context().Value(authUserCtxKey).(*claim.Auth)
	var userReq *project.User
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		printError(err, w)
		return
	}

	fetch, err := srv.FetchStoreUser(userReq)
	if err != nil {
		printError(err, w)
		return
	}

	res := responses.NewResponse(fetch)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userReq *project.User

	userId := chi.URLParam(r, "user_id")

	fetch, err := srv.FetchShowUser(userId)
	if err != nil {
		printError(err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		printError(err, w)
		return
	}

	fetch.FullName = userReq.FullName

	fetch, err = srv.FetchUpdateUser(userId, fetch)
	if err != nil {
		printError(err, w)
		return
	}

	res := responses.NewResponse(fetch)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func DestroyUser(w http.ResponseWriter, r *http.Request) {
	var userReq *project.User

	userId := chi.URLParam(r, "user_id")

	_, err := srv.FetchShowUser(userId)
	if err != nil {
		printError(err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		printError(err, w)
		return
	}

	err = srv.FetchDestroyUser(userId)
	if err != nil {
		printError(err, w)
		return
	}

	res := responses.NewResponse(nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}


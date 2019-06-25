package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/standard-go/project/internal/app/project"
	"gitlab.com/standard-go/project/internal/app/responses"
)

func IndexUser(w http.ResponseWriter, r *http.Request) {
	pageRequest := r.Context().Value(pageRequestCtxKey).(*project.PageRequest)

	fetch, count, err := srv.FetchIndexUser(pageRequest)
	if err != nil {
		printError(err, w)
		return
	}

	res := setPaginate(r, pageRequest, fetch, count, project.UserUrl)

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


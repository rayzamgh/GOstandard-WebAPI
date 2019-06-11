package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt"

	"gitlab.com/standard-go/project/internal/app/config/env"
	"gitlab.com/standard-go/project/internal/app/project"
)

//
// Contexts
//

type ctxKey struct {
	name string
}

func (k ctxKey) String() string {
	return fmt.Sprintf("tugure-location's context value: %s", k.name)
}

var (
	statusCtxKey      = &ctxKey{"Status"}
	authUserCtxKey    = &ctxKey{"AuthUser"}
)

/*
====================================
    MIDDLEWARE
====================================
*/

// CommonMiddleware tak jelas ini, dia memaksakan Content-Type dari request dengan application/json. Faedahnya adalah?
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// JWTMiddleware buat periksa jwt token yang dibuat dari service otentikasi.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer ")
		if len(token) <= 0 {
			printError(errors.New("401"), w)
			return
		}
		token = splitToken[1]
		hs256 := jwt.NewHS256(env.Get("JWT_SECRET"))

		payload, sig, err := jwt.Parse(token)
		if err != nil {
			printError(errors.New("400-Token-Parse"), w)
			return
		}
		if err = hs256.Verify(payload, sig); err != nil {
			printError(errors.New("400-Signature"), w)
			return
		}
		var jot project.Token
		if err = jwt.Unmarshal(payload, &jot); err != nil {
			printError(errors.New("400-Token-Parse"), w)
			return
		}

		auth := &project.Auth{
			User: 	jot.ToUser(),
			Token: 	token,
		}

		// Validate fields.
		iatValidator := jwt.IssuedAtValidator(now)
		expValidator := jwt.ExpirationTimeValidator(now)
		audValidator := jwt.AudienceValidator("admin")
		if err = jot.Validate(iatValidator, expValidator, audValidator); err != nil {
			switch err {
			case jwt.ErrIatValidation:
				printError(errors.New("400-Signature"), w)
				return
			case jwt.ErrExpValidation:
				printError(errors.New("401-Expired"), w)
				return
			}
		}

		ctx := context.WithValue(r.Context(), authUserCtxKey, auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package api

import (
	"github.com/go-chi/chi"
)

// InitAPIRoutes sets project API routes
func InitAPIRoutes(r chi.Router) {
	/*
	|--------------------------------------------------------------------------
	| Users Routes
	|--------------------------------------------------------------------------
	|
	*/
	r.Route("/users", func(r chi.Router) {
		r.With(PageRequestCtx).Get("/", IndexUser)
		r.Post("/", StoreUser)
		r.Route("/{user_id}", func(r chi.Router) {
			r.Get("/", ShowUser)
			r.Put("/", UpdateUser)
			r.Delete("/", DestroyUser)
		})
	})
}
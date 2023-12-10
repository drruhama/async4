package auth

import (
	"database/sql"
	routerChi "ecommerce/infra/router/chi"

	"github.com/go-chi/chi/v5"
)

// function ini untuk melakukan init terhadap semua
// hal yang dibutuhkan oleh Auth Services
func Register(router *chi.Mux, db *sql.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// seperti grouping endpoint
	// sebagai dasarnya
	router.Route("/ecommerce/auth", func(r chi.Router) {
		r.Post("/signup", handler.Register)
		r.Post("/signin", handler.Login)
		// method ini berfungsi untuk
		// mengelompokkan router yang behavior nya sama
		r.Group(func(r chi.Router) {
			// use middleware selalu di awal
			r.Use(routerChi.CheckToken)
			// middlewarenya nanti akan digunakan untuk endpoint ini
			r.Get("/profile", handler.Profile)
		})
	})
}

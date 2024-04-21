package handler

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/orgs/DonSTUhackathon/beatifullBack/service_register/pkg/dbworker"
	"github.com/orgs/DonSTUhackathon/beatifullBack/service_register/pkg/service"
)

type Handler struct {
	//service - running an application
	Services *service.Service
	//server - config the http server
	Server *service.Server
	//dbinstance - operations with database
	DBInstance *sqllogic.Database
	//sessionsstore - store and work with cookies
	SessionsStore *sessions.CookieStore
}

// Constructor of a handler
func NewHandler(services *service.Service, server *service.Server, DB *sql.DB) *Handler {
	return &Handler{Services: services, Server: server, DBInstance: &sqllogic.Database{Db: DB}, SessionsStore: NewSessionStorage()}
}

// API-Handler itself
func (h *Handler) InitRoutes() *chi.Mux {

	/////////////////////////////////////////////////////////////////////////////////////////////
	//init new router
	r := chi.NewRouter()
	// redirect /auth/ to /auth
	r.Use(middleware.RedirectSlashes)

	//serve all the api-routes

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", h.LoginUser)
		r.Post("/registration", h.RegisterUser)
		r.Put("/logout", h.LogoutUser)
		r.Post("/setdesc", h.SetDescription)
	})

	/////////////////////////////////////////////////////////////////////////////////////////////

	return r
}

package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/orgs/DonSTUhackathon/beatifullBack/service_chat/handler"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)

	/////////////////////////////////////////////////////////////////////////////////////////////
	r.Route("/chat/{chat_id}/", func(r chi.Router) {
		r.Get("/", handler.AuthMiddleware(http.HandlerFunc(handler.SocketHandler)))
		// r.Post("/sendmessage", handler.AuthMiddleware())
		// r.Post("/deletemessage", handler.AuthMiddleware())
		// r.Post("/editmessage", handler.AuthMiddleware())
		// r.Post("/closesocket", handler.AuthMiddleware())
	})

	http.HandleFunc("/start", handler.SocketHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}

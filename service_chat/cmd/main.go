package main

import (
	"database/sql"
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/orgs/DonSTUhackathon/beatifullBack/service_chat/handler"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RedirectSlashes)
	addr := flag.String("addr", ":8080", "http service address")
	var h handler.Adapter
	var err error
	h.Db, err = sql.Open("postgres", "user=admin password=4321 dbname=oggetto_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer h.Db.Close()
	/////////////////////////////////////////////////////////////////////////////////////////////

	if err != nil {
		log.Fatal(err)
	}
	defer h.Db.Close()

	r.Route("/api", func(r chi.Router) {
		r.Get("/users", h.GetUsers)
		r.Get("/chats", h.GetChats)
		r.Route("/messages", func(r chi.Router) {
			r.Get("/", h.GetMessages)
			r.Post("/send", h.CreateMessage)
			r.Post("/delete", h.DeleteMessage)
			r.Post("/edit", h.UpdateMessage)
		})
	})

	go h.HandleMessages()

	err = http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

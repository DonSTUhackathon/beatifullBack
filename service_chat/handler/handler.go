package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Adapter struct {
	Db        *sql.DB
	Clients   map[*websocket.Conn]bool
	Broadcast chan Message
	UPG       websocket.Upgrader
}

func (a Adapter) NewAdapter() *Adapter {
	a.UPG = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	a.Clients = make(map[*websocket.Conn]bool)
	a.Broadcast = make(chan Message)
	return &a
}

type User struct {
	ID        int    `json:"id"`
	GoogleID  string `json:"google_id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Position  string `json:"position"`
}

type Chat struct {
	ID1 int `json:"id1"`
	ID2 int `json:"id2"`
}

type Message struct {
	ID       int    `json:"id"`
	ChatID   int    `json:"chat_id"`
	UserID   int    `json:"user_id"`
	Content  string `json:"content"`
	IsHidden bool   `json:"is_hidden"`
}

func (a Adapter) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []User

	rows, err := a.Db.Query("SELECT id, google_id, name, birthdate, position FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.GoogleID, &user.Name, &user.Birthdate, &user.Position)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func (a Adapter) GetChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var chats []Chat

	rows, err := a.Db.Query("SELECT id1, id2 FROM chats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.ID1, &chat.ID2)
		if err != nil {
			log.Fatal(err)
		}
		chats = append(chats, chat)
	}

	json.NewEncoder(w).Encode(chats)
}

func (a Adapter) GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var messages []Message

	rows, err := a.Db.Query("SELECT id, chat_id, user_id, content, ishidden FROM messages")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.IsHidden)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, message)
	}

	json.NewEncoder(w).Encode(messages)
}

func (a Adapter) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO messages (chat_id, user_id, content, ishidden) VALUES ($1, $2, $3, $4) RETURNING id`
	err = a.Db.QueryRow(sqlStatement, message.ChatID, message.UserID, message.Content, message.IsHidden).Scan(&message.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)

	a.Broadcast <- message
}

func (a Adapter) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE messages SET chat_id=$1, user_id=$2, content=$3, ishidden=$4 WHERE id=$5 RETURNING id`
	err = a.Db.QueryRow(sqlStatement, message.ChatID, message.UserID, message.Content, message.IsHidden, params["id"]).Scan(&message.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (a Adapter) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	sqlStatement := `DELETE FROM messages WHERE id = $1`
	res, err := a.Db.Exec(sqlStatement, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a Adapter) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := a.UPG.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	a.Clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(a.Clients, ws)
			break
		}

		a.Broadcast <- msg
	}
}

func (a Adapter) HandleMessages() {
	for {
		msg := <-a.Broadcast
		for client := range a.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error: %v", err)
				client.Close()
				delete(a.Clients, client)
			}
		}
	}
}

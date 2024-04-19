package handler

import (
	"bufio"
	"github.com/gorilla/websocket"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{}
var connections = []*websocket.Conn{}

// var oauth2Config = &oauth2.Config{
// 	ClientID:     "",
// 	ClientSecret: "",
// 	RedirectURL:  "",
// 	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// 	Endpoint:     google.Endpoint,
// }

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if r.token == expected token
		next.ServeHTTP(w, r)
	})
}

func removeConn(slice []*websocket.Conn, val *websocket.Conn) []*websocket.Conn {
	index := -1
	for i, v := range slice {
		if v == val {
			index = i
			break
		}
	}
	if index != -1 {
		if index < len(slice)-1 {
			copy(slice[index:], slice[index+1:])
		}
		slice = slice[:len(slice)-1]
	}
	return slice
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	connections = append(connections, conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			connections = removeConn(connections, conn)
			log.Println("Error during message reading:", err)
			break
		}
		//work with msg
		log.Printf("Server: %s", message)
		for _, c := range connections {
			if c != conn {
				err = c.WriteMessage(messageType, message[:len(message)-1])
				if err != nil {
					connections = removeConn(connections, conn)
					log.Println("Error during message writing:", err)
					break
				}
			}
		}
	}
	connections = removeConn(connections, conn)
}

func userInput(inputChan chan string) {
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	inputChan <- msg
}

func chatInputHandler(conn *websocket.Conn, done chan interface{}) {
	for {
		inputChan := make(chan string)
		go userInput(inputChan)
		select {
		case input := <-inputChan:
			if strings.TrimSpace(input) == "" {
				continue
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Println("Error during writing to websocket:", err)
				return
			}

			select {
			case <-done:
				log.Println("Receiver Channel Closed")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel")
			}
			return
		}
	}
}
func SendMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error during message reading:", err)
		return
	}

	for _, c := range connections {
		err = c.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			return
		}
	}
}

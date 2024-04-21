package handler

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

// Function to initialize the session.CookieStore
func NewSessionStorage() *sessions.CookieStore {
	s := sessions.NewCookieStore([]byte(generateKey()))
	//read time of a session from config
	s.MaxAge(viper.GetInt("time"))
	return s
}

// Function to generate the session key
func generateKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

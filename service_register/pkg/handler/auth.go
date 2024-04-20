package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/orgs/DonSTUhackathon/beatifullBack/service_register/pkg/dbworker"
	"log"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ////////////////////////////////////////////////////////////////////////////
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var req sqllogic.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//log message when new user is being created
	log.Print("New user: " + req.Username + " " + hashPwd(req.Password))
	//set the header, so that the client will know how to deal with the response
	w.Header().Set("Content-Type", "application/json")
	//inserting the user into the database (error if something goes wrong)
	err = h.DBInstance.AddUser(req)
	if err != nil {
		log.Print("Not written!", err)
	}
	//start a session
	session, _ := h.SessionsStore.Get(r, "auth-session")
	//field of a session which represents the fact of users authentification
	session.Values["authenticated"] = true
	session.Values["password"] = hashPwd(req.Password)
	session.Values["username"] = req.Username
	//save and send the cookie
	session.Save(r, w)
	json.NewEncoder(w).Encode(map[string]string{"redirect": fmt.Sprintf("/" + req.Username)})
}

// ////////////////////////////////////////////////////////////////////////////
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	//new variable to store login and password
	var creds Credentials
	//decode the body of the request into the variable (error if wrong json structure)
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//selecting the user with such login from the DB
	row := h.DBInstance.Db.QueryRow("SELECT password FROM auth_user WHERE username = ?", creds.Username)
	var storedHashedPwd string
	//Scan DB response
	err = row.Scan(&storedHashedPwd)
	if err != nil {
		//If there is not such a login in the database, err value will be sql.ErrNoRows
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username", http.StatusUnauthorized)
		} else {
			//typical error handling
			log.Printf("Unexpected error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	//if passwords mismatch, send this error, else create session
	if storedHashedPwd != hashPwd(creds.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	//create session
	session, _ := h.SessionsStore.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["username"] = creds.Username
	session.Values["password"] = hashPwd(creds.Password)
	session.Save(r, w)
	//log message when the user logged in
	log.Print("New session: " + creds.Username + " for " + fmt.Sprint(session.Options.MaxAge) + " seconds")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"redirect": fmt.Sprintf("/" + creds.Username)})

}

//////////////////////////////////////////////////////////////////////////////

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionsStore.Get(r, "auth-session")
	session.Values["authenticated"] = 0
	session.Options.MaxAge = -1

	log.Print("Logging out")
	session.Save(r, w)

}

//////////////////////////////////////////////////////////////////////////////

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionsStore.Get(r, "auth-session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// small function to hash passwords
func hashPwd(a string) string {
	hasher := sha256.New()
	hasher.Write([]byte(a))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

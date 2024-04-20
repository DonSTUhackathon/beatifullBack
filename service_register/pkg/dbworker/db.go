package sqllogic

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	Db *sql.DB
}

type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewDBInstance() *sql.DB {
	r, _ := os.ReadFile("../configs/sql_config.txt")
	sqlDB, err := sql.Open("postgres", string(r))
	if err != nil {
		log.Print("1", err)
	}
	log.Print("Successfully connxted to DB")
	return sqlDB

}

type Profile struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}
type User struct {
	ID          int       `json:"id"`
	Password    string    `json:"password"`
	LastLogin   time.Time `json:"last_login"`
	IsSuperuser bool      `json:"is_superuser"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	IsStaff     bool      `json:"is_staff"`
	IsActive    bool      `json:"is_active"`
	DateJoined  time.Time `json:"date_joined"`
}

func DBConnect() (*Database, error) {
	token, _ := os.ReadFile("config/dbconfig.txt")
	db, err := sql.Open("mysql", string(token))
	if err != nil {
		return nil, fmt.Errorf("error occured while opening the DB: %s", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging DB: %s", err)
	}
	return &Database{db}, nil
}

func (d Database) AddUser(rd RegisterRequest) error {
	hasher := sha256.New()
	var user User

	hasher.Write([]byte(user.Password))
	user.Password = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	log.Print(string(user.Password))
	tx, err := d.Db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %s", err)
	}
	defer tx.Rollback()

	log.Print(user)
	err = tx.QueryRow(`
	INSERT INTO auth_user (password, last_login,is_superuser, username, first_name, last_name, email, is_staff, is_active,date_joined)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
`, rd.Password, time.Now(), false, rd.Username, rd.FirstName, rd.LastName, rd.Email, true, true, time.Now()).Scan(&user.ID)

	log.Print((user))
	if err != nil {
		return fmt.Errorf("error inserting user data: %s", err)
	}

	profile := Profile{
		Description: "Default user profile",
		ID:          user.ID,
	}
	log.Print(profile)

	_, err = tx.Exec(`
	INSERT INTO restauth_profile (id_id, description)
	VALUES ($1, $2)
`, user.ID, profile.Description)
	if err != nil {
		return fmt.Errorf("error inserting profile data: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %s", err)
	}

	return nil
}

func (d *Database) GetByID(ID int64) (*User, error) {
	var user User
	err := d.Db.QueryRow(`
		SELECT id, password, last_login, is_superuser, username, first_name, last_name, email, is_staff, is_active, date_joined
		FROM users
		WHERE id = ?
	`, ID).Scan(&user.ID, &user.Password, &user.LastLogin, &user.IsSuperuser, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.IsStaff, &user.IsActive, &user.DateJoined)
	if err != nil {
		return nil, fmt.Errorf("error selecting data: %s", err)
	}
	return &user, nil
}

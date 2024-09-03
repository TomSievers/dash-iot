package auth

import (
	"dash-iot/dash-iot/user"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var sessions map[string]Session

type Session struct {
	Permissions string
	expires     time.Time
}

func (s Session) Expired() bool {
	return s.expires.Before(time.Now())
}

type SessionConfig struct {
	LoginTimeout time.Duration
}

type AuthHandlerFunc func(http.ResponseWriter, *http.Request, Session)

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func AuthHandler(handler AuthHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")

		if err != nil {
			log.Printf("Session token not found: '%s'\n", err)
		} else if c != nil {
			token := c.Value
			session, exists := sessions[token]

			if exists {
				if !session.Expired() {
					handler(w, r, session)
					return
				} else {
					delete(sessions, token)
					log.Printf("Session '%s' expired \n", token)
				}
			} else {
				log.Printf("Session '%s' not found \n", token)
			}
		} else {
			log.Println("Cookie is nil")
		}

		w.Header().Add("Location", "/login")
		w.WriteHeader(http.StatusSeeOther)
	})
}

func Login(db *sql.DB, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var login loginRequest

			err := json.NewDecoder(r.Body).Decode(&login)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, "invalid json")
				return
			}

			tx, err := db.Begin()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, err.Error())
				return
			}

			res, err := user.GetUserByName(tx, login.Name)

			if err != nil {

			}

			err = bcrypt.CompareHashAndPassword(res.Password, []byte(login.Password))

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, err.Error())
				return
			}

		} else {
			handler(w, r)
		}
	})
}

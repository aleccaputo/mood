package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"mood/db"
	"mood/models"
	"net/http"
)

type Server struct {
	listenAddress string
	db            db.DB
}

func NewServer(listenAddress string, db db.DB) *Server {
	return &Server{listenAddress: listenAddress, db: db}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		loginRequest := new(models.UserLoginRequest)
		if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := s.db.LoginUser(loginRequest.Email, loginRequest.Password)
		if err != nil {
			http.Error(w, errors.New("unauthorized").Error(), http.StatusUnauthorized)
			return
		}

		token, err := createJwt(id.String())
		if err != nil {
			http.Error(w, errors.New("unauthorized").Error(), http.StatusUnauthorized)
			return
		}
		WriteJSON(w, http.StatusOK, token)
	})

	mux.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "return all comments")
	})

	mux.HandleFunc("GET /user/{id}", withJWTAuth(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		user, err := s.db.Get(parsedUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = WriteJSON(w, http.StatusOK, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "The users name is %s", user.FirstName)
	}))

	mux.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) {
		createUserRequest := new(models.CreateUserRequest)
		if err := json.NewDecoder(r.Body).Decode(createUserRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := s.db.CreateUser(models.NewUser(createUserRequest.FirstName, createUserRequest.LastName, createUserRequest.Email), createUserRequest.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenString, err := createJwt("foo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		w.Header().Set("Authorization", "Bearer "+tokenString)
		WriteJSON(w, http.StatusAccepted, userId)
		return
	})

	if httpError := http.ListenAndServe("localhost:8080", mux); httpError != nil {
		fmt.Println(httpError.Error())
	}

	return http.ListenAndServe(s.listenAddress, mux)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

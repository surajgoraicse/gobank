package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}

}
func (s *APIServer) Run() {
	router := mux.NewRouter()

	// the second paramtere accepts a function of tpe : func(http.ResponseWriter, *http.Request) but our method looks like : func(http.ResponseWriter, *http.Request) so we will create a wrapper fuction that will return the desired type of function
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))

	log.Println("JSON API server running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	

	return fmt.Errorf("method not allowed : %s", r.Method)
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	method := r.Method
	fmt.Println("here" , method)
	if method == "DELETE" {
		fmt.Println("here i am ")
		return s.handleDeleteAccount(w, r)
	}

	if method == "GET" {
		return s.handleGetAccountByID(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET : Get one account by id
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("invalid id given %s", idstr)
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return WriteJSON(w, http.StatusOK, "user not found")
		} else {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}

	return WriteJSON(w, http.StatusOK, account)
}

// GET : Get all accounts
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, accounts)
	return nil

}

// POST : create account
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	// createAccountReq := CreateAccountRequest{} //  pass pointer
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	// account response is not coming from databse
	return WriteJSON(w, http.StatusCreated, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error converting param string id  to int")
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		fmt.Println("error deleting")
		return err
	}
	WriteJSON(w, http.StatusOK, "account deleted successfully")
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// utility for sending response
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v) // what if this line return a error .
}

// utiltiy for sending error
type ApiError struct {
	Error string `json:"error"`
}

// handler function type
type apiFunc func(http.ResponseWriter, *http.Request) error

// wrapper function for handler function
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle  write json error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

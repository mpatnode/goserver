package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// AddUserResponse object  XXX: Are these required to be exportable for the json module?
type AddUserResponse struct {
	ID          string `json:"id"`
	UserCreated bool   `json:"user_created"`
	Timestamp   int64  `json:"timestamp"`
}

func isErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func help(w http.ResponseWriter, r *http.Request) {
	message := "<html><body><h1>API</h1><ol>"
	message += "<li>/0.9/user POST # email required</li>"
	message += "<li>/0.9/user/<key> GET # key can be ID or email</li>"
	message += "<li>/0.9/user/<key> DELETE # key can be ID or email</li>"
	message += "</ol></body></html>"

	w.Write([]byte(message))
}

func userCreate(w http.ResponseWriter, r *http.Request) {
	var newUser User
	var response AddUserResponse
	reqBody, err := ioutil.ReadAll(r.Body)
	if !isErr(w, err) {
		json.Unmarshal(reqBody, &newUser)
		if strings.Index(newUser.Email, "@") < 1 {
			http.Error(w, "User email is required", http.StatusNotAcceptable)
			return
		}
		if u, _ := DBGetUser(newUser.Email); u != nil {
			// XXX: Perhaps do an update here, though that's better handled by PUT semantics
			w.WriteHeader(http.StatusOK)
			response.ID = u.ID
			response.Timestamp = u.CreatedAt
			response.UserCreated = false
		} else if !isErr(w, StripeCreateCustomer(&newUser)) && !isErr(w, BTCreateCustomer(&newUser)) && !isErr(w, DBAddUser(&newUser)) {
			w.WriteHeader(http.StatusCreated)
			response.ID = newUser.ID
			response.Timestamp = time.Now().Unix()
			response.UserCreated = true
		}
		json.NewEncoder(w).Encode(response)
	}
}

func userFind(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	u, err := DBGetUser(key)

	if !isErr(w, err) {
		if u == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(*u)
	}
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	err := DBDeleteUser(key)
	if !isErr(w, err) {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", help)
	router.HandleFunc("/0.9/user", userCreate).Methods("POST")
	router.HandleFunc("/0.9/user/{key}", userFind).Methods("GET")
	router.HandleFunc("/0.9/user/{key}", userDelete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

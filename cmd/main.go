package main

import (
	"codelabx/auth"
	"codelabx/models"
	"codelabx/repos"
	"codelabx/ws"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	manager = ws.NewManager()
)

func main() {
	var r = mux.NewRouter()

	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx auth xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/authorize", authorize).Methods("POST")

	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx websocket xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	r.HandleFunc("/handShake", handShake).Methods("GET")

	go manager.ListenToRedis()

	defer log.Fatal(http.ListenAndServe(":8010", r))
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx auth xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var incomingUser models.User //get incoming user
	json.NewDecoder(r.Body).Decode(&incomingUser)

	if repos.UserExists(&incomingUser) { //checking if identical user exists
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("user Exists all ready")
		return
	}
	res := repos.AddUser(&incomingUser) // adding user to db

	if res == 1 {
		// json.NewEncoder(w).Encode(&incomingUser)

		tokenStr := auth.CreateToken(&incomingUser) //creating token

		mp := map[string]string{"token": tokenStr, " message": "Welcome to CodeLabX"}
		json.NewEncoder(w).Encode(&mp) // sending token to user
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error Happened during process please try again.")
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var incomingUser models.User
	json.NewDecoder(r.Body).Decode(&incomingUser)

	if !repos.UserExists(&incomingUser) { // checking if user exists
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User does not exist, check credentials")
		return
	}

	if repos.IsValidPassword(&incomingUser) {
		tokenStr := auth.CreateToken(&incomingUser)
		mp := map[string]string{"token": tokenStr, " message": "Welcome to CodeLabX"}
		json.NewEncoder(w).Encode(&mp) // sending token to user
		return
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("Invalid Password, Please check credentials")
	return
}

func authorize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mp := map[string]string{}
	mp["token"] = ""

	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		log.Fatal("err in reading : ", err)
	}

	if mp["token"] != "" {
		ans := auth.IsAuthorized(mp["token"])
		if ans {
			json.NewEncoder(w).Encode("Authorised token")
		} else {
			json.NewEncoder(w).Encode("Un-authorised token")
		}
		return
	}
	fmt.Println("token is empty")
	return
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx websocket xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func handShake(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("token")
	ok := auth.IsAuthorized(token)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Token was Unauthorised...")
		return
	}
	manager.ServeWs(w, r)
}

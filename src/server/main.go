package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "get, Hello " + message

		fmt.Fprintf(w, message)
	} else {
		http.Error(w, "Invalid Method, GET only", http.StatusMethodNotAllowed)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		data, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(w, string(data))
	} else {
		http.Error(w, "Error Method is "+r.Method, http.StatusMethodNotAllowed)
	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get", getHandler).Methods("GET")
	r.HandleFunc("/post", postHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", r))
}

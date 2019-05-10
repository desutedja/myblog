package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ModelUser struct {
	ID        int
	FirstName string
	LastName  string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM user ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	usr := ModelUser{}
	res := []ModelUser{}
	for selDB.Next() {
		var id int
		var firstname, lastname string
		err = selDB.Scan(&id, &firstname, &lastname)
		if err != nil {
			panic(err.Error())
		}
		usr.ID = id
		usr.FirstName = firstname
		usr.LastName = lastname
		res = append(res, usr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	defer db.Close()
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbConn()
	SID := mux.Vars(r)
	selDB, err := db.Query("SELECT * FROM user WHERE ID =?", SID["id"])
	if err != nil {
		panic(err.Error())
	}
	usr := ModelUser{}
	for selDB.Next() {
		var id int
		var firstname, lastname string
		err = selDB.Scan(&id, &firstname, &lastname)
		if err != nil {
			panic(err.Error())
		}
		usr.ID = id
		usr.FirstName = firstname
		usr.LastName = lastname
	}

	json.NewEncoder(w).Encode(usr)
	defer db.Close()
}

func createUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		var user ModelUser
		_ = json.NewDecoder(r.Body).Decode(&user)
		firstname := user.FirstName
		lastname := user.LastName
		insForm, err := db.Prepare("INSERT INTO user(FirstName, LastName) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(firstname, lastname)
		log.Println("INSERT: First Name: " + firstname + " | Last Name: " + lastname)
	}

	defer db.Close()
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "PUT" {
		var user ModelUser
		_ = json.NewDecoder(r.Body).Decode(&user)
		firstname := user.FirstName
		lastname := user.LastName
		uid := user.ID
		updtForm, err := db.Prepare("UPDATE user SET FirstName =?, LastName=? WHERE ID=?")
		if err != nil {
			panic(err.Error())
		}
		updtForm.Exec(firstname, lastname, uid)
		log.Println("UPDATE First Name: " + firstname + " | Last Name: " + lastname + " WITH ID : " + string(uid))
	}

	defer db.Close()
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	uid := mux.Vars(r)
	delForm, err := db.Prepare("DELETE FROM user WHERE ID=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(uid["id"])
	defer db.Close()
	fmt.Fprintf(w, "DELETE user Dengan ID "+uid["id"]+" Berhasil!")
	log.Println("DELETE user Dengan ID " + string(uid["id"]) + " Berhasil!")
}

func main() {
	r := mux.NewRouter()

	//routes
	r.HandleFunc("/api/users", getUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/api/users", createUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", r))

}

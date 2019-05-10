package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/kataras/go-sessions"
)

type userModel struct {
	Id        int
	UserName  string
	FirstName string
	LastName  string
	Password  string
}

func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "myblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func QueryUser(uname string) userModel {
	db := DbConn()
	usr := userModel{}
	db.QueryRow("SELECT Id,UserName,FirstName,LastName,Password FROM user WHERE UserName =?", uname).Scan(
		&usr.Id, &usr.UserName, &usr.FirstName, &usr.LastName, &usr.Password,
	)
	return usr
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}

	return true
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "Views/register.html")
		return
	}

	db := DbConn()
	uname := r.FormValue("username")
	fname := r.FormValue("firstname")
	lname := r.FormValue("lastname")
	pwd := r.FormValue("password")

	users := QueryUser(uname)

	if (userModel{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && checkErr(w, r, err) {
			regQuery, err := db.Prepare("INSERT INTO user(UserName,FirstName,LastName,Password) VALUES(?,?,?,?)")
			if err == nil {
				_, err := regQuery.Exec(uname, fname, lname, hashedPassword)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
	defer db.Close()
}

func login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) != 0 {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method != "POST" {
		http.ServeFile(w, r, "Views/login.html")
		return
	}

	uname := r.FormValue("username")
	pwd := r.FormValue("password")

	users := QueryUser(uname)

	pwdCompare := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(pwd))

	if pwdCompare == nil {
		//success
		session := sessions.Start(w, r)
		session.Set("username", users.UserName)
		session.Set("name", users.FirstName)
		http.Redirect(w, r, "/", 302)
	} else {
		//fail
		http.Redirect(w, r, "/login", 302)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/login", 302)
	}

	data := map[string]string{
		"username": session.GetString("username"),
		"message":  "Welcome on Go !",
	}

	t, err := template.ParseFiles("Views/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, data)
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}

/*func routes() {
	//routesnya
	r.HandleFunc("/register", register).Methods("POST")
	//r.HandleFunc("/login", login).Methods("POST")
	fmt.Println("Routes Running")
}*/

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", register)
	r.HandleFunc("/login", login)
	r.HandleFunc("/", home)

	fmt.Println("Server running on port :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

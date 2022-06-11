package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Context struct {
	Title string
	Name  string
	Users []*User
}

const StatsDoc = `
<!DOCTYPE html>
<html>
    <head>
        {{.Title}}
    </head>
    <body>
        <h3>Hi, {{.Name}}. Here is list of db users:</h3>
        <ul>
            {{range $key, $val := .Users}}
                <li>ID:{{$key}}, Name : {{$val.Name}}</li>
            {{end}}
        </ul>
    </body>
</html>
`

func getUsers() []*User {
	// Open up our database connection.
	db, err := sql.Open("mysql", "tester:secret@tcp(db:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	// Execute the query
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var users []*User
	for results.Next() {
		var u User
		// for each row, scan the result into our tag composite object
		err = results.Scan(&u.ID, &u.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		users = append(users, &u)
	}

	return users
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/static/stats.html")
}

func statPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content Type", "text/html")
	users := getUsers()

	templates := template.New("template")
	templates.New("doc").Parse(StatsDoc)
	context := Context{
		Title: "Statistics",
		Name:  "Admin",
		Users: users,
	}
	templates.Lookup("doc").Execute(w, context)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/statistics", statPage)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

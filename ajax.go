package edm

import (
	"fmt"
	"encoding/json"
	"net/http"
	"appengine"
	_ "google.golang.org/appengine/cloudsql"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type SqlSystem struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

func init() {
	http.HandleFunc("/ajax", handleAjax)
}

func handleAjax(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	if (action == "systemsearch") {
		handleSystemSearch(w, r)
	} else {
		fmt.Fprint(w, "No match for action")
	}
}

func handleSystemSearch(w http.ResponseWriter, r *http.Request) {
	system := r.FormValue("name")
	if system == "" {
		reportError(w, fmt.Sprintf("Expected param \"name\" not found"), 500)
		return
	}
	db, err := sqlConnect("systems")
	if err != nil {
		reportError(w, fmt.Sprintf("Error opening connection to database: %v", err), 500)
		return
	}
	rows, err := db.Query("SELECT id, name FROM systems WHERE name LIKE ? LIMIT 20", "%" + system + "%")
	if err != nil {
		reportError(w, fmt.Sprintf("Error running query: %v", err), 500)
		return
	}
	defer rows.Close()

	results := []SqlSystem{};
	for rows.Next() {
		var id int
		var name string
		_ = rows.Scan(&id, &name)
		// fmt.Fprintf(w, "Found %v (id: %v)\n", name, id)
		results = append(results, SqlSystem{id, name})
	}
	w.Header().Set("Content-Type", "application/javascript")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(results)
	if err != nil {
		reportError(w, fmt.Sprintf("Error encoding result json: %v", err), 500)
		return
	}
}

func sqlConnect(database string) (db *sql.DB, err error) {
	var datasource string
	if appengine.IsDevAppServer() {
		datasource = "readonly:readonly@tcp(173.194.252.109:3306)/" + database
	} else {
		datasource = "readonly@cloudsql(balmy-moonlight-372:galaga-sql-02)/" + database
	}
	db, err = sql.Open("mysql", datasource)
	return
}

func reportError(w http.ResponseWriter, message string, code int) {
	if !appengine.IsDevAppServer() {
		message = http.StatusText(code)
	}
	http.Error(w, message, code)
	return
}

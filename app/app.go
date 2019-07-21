package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

// Employee -
type Employee struct {
	ID            int    `json:"id"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	Email         string `json:"email"`
	Gender        string `json:"gender"`
	FavoriteColor string `json:"favoritecolor"`
}

var db *sql.DB

// Start up initialization
func init() {
	// Example conntection string
	// postgres://[USERNAME]:[PASSWORD]@[HOSTNAME]/[DATABASE]?sslmode=disable
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:pass123@localhost/sampledata?sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You are connected to your database.")
}

func main() {
	http.HandleFunc("/api", index)
	http.HandleFunc("/employees", employeeIndex)
	http.HandleFunc("/employees/find", getEmployeeByName)
	http.ListenAndServe(":8080", nil)
}

//
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "api is up!\nMETHOD="+r.Method)
}

func employeeIndex(w http.ResponseWriter, r *http.Request) {
	// Check the http method type
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	// Query the database
	rows, err := db.Query("SELECT * FROM employees;")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Close connection when finished
	defer rows.Close()

	// Scan the rows using the struct expected from the query and append to slice
	employees := make([]Employee, 0)
	for rows.Next() {
		employee := Employee{}
		err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Gender, &employee.FavoriteColor)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		employees = append(employees, employee)
	}
	// Check for error during interation
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Print the http response in pain text - DEPRECATED
	// for _, emp := range employees {
	// 	// fmt.Println(bk.isbn, bk.title, bk.author, bk.price)
	// 	fmt.Fprintf(w, "%v, %s, %s, %s, %s, %s\n", emp.ID, emp.FirstName, emp.LastName, emp.Email, emp.Gender, emp.FavoriteColor)
	// }

	// Marshal to JSON
	empsjson, err := json.Marshal(employees)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(empsjson)
}

func getEmployeeByName(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// Check the key=value in the http query string
	fname := r.FormValue("firstname")
	if fname == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM employees WHERE first_name = $1", fname)

	employee := Employee{}
	err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Gender, &employee.FavoriteColor)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	// Print the http response in pain text - DEPRECATED
	// fmt.Fprintf(w, "%v, %s, %s, %s, %s, %s\n", employee.ID, employee.FirstName, employee.LastName, employee.Email, employee.Gender, employee.FavoriteColor)

	empjson, err := json.Marshal(employee)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(empjson)
}

package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

// Employee -
type Employee struct {
	id            int
	firstName     string
	lastName      string
	email         string
	gender        string
	favoriteColor string
}

var db *sql.DB

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
	http.HandleFunc("/", index)
	http.HandleFunc("/employees", employeeIndex)
	http.HandleFunc("/employees/find", getEmployeeByName)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index is up!\nMETHOD="+r.Method)
}

func employeeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM employees;")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		employee := Employee{}
		err := rows.Scan(&employee.id, &employee.firstName, &employee.lastName, &employee.email, &employee.gender, &employee.favoriteColor) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		employees = append(employees, employee)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, emp := range employees {
		// fmt.Println(bk.isbn, bk.title, bk.author, bk.price)
		fmt.Fprintf(w, "%v, %s, %s, %s, %s, %s\n", emp.id, emp.firstName, emp.lastName, emp.email, emp.gender, emp.favoriteColor)
	}

}

func getEmployeeByName(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	fname := r.FormValue("firstname")
	if fname == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM employees WHERE first_name = $1", fname)

	employee := Employee{}
	err := row.Scan(&employee.id, &employee.firstName, &employee.lastName, &employee.email, &employee.gender, &employee.favoriteColor)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%v, %s, %s, %s, %s, %s\n", employee.id, employee.firstName, employee.lastName, employee.email, employee.gender, employee.favoriteColor)
}

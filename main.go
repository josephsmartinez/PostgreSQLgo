package main

import (
	"database/sql"
	"fmt"

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

func main() {

	// Example conntection string
	// postgres://[USERNAME]:[PASSWORD]@[HOSTNAME]/[DATABASE]?sslmode=disable
	db, err := sql.Open("postgres", "postgres://postgres:pass123@localhost/sampledata?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You are connected to your database.")

	rows, err := db.Query("SELECT * FROM employees;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		employee := Employee{}
		//
		err := rows.Scan(&employee.id, &employee.firstName, &employee.lastName, &employee.email, &employee.gender, &employee.favoriteColor) // order matters
		if err != nil {
			panic(err)
		}
		employees = append(employees, employee)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, emp := range employees {
		// fmt.Println(bk.isbn, bk.title, bk.author, bk.price)
		fmt.Printf("%v, %s, %s, %s, %s, %s\n", emp.id, emp.firstName, emp.lastName, emp.email, emp.gender, emp.favoriteColor)
	}
}

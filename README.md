# PostgreSQLgo

PostgreSQL with Golang

- Build all internal search functions first
- Take args from the command prompt to query the database.

Start app
> sh postgres-setup.sh
> go run app/app.go

Queries
> curl -i localhost:8080/employees
> curl -i localhost:8080/employees/find?firstname=Rich

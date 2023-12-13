# chi-user-account-crelog
User Account CRUD, login and logout using sessions. Made with go-chi router and PostgreSQL with Docker containerisation.

This link https://odmg.dev/project2 shows images of the request and response for some of the endpoints: create_user, user_login and fetching_user_by_id.


## Automating Program Compilation with a Makefile
- To build target use: make build-app
- To run target use: make run
- To run API inside docker container use: make docker

## Project environment variables
- HTTP_LISTEN_ADDRESS=:9090
- DB_HOST=
- DB_PORT=
- DB_USER=
- DB_PASSWORD=
- DB_NAME=
- DB_SSLMODE=

## Docker
### Installing postgres as a Docker container
- docker run --name postgresdb -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres:latest


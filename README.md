# Authentication Service

This is an authentication service implemented in Go. It uses JWT for authorization and Redis as a store.

## Project Structure

The project is structured as follows:

- `cmd/api/main.go`: The entry point of the application.
- `internal/authorization/`: Contains the authorization logic.
- `internal/authservice/`: Contains the main authentication service logic.
- `internal/config/`: Contains the logic for loading configuration.
- `internal/model/`: Contains the data models.
- `internal/store/`: Contains the data store implementations.
- `internal/transport/`: Contains the HTTP transport logic.

## Dependencies

The project uses the following dependencies:

- `github.com/golang-jwt/jwt/v5`
- `github.com/go-chi/chi/v5`
- `github.com/google/uuid`
- `github.com/redis/go-redis/v9`
- `golang.org/x/crypto`

## Building and Running the Project

You can build and run the project using Docker Compose. The `Makefile` contains the commands for building and running the project:

```sh
make run
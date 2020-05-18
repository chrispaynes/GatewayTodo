# vorChall
vorChall is an app for managing a Todo list using [Go](https://golang.org/), [Angular](https://angular.io/), [PostgreSQL](https://www.postgresql.org/), [gRPC](https://grpc.io/) and [12-Factor](https://12factor.net/).

### Quickstart
Separate Docker containers exists for the API, UI and the Database.

From the project root, run `$ make todoApp` to get the necessary services started. More Make Recipes can be found in the Makefile

- The REST API is available at `http://127.0.0.1:3000`
- The gRPC API is available at `http://127.0.0.1:3001`
- The UI is available at `http://127.0.0.1:8080` and `http://127.0.0.1:4200`
- Postgres is available at `http://127.0.0.1:5432` with [credentials](###Postgres) found below.

##### Environment Variables
Environment variables are defined in the Docker-Compose file as well as `./vorChall/db/postgres.env`.
```
$ make todoApp
```

### Postgres
- Postgres serves as a storage mechanism for Todos. Basic details about a Todo are stored in the `todo` table and are joined with a `todo_status` table that stores various statuses that a Todo can be in.
- In a local development setup, the database schema definition `./vorChall/db/ddl` is copied into the Postgres container and executed as part of the container's initialization. The database is also seeded with several dummy Todos.
- Passwords for initializing the Database are stored in `./vorChall/db/postgres.env` in the format of:
    ```
    POSTGRES_USER=read_write_user
    POSTGRES_PASSWORD=<INSERT PASSWORD HERE>
    POSTGRES_DB=todos
    POSTGRES_HOST=postgres
    ENABLE_SSL=false
    ```
- If you run into file permission issues mounting `/var/lib/postgresql/data` with the Postgres container, then remove the mount definition from the docker-compose file and restart the Postgres service.

### Angular UI
The Angular application is served by building the application for production and then serving the static HTML, CSS and JavaScript files through an [NGINX](https://www.nginx.com/) container. For development, the UI can optionally be run on the default port of `:4200` using the [Angular CLI](https://cli.angular.io/) `ng serve` command or as a separate Goroutine in the API container.

### REST / gRPC API 
The API is defined as [Protocol Buffers (Protobufs)](https://developers.google.com/protocol-buffers/docs/proto) with gRPC and REST endpoints over TCP. The Protobuf definition `./vorChall/api/proto/v1.proto` is compiled by [Proto-Gen-Go](https://github.com/golang/protobuf/tree/master/protoc-gen-go) using instructions defined in `./vorChall/api/proto/proto.go`. The compile step outputs a Go file with [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) functionality, as well as Swagger documentation.

### Swagger Documentation and Postman Collection
A [Postman](https://www.postman.com/) environment file `./vorChall/api/postman/dev.postman_environment.json` and swagger documentation `./vorChall/api/proto/api/swagger/todos.swagger.json` are available to import into Postman and interact with the API.


### API Dependencies
API Dependencies are managed with [Go Modules](https://blog.golang.org/using-go-modules) and are vendored at `./vorChall/api/vendor`.

---

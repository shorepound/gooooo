# Gooooo — Minimal Go REST API

Simple REST service example (Go 1.25). Endpoints:

- `GET /health` — health check
- `GET /items` — list items
- `POST /items` — create item (JSON)
- `GET /items/{id}` — get item
- `PUT /items/{id}` — update item (JSON)
- `DELETE /items/{id}` — delete item

Run locally:

```bash
# ensure module deps
go mod tidy
# build
go build ./...
# run
go run ./cmd/server
```

Example curl:

```bash
# create
curl -s -X POST -H "Content-Type: application/json" -d '{"name":"foo","description":"test"}' http://localhost:8080/items
# list
curl -s http://localhost:8080/items
```

Run tests:

```bash
go test ./...
```

Notes:
- Server supports graceful shutdown (SIGINT/SIGTERM).

Docker Compose (local Postgres):

```bash
# start Postgres
docker compose up -d
# stop
docker compose down
```

Migrations:

Migrations live in `./migrations` and are run automatically at server startup using `golang-migrate`.
To run migrations manually with the CLI, you can install `migrate` and run:

```bash
migrate -path ./migrations -database "${DATABASE_URL}" up
```

Environment:

- `DATABASE_URL` (optional): Postgres DSN, e.g. `postgres://postgres:password@localhost:5432/gooooo?sslmode=disable`

Docker image (build/run):

```bash
# build image
docker build -t gooooo:latest .
# run image (ensure DATABASE_URL is set appropriately)
docker run --rm -p 8080:8080 -e DATABASE_URL="${DATABASE_URL}" gooooo:latest
```

Task shortcuts (`task`):

If you use `go-task` (https://taskfile.dev/), `Taskfile.yml` contains common tasks:

```bash
task build        # builds the Go project
task run          # runs the server
task test         # runs tests
task docker-build # builds the docker image
task docker-run   # runs the docker image
task compose-up   # starts local Postgres via docker compose
task compose-down # stops compose services
```

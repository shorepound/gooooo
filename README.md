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

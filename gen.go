package gen

//go:generate go run ./cmd/web-monitor/tools/terndotenv/main.go
//go:generate sqlc generate -f ./internal/store/pgstore/sqlc.yaml


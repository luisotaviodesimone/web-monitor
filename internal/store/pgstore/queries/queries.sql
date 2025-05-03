-- name: InsertPing :one
INSERT INTO
  ping_logs (
    "host",
    "latency_avg_ms",
    "loss_rate_percent",
    "ping_count"
  )
VALUES
  ($1, $2, $3, $4)
RETURNING
  "id";

-- name: InsertHttp :one
INSERT INTO
  http_logs ("host", "latency_avg_ms", "http_count")
VALUES
  ($1, $2, $3)
RETURNING
  "id";

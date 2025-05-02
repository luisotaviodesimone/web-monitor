package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/web-monitor/cmd/web-monitor/probes"
	"github.com/web-monitor/internal/store/pgstore"
)

func get(envVar string) string {
	return os.Getenv(envVar)
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		get("DATABASE_USER"),
		get("DATABASE_PASSWORD"),
		get("DATABASE_HOST"),
		get("DATABASE_PORT"),
		get("DATABASE_NAME"),
	)
	pool, err := pgxpool.New(ctx, connString)

	if err != nil {
		slog.Error("Unable to connect to database", slog.Any("error", err))
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Unable to ping database", slog.Any("error", err))
		panic(err)
	}

	queries := pgstore.New(pool)

	slog.Info("Connected to database")

	hosts := []string{"google.com", "rnp.br", "youtube.com"}
	for {

		for _, host := range hosts {
			count := 11

			averageLatency, packetLossPercentage, err := probes.Ping(host, count)

			if err != nil {
				slog.Error("Error pinging host", slog.String("host", host), slog.Any("error", err))
			}

			queries.InsertPing(ctx, pgstore.InsertPingParams{
				LatencyAvgMs:    pgtype.Int4{Int32: int32(averageLatency)},
				LossRatePercent: pgtype.Int4{Int32: int32(packetLossPercentage)},
				PingCount:       pgtype.Int4{Int32: int32(count)},
			})

			averageLatency, err = probes.HttpPing(host, count)

			if err != nil {
				slog.Error("Error pinging host", slog.String("host", host), slog.Any("error", err))
			}

			queries.InsertHttp(ctx, pgstore.InsertHttpParams{
				LatencyAvgMs: pgtype.Int4{Int32: int32(averageLatency)},
				HttpCount:    pgtype.Int4{Int32: int32(count)},
			})
		}
	}
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func WaitForPostgres(dsn string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error
	for range ticker.C {
		if time.Now().After(deadline) {
			break
		}

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			lastErr = err
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		pingErr := db.PingContext(ctx)
		cancel()
		_ = db.Close()

		if pingErr == nil {
			return nil
		}
		lastErr = pingErr
	}
	return fmt.Errorf("postgres did not rise for %s: %w", timeout, lastErr)
}

package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/patato8984/Shop/internal/shared/dto"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"go.uber.org/zap"
)

type Event struct {
	Topic      string
	Event_type string          `json:"type"`
	Payload    json.RawMessage `json:"payload"`
}
type OutboxRepo struct {
	db *sql.DB
}
type Worker struct {
	db        *sql.DB
	publisher *shared_events.EventPublisher
	log       *zap.Logger
}

func NewOutboxRepo(db *sql.DB) *OutboxRepo {
	return &OutboxRepo{db: db}
}
func NewOutboxWorker(db *sql.DB, publisher *shared_events.EventPublisher, log *zap.Logger) *Worker {
	return &Worker{db: db, publisher: publisher, log: log}
}
func (r *OutboxRepo) Add(ctx context.Context, topic, eventType string, payLoad any) error {
	db := dto.Getter(ctx, r.db)
	data, _ := json.Marshal(&payLoad)
	_, err := db.ExecContext(ctx, "INSERT INTO outbox (topic, event_type, payload) VALUES ($1,$2,$3)", topic, eventType, data)
	return err
}
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	w.log.Info("Outbox worker started")

	for {
		select {
		case <-ctx.Done():
			w.log.Info("Outbox worker stopped by context")
			return
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				w.log.Error("worker batch error", zap.Error(err))
			}
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	rows, err := tx.QueryContext(ctx, "SELECT id, topic, event_type, payload FROM outbox WHERE status = 'new' ORDER BY created_at ASC LIMIT 50 FOR UPDATE SKIP LOCKED")
	if err != nil {
		return fmt.Errorf("select outbox: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var e Event
		var id uuid.UUID
		if err := rows.Scan(&id, &e.Topic, &e.Event_type, &e.Payload); err != nil {
			return err
		}
		if err := w.publisher.Publisher(ctx, e.Topic, e.Event_type, e.Payload); err != nil {
			w.log.Warn("failed to publish event", zap.Error(err))
			return err
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil
	}
	_, err = tx.ExecContext(ctx, "UPDATE outbox SET status = 'processed', processed_at = NOW() WHERE id = ANY($1)", pq.Array(ids))
	if err != nil {
		return fmt.Errorf("update outbox: %w", err)
	}

	return tx.Commit()
}

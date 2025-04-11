package postgres

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type ExecutionLogRepository struct {
    db *pgxpool.Pool
}

func NewExecutionLogRepository(db *pgxpool.Pool) *ExecutionLogRepository {
    return &ExecutionLogRepository{db: db}
}

func (r *ExecutionLogRepository) Save(exec entity.ExecutionLog) error {
    query := `
        INSERT INTO executions (
            id, bot_id, entry_price, entry_time, exit_price, exit_time,
            duration, profit, roi_pct, strategy, created_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, now()
        )
    `
    _, err := r.db.Exec(context.Background(), query,
        uuid.New(), exec.BotID, exec.Entry.Price, exec.Entry.Timestamp,
        exec.Exit.Price, exec.Exit.Timestamp, exec.Duration, exec.Profit, exec.ROIPct,
        exec.Strategy.Name,
    )
    return err
}

func (r *ExecutionLogRepository) GetAll() ([]entity.ExecutionLog, error) {
    query := `
        SELECT bot_id, entry_price, entry_time, exit_price, exit_time,
               duration, profit, roi_pct, strategy
        FROM executions
        ORDER BY created_at DESC
    `
    rows, err := r.db.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []entity.ExecutionLog
    for rows.Next() {
        var e entity.ExecutionLog
        err := rows.Scan(
            &e.BotID, &e.Entry.Price, &e.Entry.Timestamp,
            &e.Exit.Price, &e.Exit.Timestamp, &e.Duration,
            &e.Profit, &e.ROIPct, &e.Strategy.Name,
        )
        if err != nil {
            return nil, err
        }
        logs = append(logs, e)
    }

    return logs, nil
}
package postgres

import (
    "context"
    "encoding/json"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type DecisionLogRepository struct {
    db *pgxpool.Pool
}

func NewDecisionLogRepository(db *pgxpool.Pool) *DecisionLogRepository {
    return &DecisionLogRepository{db: db}
}

func (r *DecisionLogRepository) Save(log entity.DecisionLog) error {
    indicatorsJSON, _ := json.Marshal(log.Indicators)
    contextJSON, _ := json.Marshal(log.Context)

    query := `
        INSERT INTO decisions (
            id, bot_id, symbol, interval, timestamp, decision, price,
            indicators, context, strategy, created_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7,
            $8, $9, $10, now()
        )
    `
    _, err := r.db.Exec(context.Background(), query,
        uuid.New(), log.BotID, log.Symbol, log.Interval,
        log.Timestamp, log.Decision, log.Indicators["price"],
        indicatorsJSON, contextJSON, log.Strategy.Name,
    )
    return err
}
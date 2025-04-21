// internal/infra/repository/postgres/postgres_bot_repository.go

package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
)

type BotRepository struct {
	db *pgxpool.Pool
}

func NewBotRepository(db *pgxpool.Pool) *BotRepository {
	return &BotRepository{db: db}
}

func (r *BotRepository) Create(bot *entity.Bot) (*entity.BotWithStrategy, error) {
	query := `
    INSERT INTO bots (id, account_id, strategy_id, symbol, interval, autonomous, active, config_json, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
`
	_, err := r.db.Exec(context.Background(), query,
		bot.ID, bot.AccountID, bot.StrategyID, bot.Symbol, bot.Interval,
		bot.Autonomous, bot.Active, bot.ConfigJSON,
	)
	if err != nil {
		return nil, err
	}

	borCreated, err := r.GetByID(bot.ID)
	if err != nil {
		return nil, err
	}

	return borCreated, nil
}

func (r *BotRepository) GetByID(id uuid.UUID) (*entity.BotWithStrategy, error) {
	query := `
	SELECT 
		b.id, b.account_id, b.symbol, b.interval, b.autonomous, b.active,
		b.config_json, b.created_at, b.updated_at,
		s.id AS strategy_id, s.name AS strategy_name
	FROM bots b
	JOIN strategies s ON b.strategy_id = s.id
	WHERE b.id = $1
`

	row := r.db.QueryRow(context.Background(), query, id)

	var b entity.BotWithStrategy
	err := row.Scan(
		&b.ID, &b.AccountID, &b.Symbol, &b.Interval, &b.Autonomous, &b.Active,
		&b.ConfigJSON, &b.CreatedAt, &b.UpdatedAt,
		&b.StrategyID, &b.StrategyName,
	)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *BotRepository) GetByAccountID(accountID uuid.UUID) ([]entity.BotWithStrategy, error) {
	logger.Debug("Buscando bots para o account_id: ", accountID)

	query := `
		SELECT 
			b.id, b.account_id, b.symbol, b.interval, b.autonomous, b.active,
			b.config_json, b.created_at, b.updated_at,
			s.id AS strategy_id, s.name AS strategy_name
		FROM bots b
		JOIN strategies s ON b.strategy_id = s.id
		WHERE b.account_id = $1
	`

	rows, err := r.db.Query(context.Background(), query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []entity.BotWithStrategy
	for rows.Next() {
		var b entity.BotWithStrategy
		if err := rows.Scan(
			&b.ID, &b.AccountID, &b.Symbol, &b.Interval, &b.Autonomous, &b.Active,
			&b.ConfigJSON, &b.CreatedAt, &b.UpdatedAt,
			&b.StrategyID, &b.StrategyName,
		); err != nil {
			return nil, err
		}
		bots = append(bots, b)
	}

	logger.Debug("Total de bots encontrados: ", len(bots))

	return bots, nil
}

func (r *BotRepository) Update(bot *entity.Bot) (*entity.BotWithStrategy, error) {
	query := `
		UPDATE bots
		SET symbol = $1, interval = $2, strategy_id = $3, autonomous = $4,
			active = $5, config_json = $6, updated_at = now()
		WHERE id = $7
	`
	_, err := r.db.Exec(context.Background(), query,
		bot.Symbol, bot.Interval, bot.StrategyID,
		bot.Autonomous, bot.Active, bot.ConfigJSON, bot.ID,
	)
	if err != nil {
		return nil, err
	}

	botUpdate, err := r.GetByID(bot.ID)
	if err != nil {
		return nil, err
	}

	return botUpdate, nil
}

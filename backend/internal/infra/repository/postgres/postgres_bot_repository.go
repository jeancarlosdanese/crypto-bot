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

func (r *BotRepository) Create(bot *entity.Bot) (*entity.Bot, error) {
	query := `
        INSERT INTO bots (id, account_id, symbol, interval, strategy_name, autonomous, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
    `
	_, err := r.db.Exec(context.Background(), query,
		bot.ID, bot.AccountID, bot.Symbol, bot.Interval, bot.StrategyName,
		bot.Autonomous, bot.Active,
	)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func (r *BotRepository) GetByID(id uuid.UUID) (*entity.Bot, error) {
	query := `SELECT id, account_id, symbol, interval, strategy_name, autonomous, active FROM bots WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, id)

	var b entity.Bot
	err := row.Scan(&b.ID, &b.AccountID, &b.Symbol, &b.Interval, &b.StrategyName, &b.Autonomous, &b.Active)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *BotRepository) GetByAccountID(accountID uuid.UUID) ([]entity.Bot, error) {
	logger.Debug("Buscando bots para o account_id: ", accountID)

	query := `SELECT id, account_id, symbol, interval, strategy_name, autonomous, active FROM bots WHERE account_id = $1`
	rows, err := r.db.Query(context.Background(), query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []entity.Bot
	for rows.Next() {
		var b entity.Bot
		if err := rows.Scan(&b.ID, &b.AccountID, &b.Symbol, &b.Interval, &b.StrategyName, &b.Autonomous, &b.Active); err != nil {
			return nil, err
		}
		bots = append(bots, b)
	}

	logger.Debug("Total de bots encontrados: ", len(bots))

	return bots, nil
}

func (r *BotRepository) Update(bot *entity.Bot) (*entity.Bot, error) {
	query := `
        UPDATE bots
        SET symbol = $1, interval = $2, strategy_name = $3, autonomous = $4, active = $5, updated_at = now()
        WHERE id = $6
    `
	_, err := r.db.Exec(context.Background(), query,
		bot.Symbol, bot.Interval, bot.StrategyName, bot.Autonomous, bot.Active, bot.ID,
	)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

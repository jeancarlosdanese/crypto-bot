package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type PositionRepository struct {
	db *pgxpool.Pool
}

func NewPositionRepository(db *pgxpool.Pool) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) Save(p entity.OpenPosition) error {
	query := `
        INSERT INTO positions (id, bot_id, entry_price, timestamp)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (bot_id) DO UPDATE SET entry_price = EXCLUDED.entry_price, timestamp = EXCLUDED.timestamp
    `
	_, err := r.db.Exec(context.Background(), query,
		uuid.New(), p.BotID, p.EntryPrice, p.Timestamp,
	)
	return err
}

func (r *PositionRepository) GetAll() ([]entity.OpenPosition, error) {
	query := `SELECT bot_id, entry_price, timestamp FROM positions`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []entity.OpenPosition
	for rows.Next() {
		var p entity.OpenPosition
		err := rows.Scan(&p.BotID, &p.EntryPrice, &p.Timestamp)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}

	return positions, nil
}

func (r *PositionRepository) Get(botID uuid.UUID) (*entity.OpenPosition, error) {
	query := `SELECT bot_id, entry_price, timestamp FROM positions WHERE bot_id = $1`
	row := r.db.QueryRow(context.Background(), query, botID)

	var p entity.OpenPosition
	err := row.Scan(&p.BotID, &p.EntryPrice, &p.Timestamp)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PositionRepository) Delete(botID uuid.UUID) error {
	query := `DELETE FROM positions WHERE bot_id = $1`
	_, err := r.db.Exec(context.Background(), query, botID)
	return err
}

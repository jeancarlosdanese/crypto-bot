// internal/infra/repository/postgres/postgres_account_repository.go

package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	query := `
        INSERT INTO accounts (id, name, email, whatsapp, api_key, binance_api_key, binance_api_secret, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
        RETURNING id
    `
	err := r.db.QueryRow(ctx, query,
		account.ID, account.Name, account.Email, account.WhatsApp,
		account.APIKey, account.BinanceAPIKey, account.BinanceAPISecret,
	).Scan(&account.ID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	query := `SELECT id, name, email, whatsapp, api_key, binance_api_key, binance_api_secret FROM accounts WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var a entity.Account
	err := row.Scan(&a.ID, &a.Name, &a.Email, &a.WhatsApp, &a.APIKey, &a.BinanceAPIKey, &a.BinanceAPISecret)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountRepository) GetAll(ctx context.Context) ([]*entity.Account, error) {
	query := `SELECT id, name, email, whatsapp, api_key, binance_api_key, binance_api_secret FROM accounts ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*entity.Account
	for rows.Next() {
		var a entity.Account
		err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.WhatsApp, &a.APIKey, &a.BinanceAPIKey, &a.BinanceAPISecret)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateByID(ctx context.Context, id uuid.UUID, jsonData []byte) (*entity.Account, error) {
	var updates map[string]interface{}
	if err := json.Unmarshal(jsonData, &updates); err != nil {
		return nil, fmt.Errorf("erro ao fazer unmarshal do JSON: %w", err)
	}

	setClause := ""
	args := []interface{}{}
	i := 1
	for k, v := range updates {
		setClause += fmt.Sprintf("%s = $%d,", k, i)
		args = append(args, v)
		i++
	}
	if setClause == "" {
		return nil, fmt.Errorf("nenhum campo para atualizar")
	}
	setClause = setClause[:len(setClause)-1]
	args = append(args, id)

	query := fmt.Sprintf("UPDATE accounts SET %s, updated_at = now() WHERE id = $%d", setClause, len(args))
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *AccountRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *entity.Account) error {
	query := `
        INSERT INTO accounts (id, name, email, whatsapp, is_admin, api_key, binance_api_key, binance_api_secret, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
    `
	_, err := r.db.Exec(context.Background(), query,
		account.ID, account.Name, account.Email, account.WhatsApp,
		account.IsAdmin, account.APIKey, account.BinanceAPIKey, account.BinanceAPISecret,
	)
	return err
}

func (r *AccountRepository) GetByEmail(email string) (*entity.Account, error) {
	query := `SELECT id, name, email, whatsapp, is_admin, api_key, binance_api_key, binance_api_secret FROM accounts WHERE email = $1`
	row := r.db.QueryRow(context.Background(), query, email)

	var a entity.Account
	err := row.Scan(&a.ID, &a.Name, &a.Email, &a.WhatsApp, &a.IsAdmin, &a.APIKey, &a.BinanceAPIKey, &a.BinanceAPISecret)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountRepository) GetByWhatsApp(whatsapp string) (*entity.Account, error) {
	query := `SELECT id, name, email, whatsapp, is_admin, api_key, binance_api_key, binance_api_secret FROM accounts WHERE whatsapp = $1`
	row := r.db.QueryRow(context.Background(), query, whatsapp)

	var a entity.Account
	err := row.Scan(&a.ID, &a.Name, &a.Email, &a.WhatsApp, &a.IsAdmin, &a.APIKey, &a.BinanceAPIKey, &a.BinanceAPISecret)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AccountRepository) GetByAPIKey(apiKey string) (*entity.Account, error) {
	query := `SELECT id, name, email, whatsapp, is_admin, api_key, binance_api_key, binance_api_secret FROM accounts WHERE api_key = $1`
	row := r.db.QueryRow(context.Background(), query, apiKey)

	var a entity.Account
	err := row.Scan(&a.ID, &a.Name, &a.Email, &a.WhatsApp, &a.IsAdmin, &a.APIKey, &a.BinanceAPIKey, &a.BinanceAPISecret)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

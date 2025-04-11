// internal/domain/entity/account.go

package entity

import "github.com/google/uuid"

type Account struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	WhatsApp         string    `json:"whatsapp"`
	APIKey           *string   `json:"api_key"`
	BinanceAPIKey    *string   `json:"binance_api_key"`
	BinanceAPISecret *string   `json:"binance_api_secret"`
}

// IsAdmin verifica se a conta é admin baseada no ID fixo
func (a *Account) IsAdmin() bool {
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	return a.ID == adminID
}

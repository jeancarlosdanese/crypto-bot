// internal/domain/entity/account.go

package entity

import "github.com/google/uuid"

type Account struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	WhatsApp         string    `json:"whatsapp"`
	IsAdmin          bool      `json:"is_admin"`
	APIKey           *string   `json:"api_key"`
	BinanceAPIKey    *string   `json:"binance_api_key"`
	BinanceAPISecret *string   `json:"binance_api_secret"`
}

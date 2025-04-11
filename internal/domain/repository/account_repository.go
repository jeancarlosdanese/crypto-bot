// internal/domain/repository/account_repository.go

package repository

import "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"

type AccountRepository interface {
	GetByEmail(email string) (*entity.Account, error)
	GetByWhatsApp(whatsapp string) (*entity.Account, error)
	GetByAPIKey(apiKey string) (*entity.Account, error)
	Create(account *entity.Account) error
}

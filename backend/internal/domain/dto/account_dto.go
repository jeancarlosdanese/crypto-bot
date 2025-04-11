package dto

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
	"github.com/jeancarlosdanese/crypto-bot/internal/utils"
)

// AccountCreateDTO define os campos necessários para criar uma conta
type AccountCreateDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	WhatsApp string `json:"whatsapp"`
}

// AccountUpdateDTO define os campos permitidos para atualizar uma conta
type AccountUpdateDTO struct {
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	WhatsApp         string  `json:"whatsapp"`
	BinanceAPIKey    *string `json:"binance_api_key"`
	BinanceAPISecret *string `json:"binance_api_secret"`
}

// AccountResponseDTO define a estrutura de resposta para a conta
type AccountResponseDTO struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	WhatsApp         string  `json:"whatsapp"`
	APIKey           *string `json:"api_key,omitempty"`
	BinanceAPIKey    *string `json:"binance_api_key,omitempty"`
	BinanceAPISecret *string `json:"binance_api_secret,omitempty"`
}

// Construtor para resposta formatada
func NewAccountResponseDTO(account *entity.Account) AccountResponseDTO {
	return AccountResponseDTO{
		ID:               account.ID.String(),
		Name:             account.Name,
		Email:            account.Email,
		WhatsApp:         utils.FormatWhatsApp(account.WhatsApp),
		APIKey:           account.APIKey,
		BinanceAPIKey:    account.BinanceAPIKey,
		BinanceAPISecret: account.BinanceAPISecret,
	}
}

// Construtor para Account (entidade)
func (a *AccountCreateDTO) ToEntity() *entity.Account {
	return &entity.Account{
		ID:       uuid.New(),
		Name:     a.Name,
		Email:    a.Email,
		WhatsApp: a.WhatsApp,
		APIKey:   utils.GenerateAPIKeyPtr(),
	}
}

// Validação ao criar conta
func (a *AccountCreateDTO) Validate() error {
	if len(a.Name) < 3 || len(a.Name) > 100 {
		return errors.New("o nome deve ter entre 3 e 100 caracteres")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(a.Email) {
		return errors.New("e-mail inválido")
	}

	// Sanitizar WhatsApp
	re := regexp.MustCompile(`\D`)
	a.WhatsApp = re.ReplaceAllString(a.WhatsApp, "")

	if len(a.WhatsApp) < 10 || len(a.WhatsApp) > 15 {
		return errors.New("o WhatsApp deve ter entre 10 e 15 dígitos")
	}

	return nil
}

// Validação ao atualizar conta
func (a *AccountUpdateDTO) Validate() error {
	if a.Name != "" && (len(a.Name) < 3 || len(a.Name) > 100) {
		return errors.New("o nome deve ter entre 3 e 100 caracteres")
	}

	if a.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(a.Email) {
			return errors.New("e-mail inválido")
		}
	}

	if a.WhatsApp != "" {
		re := regexp.MustCompile(`\D`)
		a.WhatsApp = re.ReplaceAllString(a.WhatsApp, "")

		if len(a.WhatsApp) < 10 || len(a.WhatsApp) > 15 {
			return errors.New("o WhatsApp deve ter entre 10 e 15 dígitos")
		}
	}

	return nil
}

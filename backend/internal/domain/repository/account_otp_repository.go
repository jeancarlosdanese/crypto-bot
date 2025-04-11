// internal/db/account_otp_repository.go

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"
)

type AccountOTPRepository interface {
	FindValidOTP(ctx context.Context, identifier string, otp string) (*uuid.UUID, error)
	CleanExpiredOTPs(ctx context.Context) error
	FindByEmailOrWhatsApp(ctx context.Context, identifier string) (*entity.Account, error)
	StoreOTP(ctx context.Context, accountID string, otp string) error
	GetOTPAttempts(ctx context.Context, identifier string) (int, error)
	IncrementOTPAttempts(ctx context.Context, identifier string) error
	ResetOTPAttempts(ctx context.Context, identifier string) error
}

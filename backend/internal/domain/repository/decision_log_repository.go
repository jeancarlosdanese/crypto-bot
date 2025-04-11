// internal/domain/repository/decision_log_repository.go

package repository

import "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"

type DecisionLogRepository interface {
	Save(log entity.DecisionLog) error
}

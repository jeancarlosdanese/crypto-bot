// internal/domain/repository/execution_log_repository.go

package repository

import "github.com/jeancarlosdanese/crypto-bot/internal/domain/entity"

type ExecutionLogRepository interface {
	Save(log entity.ExecutionLog) error
	GetAll() ([]entity.ExecutionLog, error)
}

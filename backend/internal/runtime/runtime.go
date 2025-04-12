// internal/runtime/runtime.go

package runtime

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jeancarlosdanese/crypto-bot/internal/app/usecases"
)

// Mapa global de bots em execução (em memória)
var BotsMap = struct {
	sync.RWMutex
	Items map[uuid.UUID]*usecases.StrategyUseCase
}{
	Items: make(map[uuid.UUID]*usecases.StrategyUseCase),
}

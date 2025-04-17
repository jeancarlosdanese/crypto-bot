// internal/domain/entity/strategy_context.go

package entity

type StrategyContext struct {
	Account          Account
	Bot              Bot
	PositionQuantity float64
	LastEntryPrice   float64
	LastEntryTime    int64
}

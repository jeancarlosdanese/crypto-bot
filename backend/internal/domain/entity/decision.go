// internal/domain/entity/decision.go

package entity

type Decision struct {
	ID           string  `bson:"id,omitempty"`
	StrategyID   string  `bson:"strategy_id"`   // ID da estratégia associada à decisão
	AssetID      string  `bson:"asset_id"`      // ID do ativo associado à decisão
	Quantity     float64 `bson:"quantity"`      // Quantidade de ativos a serem comprados ou vendidos
	Profit       float64 `bson:"profit"`        // Lucro ou perda da decisão
	StopLoss     float64 `bson:"stop_loss"`     // Preço de stop loss
	TakeProfit   float64 `bson:"take_profit"`   // Preço de realização de lucro
	Price        float64 `bson:"price"`         // Preço no momento da decisão
	DecisionType string  `bson:"decision_type"` // Ex: "BUY", "SELL", "HOLD"
	CreatedAt    string  `bson:"created_at"`    // Data e hora da decisão
}

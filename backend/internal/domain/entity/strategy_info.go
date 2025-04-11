// internal/domain/entity/strategy_info.go

package entity

type StrategyInfo struct {
	Name       string         `bson:"name"`
	Version    string         `bson:"version"`
	Parameters map[string]any `bson:"parameters"`
}

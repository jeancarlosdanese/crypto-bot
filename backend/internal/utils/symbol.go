// internal/utils/symbol.go

package utils

import "strings"

// FormatForBinance converte "BNB/USDT" para "BNBUSDT"
func FormatForBinance(symbol string) string {
	return strings.ToUpper(strings.ReplaceAll(symbol, "/", ""))
}

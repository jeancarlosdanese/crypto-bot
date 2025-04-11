// internal/infra/config/env.go

package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
)

func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		logger.Warn(fmt.Sprintf(".env file not found (normal in Docker): %v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		logger.Warn("Erro ao ler .env: %v", err)
	}
}

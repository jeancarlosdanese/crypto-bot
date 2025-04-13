// cmd/api/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jeancarlosdanese/crypto-bot/internal/domain/repository"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/config"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/database"
	"github.com/jeancarlosdanese/crypto-bot/internal/infra/repository/postgres"
	"github.com/jeancarlosdanese/crypto-bot/internal/logger"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/middlewares"
	"github.com/jeancarlosdanese/crypto-bot/internal/server/routes"
)

func main() {
	logger.InitLogger()
	log.Println("🚀 Iniciando Robô de Crypto...")

	config.LoadEnv(".env")

	// PostgreSQL
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("Erro ao conectar no PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Repositórios
	accountRepo := postgres.NewAccountRepository(pool)
	botRepo := postgres.NewBotRepository(pool)
	positionRepo := postgres.NewPositionRepository(pool)
	executionRepo := postgres.NewExecutionLogRepository(pool)
	decisionRepo := postgres.NewDecisionLogRepository(pool)
	otpRepo := postgres.NewAccountOTPRepository(pool)

	go startBots(
		context.Background(),
		accountRepo,
		botRepo,
		positionRepo,
		decisionRepo,
		executionRepo,
	)

	// 🌐 Iniciar servidor HTTP com rotas REST
	go startHTTPServer(accountRepo, botRepo, otpRepo, pool)

	// 🛑 Aguardar sinal do SO para desligar
	waitForShutdown()
}

func startHTTPServer(
	accountRepo repository.AccountRepository,
	botRepo repository.BotRepository,
	otpRepo repository.AccountOTPRepository,
	db *pgxpool.Pool,
) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// mux := http.NewServeMux()
	router := middlewares.CORSMiddleware(
		routes.NewRouter(
			otpRepo,
			accountRepo,
			botRepo,
		),
	)

	// mux.Handle("/", router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("🌐 Servidor HTTP iniciado em http://localhost:%s", port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Erro no servidor HTTP: %v", err)
	}
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("⚠️ Encerrando servidor e bots...")
	time.Sleep(2 * time.Second)
	log.Println("✅ Encerrado com sucesso.")
}

package main


import (
	"log"
	stdhttp "net/http"

	"github.com/al4nzh/pollingservice.git/internal/config"
	"github.com/al4nzh/pollingservice.git/internal/db"
	apphttp "github.com/al4nzh/pollingservice.git/internal/http"
	"github.com/al4nzh/pollingservice.git/internal/fetcher"
	"github.com/al4nzh/pollingservice.git/internal/repository"
	"github.com/al4nzh/pollingservice.git/internal/service"

	"github.com/joho/godotenv"
	"context"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	pool, err := db.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewMarketRepository(pool)

	client := &stdhttp.Client{
		Timeout: cfg.HTTPTimeout,
	}

	csfloatFetcher := fetcher.NewCSFloatFetcher(client, cfg.CSFloatAPIKey)

	pollingService := service.NewPollingService(repo, csfloatFetcher)
	

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pollingService.StartScheduler(ctx, cfg.PollInterval)

	handler := apphttp.NewHandler(pollingService)

	mux := stdhttp.NewServeMux()
	apphttp.RegisterRoutes(mux, handler)

	server := &stdhttp.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Println("polling service running on port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
package db

import (
	"context"
	"log"
	"time"

	"github.com/al4nzh/pollingservice.git/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(cfg config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	   pool, err := pgxpool.New(ctx, cfg.DBDSN)
	   if err != nil {
		   return nil, err
	   }

	   // Проверяем соединение
	   if err := pool.Ping(ctx); err != nil {
		   return nil, err
	   }

	   log.Println("connected to postgres")

	   return pool, nil
}
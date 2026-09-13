package config

import (
	"context"
	"log/slog"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
	"neupaneanish.com.np/profile/internal/repository"
)

type Config struct {
	Pool       *pgxpool.Pool
	Client     valkey.Client
	Repository repository.Querier
	Logger     *slog.Logger
	Resolver   *net.Resolver
}

func NewConfig(ctx context.Context, env *Env, logger *slog.Logger) (*Config, error) {
	pool, poolErr := NewDatabase(ctx, env.DatabaseURL)
	if poolErr != nil {
		return nil, poolErr
	}

	client, clientErr := NewValkey(ctx, env.ValkeyURL)
	if clientErr != nil {
		return nil, clientErr
	}

	return &Config{
		Pool:       pool,
		Client:     client,
		Repository: repository.New(pool),
		Logger:     logger,
		Resolver:   net.DefaultResolver,
	}, nil
}

func (c *Config) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
	if c.Client != nil {
		c.Client.Close()
	}
}

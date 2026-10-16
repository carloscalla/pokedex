package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"pokedex-backend/internal/auth"
	"pokedex-backend/internal/catalog"
	"pokedex-backend/internal/config"
	"pokedex-backend/internal/db"
	"pokedex-backend/internal/httpapi"
	"pokedex-backend/internal/pokeapi"
	"pokedex-backend/internal/sync"
	"pokedex-backend/internal/team"
)

const jwtTTL = 24 * time.Hour

func main() {
	forceSync := flag.Bool("sync", false, "force a full catalog re-sync before starting, without touching users/team_members")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config error", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database error", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	pokeapiClient := pokeapi.NewClient(cfg.PokeAPIBaseURL)

	catalogRepo := catalog.NewRepository(pool)
	catalogService := catalog.NewService(catalogRepo)

	authRepo := auth.NewRepository(pool)
	jwtService := auth.NewJWTService(cfg.JWTSecret, jwtTTL)
	authService := auth.NewService(authRepo, jwtService)

	teamRepo := team.NewRepository(pool)
	teamService := team.NewService(teamRepo, catalogRepo)

	count, err := catalogRepo.Count(ctx)
	if err != nil {
		logger.Error("checking catalog", "error", err)
		os.Exit(1)
	}

	if count == 0 {
		logger.Info("catalog empty, running first-time sync (this will take a few minutes)")
		if err := sync.Run(ctx, pool, pokeapiClient, logger, cfg.SyncConcurrency); err != nil {
			logger.Error("initial sync failed", "error", err)
			os.Exit(1)
		}
	} else {
		logger.Info("catalog already populated, skipping sync", "count", count)
	}

	if *forceSync {
		logger.Info("--sync flag set, forcing re-sync")
		if err := sync.Run(ctx, pool, pokeapiClient, logger, cfg.SyncConcurrency); err != nil {
			logger.Error("forced sync failed", "error", err)
			os.Exit(1)
		}
	}

	syncTrigger := func() {
		if err := sync.Run(context.Background(), pool, pokeapiClient, logger, cfg.SyncConcurrency); err != nil {
			logger.Error("manual sync failed", "error", err)
		}
	}

	router := httpapi.NewRouter(httpapi.Deps{
		CatalogService: catalogService,
		AuthService:    authService,
		TeamService:    teamService,
		JWTService:     jwtService,
		AllowedOrigin:  cfg.AllowedOrigin,
		AdminToken:     cfg.AdminToken,
		SyncTrigger:    syncTrigger,
	})

	addr := ":" + cfg.Port
	logger.Info("server listening", "addr", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

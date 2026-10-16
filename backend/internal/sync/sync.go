package sync

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"pokedex-backend/internal/pokeapi"
)

const (
	maxFetchAttempts = 3
	retryBackoffBase = 200 * time.Millisecond
	progressEvery    = 100
)

// Run performs the full catalog ETL: fetch every Pokemon's {name,url} pair,
// then fetch+upsert full detail for each one with bounded concurrency. It is
// idempotent (upsert-based) and safe to re-run at any time.
func Run(ctx context.Context, pool *pgxpool.Pool, client *pokeapi.Client, logger *slog.Logger, concurrency int) error {
	names, err := client.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("listing pokemon catalog: %w", err)
	}
	total := len(names)
	logger.Info("sync starting", "total", total, "concurrency", concurrency)

	var completed atomic.Int64
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(concurrency)

	for _, nr := range names {
		nr := nr
		g.Go(func() error {
			detail, err := fetchWithRetry(gctx, client, nr.Name)
			if err != nil {
				// A single flaky entry out of ~1300 is an acceptable, visible
				// loss; aborting the whole sync over one upstream hiccup is
				// the wrong tradeoff here.
				logger.Warn("skipping pokemon after repeated failures", "name", nr.Name, "error", err)
				return nil
			}

			if err := upsert(gctx, pool, detail); err != nil {
				// A write failure, by contrast, likely signals a systemic
				// problem (bad connection, schema drift) worth failing loudly for.
				return fmt.Errorf("upserting %s: %w", nr.Name, err)
			}

			n := completed.Add(1)
			if n%progressEvery == 0 {
				logger.Info("sync progress", "done", n, "total", total)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	logger.Info("sync complete", "done", completed.Load(), "total", total)
	return nil
}

func fetchWithRetry(ctx context.Context, client *pokeapi.Client, name string) (*pokeapi.PokemonDetail, error) {
	var lastErr error
	for attempt := 1; attempt <= maxFetchAttempts; attempt++ {
		detail, err := client.GetPokemon(ctx, name)
		if err == nil {
			return detail, nil
		}
		lastErr = err

		if attempt < maxFetchAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryBackoffBase * time.Duration(attempt)):
			}
		}
	}
	return nil, lastErr
}

func upsert(ctx context.Context, pool *pgxpool.Pool, d *pokeapi.PokemonDetail) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO pokemon (id, name, sprite_url, height, weight) VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, sprite_url=EXCLUDED.sprite_url,
		                                 height=EXCLUDED.height, weight=EXCLUDED.weight`,
		d.ID, d.Name, d.Sprites.FrontDefault, d.Height, d.Weight,
	)
	if err != nil {
		return fmt.Errorf("upserting pokemon row: %w", err)
	}

	if _, err := tx.Exec(ctx, "DELETE FROM pokemon_types WHERE pokemon_id = $1", d.ID); err != nil {
		return fmt.Errorf("clearing types: %w", err)
	}
	for _, t := range d.Types {
		if _, err := tx.Exec(ctx,
			"INSERT INTO pokemon_types (pokemon_id, type, slot) VALUES ($1,$2,$3)",
			d.ID, t.Type.Name, t.Slot,
		); err != nil {
			return fmt.Errorf("inserting type: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, "DELETE FROM pokemon_stats WHERE pokemon_id = $1", d.ID); err != nil {
		return fmt.Errorf("clearing stats: %w", err)
	}
	for _, s := range d.Stats {
		if _, err := tx.Exec(ctx,
			"INSERT INTO pokemon_stats (pokemon_id, stat_name, base_stat) VALUES ($1,$2,$3)",
			d.ID, s.Stat.Name, s.BaseStat,
		); err != nil {
			return fmt.Errorf("inserting stat: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, "DELETE FROM pokemon_abilities WHERE pokemon_id = $1", d.ID); err != nil {
		return fmt.Errorf("clearing abilities: %w", err)
	}
	for _, a := range d.Abilities {
		if _, err := tx.Exec(ctx,
			"INSERT INTO pokemon_abilities (pokemon_id, ability, slot) VALUES ($1,$2,$3)",
			d.ID, a.Ability.Name, a.Slot,
		); err != nil {
			return fmt.Errorf("inserting ability: %w", err)
		}
	}

	return tx.Commit(ctx)
}

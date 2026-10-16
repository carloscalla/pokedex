package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Count returns the total number of rows in the pokemon table, regardless of
// any filter. Used at startup to decide whether the catalog needs syncing.
func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM pokemon").Scan(&count)
	return count, err
}

// filterClauses builds the shared JOIN/WHERE fragments and args used by both
// List and CountFiltered, so the two queries can never drift out of sync with
// each other (same predicate => same row set => correct totalPages math).
type filterClauses struct {
	joins string
	where string
	args  []any
}

func buildFilterClauses(p ListParams) filterClauses {
	var joins []string
	var conditions []string
	var args []any

	if p.Type != "" {
		args = append(args, p.Type)
		joins = append(joins, fmt.Sprintf("JOIN pokemon_types pt ON pt.pokemon_id = p.id AND pt.type = $%d", len(args)))
	}
	if p.Search != "" {
		args = append(args, p.Search)
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE '%%' || $%d || '%%'", len(args)))
	}

	where := "TRUE"
	if len(conditions) > 0 {
		where = strings.Join(conditions, " AND ")
	}

	return filterClauses{joins: strings.Join(joins, " "), where: where, args: args}
}

func (r *Repository) CountFiltered(ctx context.Context, p ListParams) (int, error) {
	fc := buildFilterClauses(p)
	query := fmt.Sprintf("SELECT COUNT(*) FROM pokemon p %s WHERE %s", fc.joins, fc.where)
	var count int
	err := r.pool.QueryRow(ctx, query, fc.args...).Scan(&count)
	return count, err
}

func (r *Repository) List(ctx context.Context, p ListParams) ([]PokemonSummary, error) {
	fc := buildFilterClauses(p)
	joins := fc.joins
	args := fc.args

	orderBy := "ORDER BY p.id ASC"
	switch {
	case p.SortStat == "name":
		dir := "ASC"
		if p.SortDesc {
			dir = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY p.name %s", dir)
	case ValidSortStats[p.SortStat]:
		args = append(args, p.SortStat)
		statParam := len(args)
		joins = joins + fmt.Sprintf(" LEFT JOIN pokemon_stats ps ON ps.pokemon_id = p.id AND ps.stat_name = $%d", statParam)
		dir := "ASC"
		if p.SortDesc {
			dir = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY ps.base_stat %s", dir)
	}

	args = append(args, p.PageSize, Offset(p.Page, p.PageSize))
	limitParam, offsetParam := len(args)-1, len(args)

	query := fmt.Sprintf(
		"SELECT p.id, p.name, p.sprite_url FROM pokemon p %s WHERE %s %s LIMIT $%d OFFSET $%d",
		joins, fc.where, orderBy, limitParam, offsetParam,
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []PokemonSummary
	var ids []int
	for rows.Next() {
		var s PokemonSummary
		if err := rows.Scan(&s.ID, &s.Name, &s.SpriteURL); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
		ids = append(ids, s.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	typesByID, err := r.typesForIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range summaries {
		summaries[i].Types = typesByID[summaries[i].ID]
	}

	return summaries, nil
}

func (r *Repository) typesForIDs(ctx context.Context, ids []int) (map[int][]string, error) {
	result := make(map[int][]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	rows, err := r.pool.Query(ctx,
		"SELECT pokemon_id, type FROM pokemon_types WHERE pokemon_id = ANY($1) ORDER BY pokemon_id, slot",
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var t string
		if err := rows.Scan(&id, &t); err != nil {
			return nil, err
		}
		result[id] = append(result[id], t)
	}
	return result, rows.Err()
}

func (r *Repository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pokemon WHERE id = $1)", id).Scan(&exists)
	return exists, err
}

func (r *Repository) GetByName(ctx context.Context, name string) (*PokemonDetail, error) {
	var d PokemonDetail
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, sprite_url, height, weight FROM pokemon WHERE name = $1",
		name,
	).Scan(&d.ID, &d.Name, &d.SpriteURL, &d.Height, &d.Weight)
	if err != nil {
		return nil, ErrNotFound
	}

	typesByID, err := r.typesForIDs(ctx, []int{d.ID})
	if err != nil {
		return nil, err
	}
	d.Types = typesByID[d.ID]

	stats, err := r.statsFor(ctx, d.ID)
	if err != nil {
		return nil, err
	}
	d.Stats = stats

	abilities, err := r.abilitiesFor(ctx, d.ID)
	if err != nil {
		return nil, err
	}
	d.Abilities = abilities

	return &d, nil
}

func (r *Repository) statsFor(ctx context.Context, pokemonID int) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, "SELECT stat_name, base_stat FROM pokemon_stats WHERE pokemon_id = $1", pokemonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var name string
		var value int
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		stats[name] = value
	}
	return stats, rows.Err()
}

func (r *Repository) abilitiesFor(ctx context.Context, pokemonID int) ([]string, error) {
	rows, err := r.pool.Query(ctx, "SELECT ability FROM pokemon_abilities WHERE pokemon_id = $1 ORDER BY slot", pokemonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var abilities []string
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		abilities = append(abilities, a)
	}
	return abilities, rows.Err()
}

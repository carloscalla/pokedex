package team

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Count(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM team_members WHERE user_id = $1", userID).Scan(&count)
	return count, err
}

func (r *Repository) Exists(ctx context.Context, userID int64, pokemonID int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM team_members WHERE user_id = $1 AND pokemon_id = $2)",
		userID, pokemonID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) Add(ctx context.Context, userID int64, pokemonID int, position int) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO team_members (user_id, pokemon_id, position) VALUES ($1,$2,$3)",
		userID, pokemonID, position,
	)
	return err
}

func (r *Repository) Remove(ctx context.Context, userID int64, pokemonID int) (int64, error) {
	tag, err := r.pool.Exec(ctx, "DELETE FROM team_members WHERE user_id = $1 AND pokemon_id = $2", userID, pokemonID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (r *Repository) List(ctx context.Context, userID int64) ([]Member, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tm.pokemon_id, p.name, p.sprite_url, tm.position
		FROM team_members tm
		JOIN pokemon p ON p.id = tm.pokemon_id
		WHERE tm.user_id = $1
		ORDER BY tm.position
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	var ids []int
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.PokemonID, &m.Name, &m.SpriteURL, &m.Position); err != nil {
			return nil, err
		}
		members = append(members, m)
		ids = append(ids, m.PokemonID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	typesByID, err := r.typesForIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range members {
		members[i].Types = typesByID[members[i].PokemonID]
	}

	return members, nil
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

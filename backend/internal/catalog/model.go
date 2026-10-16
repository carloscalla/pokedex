package catalog

import "errors"

var ErrNotFound = errors.New("pokemon not found")

type PokemonSummary struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	SpriteURL string   `json:"spriteUrl"`
	Types     []string `json:"types"`
}

type PokemonDetail struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	SpriteURL string         `json:"spriteUrl"`
	Height    int            `json:"height"`
	Weight    int            `json:"weight"`
	Types     []string       `json:"types"`
	Stats     map[string]int `json:"stats"`
	Abilities []string       `json:"abilities"`
}

// ListParams is the validated/whitelisted form of the list endpoint's query params.
type ListParams struct {
	Page     int
	PageSize int
	Search   string
	Type     string // "" means no type filter
	SortStat string // "" (id order) | "name" | one of the six base-stat names
	SortDesc bool
}

type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// ValidTypes are the 18 canonical Pokemon types PokeAPI uses.
var ValidTypes = map[string]bool{
	"normal": true, "fighting": true, "flying": true, "poison": true,
	"ground": true, "rock": true, "bug": true, "ghost": true, "steel": true,
	"fire": true, "water": true, "grass": true, "electric": true, "psychic": true,
	"ice": true, "dragon": true, "dark": true, "fairy": true,
}

// ValidSortStats are the six base stats sortable via the type filter/sort query param.
var ValidSortStats = map[string]bool{
	"hp": true, "attack": true, "defense": true,
	"special-attack": true, "special-defense": true, "speed": true,
}

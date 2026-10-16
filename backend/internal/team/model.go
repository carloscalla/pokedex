package team

import "errors"

const MaxTeamSize = 6

var (
	ErrPokemonNotFound = errors.New("pokemon not found")
	ErrTeamFull        = errors.New("team is full (max 6 pokemon)")
	ErrAlreadyOnTeam   = errors.New("pokemon already on team")
	ErrNotOnTeam       = errors.New("pokemon not on team")
)

type Member struct {
	PokemonID int      `json:"pokemonId"`
	Name      string   `json:"name"`
	SpriteURL string   `json:"spriteUrl"`
	Types     []string `json:"types"`
	Position  int      `json:"position"`
}

type TeamResponse struct {
	Members []Member `json:"members"`
	Size    int      `json:"size"`
	Max     int      `json:"max"`
}

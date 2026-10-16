package pokeapi

// NamedResource is PokeAPI's {name, url} pair, used in list responses.
type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type listResponse struct {
	Count   int             `json:"count"`
	Results []NamedResource `json:"results"`
}

// PokemonDetail is the subset of PokeAPI's /pokemon/{name} response we care about.
type PokemonDetail struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Abilities []struct {
		Slot    int `json:"slot"`
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
}

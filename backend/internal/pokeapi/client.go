package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// ListAll fetches the full {name,url} catalog in a single request. PokeAPI's
// total count (~1300) comfortably fits under this limit, so one page covers everything.
func (c *Client) ListAll(ctx context.Context) ([]NamedResource, error) {
	u := fmt.Sprintf("%s/pokemon?limit=10000&offset=0", c.baseURL)
	var resp listResponse
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, fmt.Errorf("listing pokemon: %w", err)
	}
	return resp.Results, nil
}

// GetPokemon fetches full detail for one Pokemon by name or id.
func (c *Client) GetPokemon(ctx context.Context, nameOrID string) (*PokemonDetail, error) {
	u := fmt.Sprintf("%s/pokemon/%s", c.baseURL, url.PathEscape(nameOrID))
	var detail PokemonDetail
	if err := c.getJSON(ctx, u, &detail); err != nil {
		return nil, fmt.Errorf("fetching pokemon %s: %w", nameOrID, err)
	}
	return &detail, nil
}

func (c *Client) getJSON(ctx context.Context, u string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d from %s", res.StatusCode, u)
	}
	return json.NewDecoder(res.Body).Decode(dest)
}

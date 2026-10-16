package catalog

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidType = errors.New("invalid type filter")
	ErrInvalidSort = errors.New("invalid sort parameter")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListQuery is the raw, unvalidated form of the list endpoint's query params.
type ListQuery struct {
	Page     string
	PageSize string
	Search   string
	Type     string
	Sort     string
}

func (s *Service) List(ctx context.Context, q ListQuery) (PaginatedResponse[PokemonSummary], error) {
	params, err := parseListParams(q)
	if err != nil {
		return PaginatedResponse[PokemonSummary]{}, err
	}

	total, err := s.repo.CountFiltered(ctx, params)
	if err != nil {
		return PaginatedResponse[PokemonSummary]{}, fmt.Errorf("counting pokemon: %w", err)
	}

	items, err := s.repo.List(ctx, params)
	if err != nil {
		return PaginatedResponse[PokemonSummary]{}, fmt.Errorf("listing pokemon: %w", err)
	}
	if items == nil {
		items = []PokemonSummary{}
	}

	return PaginatedResponse[PokemonSummary]{
		Items:      items,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: TotalPages(total, params.PageSize),
	}, nil
}

func (s *Service) GetDetail(ctx context.Context, name string) (*PokemonDetail, error) {
	return s.repo.GetByName(ctx, strings.ToLower(strings.TrimSpace(name)))
}

func parseListParams(q ListQuery) (ListParams, error) {
	page, err := strconv.Atoi(q.Page)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize := PageSize
	if q.PageSize != "" {
		if n, err := strconv.Atoi(q.PageSize); err == nil && n > 0 {
			pageSize = n
			if pageSize > 100 {
				pageSize = 100
			}
		}
	}

	typeFilter := strings.ToLower(strings.TrimSpace(q.Type))
	if typeFilter != "" && !ValidTypes[typeFilter] {
		return ListParams{}, ErrInvalidType
	}

	sortStat, desc, err := parseSort(q.Sort)
	if err != nil {
		return ListParams{}, err
	}

	return ListParams{
		Page:     page,
		PageSize: pageSize,
		Search:   strings.TrimSpace(q.Search),
		Type:     typeFilter,
		SortStat: sortStat,
		SortDesc: desc,
	}, nil
}

// parseSort accepts "attack" (ascending) or "-attack" (descending), and the
// special case "name"/"-name". Anything else is rejected rather than passed
// through, since the column/direction end up interpolated into SQL and can
// never be safely parameterized as a bind value.
func parseSort(raw string) (stat string, desc bool, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false, nil
	}
	if strings.HasPrefix(raw, "-") {
		desc = true
		raw = raw[1:]
	}
	if raw == "name" || ValidSortStats[raw] {
		return raw, desc, nil
	}
	return "", false, ErrInvalidSort
}

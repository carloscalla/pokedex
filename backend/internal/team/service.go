package team

import (
	"context"

	"pokedex-backend/internal/catalog"
)

type Service struct {
	repo        *Repository
	catalogRepo *catalog.Repository
}

func NewService(repo *Repository, catalogRepo *catalog.Repository) *Service {
	return &Service{repo: repo, catalogRepo: catalogRepo}
}

func (s *Service) List(ctx context.Context, userID int64) (TeamResponse, error) {
	members, err := s.repo.List(ctx, userID)
	if err != nil {
		return TeamResponse{}, err
	}
	if members == nil {
		members = []Member{}
	}
	return TeamResponse{Members: members, Size: len(members), Max: MaxTeamSize}, nil
}

func (s *Service) Add(ctx context.Context, userID int64, pokemonID int) (TeamResponse, error) {
	exists, err := s.catalogRepo.Exists(ctx, pokemonID)
	if err != nil {
		return TeamResponse{}, err
	}
	if !exists {
		return TeamResponse{}, ErrPokemonNotFound
	}

	count, err := s.repo.Count(ctx, userID)
	if err != nil {
		return TeamResponse{}, err
	}
	if count >= MaxTeamSize {
		return TeamResponse{}, ErrTeamFull
	}

	already, err := s.repo.Exists(ctx, userID, pokemonID)
	if err != nil {
		return TeamResponse{}, err
	}
	if already {
		return TeamResponse{}, ErrAlreadyOnTeam
	}

	if err := s.repo.Add(ctx, userID, pokemonID, count+1); err != nil {
		return TeamResponse{}, err
	}

	return s.List(ctx, userID)
}

func (s *Service) Remove(ctx context.Context, userID int64, pokemonID int) (TeamResponse, error) {
	rows, err := s.repo.Remove(ctx, userID, pokemonID)
	if err != nil {
		return TeamResponse{}, err
	}
	if rows == 0 {
		return TeamResponse{}, ErrNotOnTeam
	}
	return s.List(ctx, userID)
}

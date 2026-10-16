package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
	jwt  *JWTService
}

func NewService(repo *Repository, jwt *JWTService) *Service {
	return &Service{repo: repo, jwt: jwt}
}

func (s *Service) Signup(ctx context.Context, username, password string) (string, *User, error) {
	if len(password) < 8 {
		return "", nil, ErrWeakPassword
	}

	existing, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if existing != nil {
		return "", nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	user, err := s.repo.Create(ctx, username, string(hash))
	if err != nil {
		return "", nil, err
	}

	token, err := s.jwt.Issue(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Login deliberately returns the same ErrInvalidCredentials whether the
// username doesn't exist or the password is wrong, so the response never
// leaks which one was incorrect.
func (s *Service) Login(ctx context.Context, username, password string) (string, *User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.jwt.Issue(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

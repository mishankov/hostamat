package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type contextKey string

const (
	currentUserKey contextKey = "user"
)

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, username, password string) error {
	salt := uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+":"+salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashedPassword),
		Salt:     salt,
		Created:  time.Now(),
		Updated:  time.Now(),
		Status:   StatusActive,
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := s.repo.Get(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+":"+user.Salt))
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return user, nil
}

func (s *Service) GetCurrentUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(currentUserKey).(*User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}

	return user, nil
}

func (s *Service) SetCurrentUserToContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

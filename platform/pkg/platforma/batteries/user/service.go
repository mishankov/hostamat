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
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type authStorage interface {
	GetUserIdFromSessionId(context.Context, string) (string, error)
	CreateSessionForUser(context.Context, string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type Service struct {
	repo              repository
	authStorage       authStorage
	sessionCookieName string
}

func NewService(repo repository, authStorage authStorage, sessionCookieName string) *Service {
	return &Service{
		repo:              repo,
		authStorage:       authStorage,
		sessionCookieName: sessionCookieName,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*User, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetFromSession(ctx context.Context, sessionId string) (*User, error) {
	userId, err := s.authStorage.GetUserIdFromSessionId(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id from request: %w", err)
	}

	if userId == "" {
		return nil, nil
	}

	return s.Get(ctx, userId)
}

func (s *Service) CreateWithLoginAndPassword(ctx context.Context, username, password string) error {
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

func (s *Service) CreateSessionFromUsernameAndPassword(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+":"+user.Salt))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	session, err := s.authStorage.CreateSessionForUser(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	return session, err
}

func (s *Service) DeleteSession(ctx context.Context, sessionId string) error {
	return s.authStorage.DeleteSession(ctx, sessionId)
}

func (s *Service) CookieName() string {
	return s.sessionCookieName
}

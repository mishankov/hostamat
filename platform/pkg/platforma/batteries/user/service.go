package user

import (
	"context"
	"fmt"
	"hostamat/pkg/platforma/batteries/session"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	Get(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type sessionService interface {
	Get(ctx context.Context, id string) (*session.Session, error)
	Create(ctx context.Context, session *session.Session) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repo       repository
	sessionSvc sessionService
}

func NewService(repo repository, sessionSvc sessionService) *Service {
	return &Service{
		repo:       repo,
		sessionSvc: sessionSvc,
	}
}

func (s *Service) Get(ctx context.Context, id string) (*User, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetBySessionId(ctx context.Context, sessionId string) (*User, error) {
	session, err := s.sessionSvc.Get(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return s.repo.Get(ctx, session.User)
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

func (s *Service) Login(ctx context.Context, username, password string) (*session.Session, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+":"+user.Salt))
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	session := &session.Session{
		ID:      uuid.New().String(),
		User:    user.ID,
		Created: time.Now(),
		Expires: time.Now().Add(100 * 24 * time.Hour),
	}

	err = s.sessionSvc.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *Service) Logout(ctx context.Context, sessionId string) error {
	return s.sessionSvc.Delete(ctx, sessionId)
}

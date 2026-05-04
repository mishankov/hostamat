package feedback

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/platforma-dev/platforma/auth"
)

type repository interface {
	Create(ctx context.Context, item *Feedback) error
	GetUserContext(ctx context.Context, userID string) (*UserContext, error)
	CountByUserSince(ctx context.Context, userID string, since time.Time) (int, error)
	DeleteByUser(ctx context.Context, userID string) error
}

type Service struct {
	repo repository
	now  func() time.Time
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
}

func (s *Service) Create(ctx context.Context, message, url string) (*Feedback, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUserNotFound
	}

	message = strings.TrimSpace(message)
	if message == "" {
		return nil, errMissingMessage
	}
	if len([]rune(message)) > MaxMessageLength {
		return nil, errLongMessage
	}

	userContext, err := s.repo.GetUserContext(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	userContext.Email = strings.TrimSpace(userContext.Email)
	if userContext.Email == "" {
		return nil, errMissingEmail
	}

	url = strings.TrimSpace(url)
	if len([]rune(url)) > MaxURLLength {
		return nil, errLongURL
	}

	now := s.now()
	since := now.Add(-time.Hour)
	recentCount, err := s.repo.CountByUserSince(ctx, userContext.ID, since)
	if err != nil {
		return nil, err
	}
	if recentCount >= MaxSubmissionsPerHour {
		return nil, errRateLimited
	}

	item := &Feedback{
		ID:        uuid.NewString(),
		UserID:    userContext.ID,
		Username:  userContext.Username,
		Email:     userContext.Email,
		URL:       url,
		Message:   message,
		CreatedAt: now,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) DeleteByUser(ctx context.Context, userID string) error {
	return s.repo.DeleteByUser(ctx, userID)
}

package feedback_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"hostamat/internal/feedback"

	"github.com/platforma-dev/platforma/auth"
)

type repositoryMock struct {
	user          *feedback.UserContext
	created       *feedback.Feedback
	deletedUserID string
	recentCount   int
	err           error
}

func (m *repositoryMock) Create(ctx context.Context, item *feedback.Feedback) error {
	m.created = item
	return m.err
}

func (m *repositoryMock) GetUserContext(ctx context.Context, userID string) (*feedback.UserContext, error) {
	if m.user == nil {
		return nil, errors.New("not found")
	}
	return m.user, nil
}

func (m *repositoryMock) CountByUserSince(ctx context.Context, userID string, since time.Time) (int, error) {
	return m.recentCount, m.err
}

func (m *repositoryMock) DeleteByUser(ctx context.Context, userID string) error {
	m.deletedUserID = userID
	return m.err
}

func TestCreateStoresTrimmedMessageWithUserContext(t *testing.T) {
	t.Parallel()

	repo := &repositoryMock{
		user: &feedback.UserContext{
			ID:       "user-1",
			Username: "alice",
			Email:    "alice@example.com",
		},
	}
	service := feedback.NewService(repo)

	item, err := service.Create(testContext("user-1"), "  hello\n\nworld  ", " /dashboard/sites ")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if item.Message != "hello\n\nworld" {
		t.Fatalf("Create() message = %q", item.Message)
	}
	if item.URL != "/dashboard/sites" {
		t.Fatalf("Create() URL = %q", item.URL)
	}
	if item.UserID != "user-1" || item.Username != "alice" || item.Email != "alice@example.com" {
		t.Fatalf("Create() user context = %#v", item)
	}
	if repo.created == nil || repo.created.ID == "" || repo.created.CreatedAt.IsZero() {
		t.Fatalf("Create() did not persist complete item: %#v", repo.created)
	}
}

func TestCreateRejectsMissingUser(t *testing.T) {
	t.Parallel()

	service := feedback.NewService(&repositoryMock{})

	_, err := service.Create(context.Background(), "hello", "/dashboard")
	if !errors.Is(err, auth.ErrUserNotFound) {
		t.Fatalf("Create() error = %v, want auth.ErrUserNotFound", err)
	}
}

func TestCreateRejectsInvalidMessages(t *testing.T) {
	t.Parallel()

	service := feedback.NewService(&repositoryMock{
		user: &feedback.UserContext{ID: "user-1", Username: "alice", Email: "alice@example.com"},
	})

	for _, message := range []string{"", "   \n\t"} {
		_, err := service.Create(testContext("user-1"), message, "/dashboard")
		if err == nil {
			t.Fatalf("Create(%q) error = nil, want validation error", message)
		}
	}

	_, err := service.Create(testContext("user-1"), strings.Repeat("я", feedback.MaxMessageLength+1), "/dashboard")
	if err == nil {
		t.Fatal("Create(long message) error = nil, want validation error")
	}
}

func TestCreateRejectsLongURL(t *testing.T) {
	t.Parallel()

	service := feedback.NewService(&repositoryMock{
		user: &feedback.UserContext{ID: "user-1", Username: "alice", Email: "alice@example.com"},
	})

	_, err := service.Create(testContext("user-1"), "hello", strings.Repeat("a", feedback.MaxURLLength+1))
	if err == nil {
		t.Fatal("Create(long URL) error = nil, want validation error")
	}
}

func TestCreateRejectsMissingEmail(t *testing.T) {
	t.Parallel()

	service := feedback.NewService(&repositoryMock{
		user: &feedback.UserContext{ID: "user-1", Username: "alice", Email: " "},
	})

	_, err := service.Create(testContext("user-1"), "hello", "/dashboard")
	if err == nil {
		t.Fatal("Create() error = nil, want missing email error")
	}
}

func TestCreateRejectsMoreThanFiveSubmissionsPerHour(t *testing.T) {
	t.Parallel()

	service := feedback.NewService(&repositoryMock{
		user:        &feedback.UserContext{ID: "user-1", Username: "alice", Email: "alice@example.com"},
		recentCount: feedback.MaxSubmissionsPerHour,
	})

	_, err := service.Create(testContext("user-1"), "hello", "/dashboard")
	if err == nil {
		t.Fatal("Create() error = nil, want rate limit error")
	}
}

func TestDeleteByUserDeletesFeedback(t *testing.T) {
	t.Parallel()

	repo := &repositoryMock{}
	service := feedback.NewService(repo)

	if err := service.DeleteByUser(context.Background(), "user-1"); err != nil {
		t.Fatalf("DeleteByUser() error = %v", err)
	}

	if repo.deletedUserID != "user-1" {
		t.Fatalf("deleted user id = %q, want user-1", repo.deletedUserID)
	}
}

func testContext(userID string) context.Context {
	return context.WithValue(context.Background(), auth.UserContextKey, &auth.User{
		ID:       userID,
		Username: "ignored",
		Created:  time.Now(),
		Updated:  time.Now(),
	})
}

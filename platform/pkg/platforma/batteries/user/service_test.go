package user_test

import (
	"context"
	"errors"
	"hostamat/pkg/platforma/batteries/session"
	"hostamat/pkg/platforma/batteries/user"
	"testing"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepository{usersByUsername: map[string]*user.User{}, usersById: map[string]*user.User{}}
	sessionSvc := &mockSessionService{}
	service := user.NewService(repo, sessionSvc)

	username := "testUser"
	password := "testPassword"

	err := service.Create(ctx, username, password)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := service.Login(ctx, username, password)
	if err != nil {
		t.Fatalf("failed to login user: %v", err)
	}

	if user == nil {
		t.Fatalf("user is nil")
	}

	if sessionSvc.callCreateCount != 1 {
		t.Fatalf("session service create method not called")
	}
}

type mockSessionService struct {
	callCreateCount int
}

func (m *mockSessionService) Get(ctx context.Context, id string) (*session.Session, error) {
	return nil, nil
}

func (m *mockSessionService) Create(ctx context.Context, session *session.Session) error {
	m.callCreateCount++
	return nil
}

type mockRepository struct {
	usersByUsername map[string]*user.User
	usersById       map[string]*user.User
}

func (m *mockRepository) Get(ctx context.Context, id string) (*user.User, error) {
	if user, ok := m.usersById[id]; ok {
		return user, nil
	}

	return nil, errors.New("user not found")
}

func (m *mockRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	if user, ok := m.usersByUsername[username]; ok {
		return user, nil
	}

	return nil, errors.New("user not found")
}

func (m *mockRepository) Create(ctx context.Context, user *user.User) error {
	m.usersByUsername[user.Username] = user
	m.usersById[user.ID] = user
	return nil
}

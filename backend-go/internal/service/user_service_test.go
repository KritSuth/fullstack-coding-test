package service_test
 
import (
	"context"
	"errors"
	"testing"
 
	"github.com/yourorg/user-api/internal/model"
	"github.com/yourorg/user-api/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
 
// mockUserRepository is a test double for repository.UserRepository.
type mockUserRepository struct {
	users  map[string]*model.User
	emails map[string]*model.User
}
 
func newMockRepo() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*model.User),
		emails: make(map[string]*model.User),
	}
}
 
func (m *mockUserRepository) Create(_ context.Context, user *model.User) (*model.User, error) {
	if _, exists := m.emails[user.Email]; exists {
		return nil, errors.New("duplicate email")
	}
	user.ID = primitive.NewObjectID()
	m.users[user.ID.Hex()] = user
	m.emails[user.Email] = user
	return user, nil
}
 
func (m *mockUserRepository) FindByID(_ context.Context, id string) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}
 
func (m *mockUserRepository) FindByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := m.emails[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}
 
func (m *mockUserRepository) FindAll(_ context.Context) ([]*model.User, error) {
	result := make([]*model.User, 0, len(m.users))
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}
 
func (m *mockUserRepository) Update(_ context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	return u, nil
}
 
func (m *mockUserRepository) Delete(_ context.Context, id string) error {
	delete(m.users, id)
	return nil
}
 
func (m *mockUserRepository) Count(_ context.Context) (int64, error) {
	return int64(len(m.users)), nil
}
 
// Tests
 
func TestRegister(t *testing.T) {
	svc := service.NewUserService(newMockRepo())
 
	user, err := svc.Register(context.Background(), &model.CreateUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", user.Email)
	}
	if user.Password == "secret123" {
		t.Error("password must be hashed, not stored in plain text")
	}
}
 
func TestLogin_Success(t *testing.T) {
	svc := service.NewUserService(newMockRepo())
 
	_, _ = svc.Register(context.Background(), &model.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "pass1234",
	})
 
	token, err := svc.Login(context.Background(), &model.LoginRequest{
		Email:    "bob@example.com",
		Password: "pass1234",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Error("expected a non-empty JWT token")
	}
}
 
func TestLogin_WrongPassword(t *testing.T) {
	svc := service.NewUserService(newMockRepo())
 
	_, _ = svc.Register(context.Background(), &model.CreateUserRequest{
		Name:     "Carol",
		Email:    "carol@example.com",
		Password: "correct",
	})
 
	_, err := svc.Login(context.Background(), &model.LoginRequest{
		Email:    "carol@example.com",
		Password: "wrong",
	})
	if err == nil {
		t.Error("expected error for wrong password")
	}
}
 
func TestCount(t *testing.T) {
	svc := service.NewUserService(newMockRepo())
 
	for i := 0; i < 3; i++ {
		svc.Register(context.Background(), &model.CreateUserRequest{
			Name:     "User",
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "password",
		})
	}
 
	count, err := svc.Count(context.Background())
	if err != nil {
		t.Fatalf("count error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 users, got %d", count)
	}
}
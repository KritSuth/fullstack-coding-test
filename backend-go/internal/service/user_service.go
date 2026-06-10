package service
 
import (
	"context"
	"errors"
	"os"
	"time"
 
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourorg/user-api/internal/model"
	"github.com/yourorg/user-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/mongo"
)
 
type UserService struct {
	repo repository.UserRepository
}
 
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
 
func (s *UserService) Register(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}
	return s.repo.Create(ctx, user)
}
 
func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return generateJWT(user.ID.Hex())
}
 
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}
 
func (s *UserService) List(ctx context.Context) ([]*model.User, error) {
	return s.repo.FindAll(ctx)
}
 
func (s *UserService) Update(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	return s.repo.Update(ctx, id, req)
}
 
func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
 
func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
 
func generateJWT(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme"
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
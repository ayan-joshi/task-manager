package service

import (
	"strings"

	"github.com/google/uuid"

	"taskmanager/internal/domain"
	"taskmanager/internal/dto"
	"taskmanager/internal/repository"
	"taskmanager/internal/utils"
)

// AuthService handles registration and authentication.
type AuthService struct {
	users repository.UserRepository
	jwt   *utils.JWTManager
}

// NewAuthService constructs an AuthService.
func NewAuthService(users repository.UserRepository, jwt *utils.JWTManager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

// Signup creates a new account with a hashed password and returns a signed token.
// Email uniqueness is enforced both here (friendly error) and by a unique index
// in the database (race-condition safety net).
func (s *AuthService) Signup(req dto.SignupRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	existing, err := s.users.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, utils.ErrEmailTaken
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    email,
		Password: hash,
		Role:     domain.RoleUser,
	}
	if err := s.users.Create(user); err != nil {
		// A unique-violation here means a concurrent signup beat us to it.
		return nil, utils.ErrEmailTaken
	}

	return s.tokenResponse(user)
}

// Login verifies credentials and returns a signed token. The same error is
// returned whether the email is unknown or the password is wrong, so the
// endpoint does not reveal which emails are registered.
func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(normalizeEmail(req.Email))
	if err != nil {
		return nil, err
	}
	if user == nil || !utils.CheckPassword(user.Password, req.Password) {
		return nil, utils.ErrInvalidCreds
	}
	return s.tokenResponse(user)
}

// Me returns the user identified by id, or an unauthorized error if the user no
// longer exists (e.g. the account was deleted after the token was issued).
func (s *AuthService) Me(id uuid.UUID) (*domain.User, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrUnauthorized
	}
	return user, nil
}

func (s *AuthService) tokenResponse(user *domain.User) (*dto.AuthResponse, error) {
	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: token, User: user}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

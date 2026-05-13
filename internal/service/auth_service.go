package service

import (
	"errors"

	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/repository"
	"Sevima-AI-Content-Creator/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(refreshToken string) (*model.AuthResponse, error)
	GetProfile(userID string) (*model.UserInfo, error)
	ChangePassword(userID string, req *model.ChangePasswordRequest) error
	DeleteAccount(userID string) error
	RestoreAccount(refreshToken string) (*model.UserInfo, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check email already exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// First user is admin (atomic check for first user)
	role := "user"
	count, err := s.userRepo.Count()
	if err == nil && count == 0 {
		role = "admin"
	}
	// Note: For production, wrap Count() + Create() in a database transaction to prevent race condition
	// where two simultaneous registrations could both become admin

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     role,
		Credits:  1000,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.generateTokenResponse(user)
}

func (s *authService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return s.generateTokenResponse(user)
}

func (s *authService) RefreshToken(refreshToken string) (*model.AuthResponse, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateTokenResponse(user)
}

func (s *authService) GetProfile(userID string) (*model.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &model.UserInfo{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Credits:   user.Credits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authService) ChangePassword(userID string, req *model.ChangePasswordRequest) error {
	// Get current user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, string(hashed)); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *authService) DeleteAccount(userID string) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	_ = user // verification that user exists

	// Soft delete user
	if err := s.userRepo.Delete(userID); err != nil {
		return errors.New("failed to delete account")
	}

	return nil
}

func (s *authService) RestoreAccount(refreshToken string) (*model.UserInfo, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Find user including deleted (Unscoped)
	user, err := s.userRepo.FindByIDIncludeDeleted(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	_ = user // verification that user exists

	// Restore user (set deleted_at to null)
	if err := s.userRepo.Restore(claims.UserID); err != nil {
		return nil, errors.New("failed to restore account")
	}

	return &model.UserInfo{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Credits:   user.Credits,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authService) generateTokenResponse(user *model.User) (*model.AuthResponse, error) {
	userIDStr := user.ID.String()

	accessToken, expiresIn, err := utils.GenerateAccessToken(userIDStr, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(userIDStr, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User: model.UserInfo{
			ID:        userIDStr,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Credits:   user.Credits,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

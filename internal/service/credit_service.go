package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"Sevima-AI-Content-Creator/internal/repository"
)

type CreditService interface {
	GetCredits(userID string) (int, error)
	AddCredits(adminUserID, targetUserID string, amount int) (int, error)
	GetUserCredits(ctx context.Context, userID uuid.UUID) (int, error)
	DeductCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error
}

type creditService struct {
	userRepo repository.UserRepository
}

func NewCreditService(userRepo repository.UserRepository) CreditService {
	return &creditService{userRepo}
}

func (s *creditService) GetCredits(userID string) (int, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, errors.New("user not found")
	}
	return user.Credits, nil
}

func (s *creditService) AddCredits(adminUserID, targetUserID string, amount int) (int, error) {
	// Verify admin
	admin, err := s.userRepo.FindByID(adminUserID)
	if err != nil {
		return 0, errors.New("admin user not found")
	}
	if admin.Role != "admin" {
		return 0, errors.New("unauthorized: only admins can add credits")
	}

	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}

	// Find target user
	target, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return 0, errors.New("target user not found")
	}

	newCredits := target.Credits + amount
	if err := s.userRepo.UpdateCredits(targetUserID, newCredits); err != nil {
		return 0, errors.New("failed to update credits")
	}

	return newCredits, nil
}

func (s *creditService) GetUserCredits(ctx context.Context, userID uuid.UUID) (int, error) {
	user, err := s.userRepo.FindByID(userID.String())
	if err != nil {
		return 0, errors.New("user not found")
	}
	return user.Credits, nil
}

func (s *creditService) DeductCredits(ctx context.Context, userID uuid.UUID, amount int, reason string) error {
	user, err := s.userRepo.FindByID(userID.String())
	if err != nil {
		return errors.New("user not found")
	}
	if user.Credits < amount {
		return errors.New("insufficient credits")
	}
	return s.userRepo.UpdateCredits(userID.String(), user.Credits-amount)
}

package service

import (
	"context"
	"fmt"
	"strings"

	"crud_app/dto"
	"crud_app/repository"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks_$GOPACKAGE/mock_$GOFILE

type UserValidator interface {
	Create(ctx context.Context, user *dto.User) error
	Update(ctx context.Context, user *dto.User, id uint) error
	Delete(ctx context.Context, id uint) error
}

type userValidator struct {
	userRepo repository.UserRepo
}

func NewUserValidator(userRepo repository.UserRepo) UserValidator {
	return &userValidator{userRepo: userRepo}
}

func (v *userValidator) Create(ctx context.Context, user *dto.User) error {
	if err := v.validateNewUserData(user); err != nil {
		return err
	}

	return nil
}

func (v *userValidator) Update(ctx context.Context, user *dto.User, id uint) error {
	if err := v.validateNewUserData(user); err != nil {
		return err
	}
	if err := v.validateUserExists(ctx, id); err != nil {
		return err
	}

	return nil
}

func (v *userValidator) Delete(ctx context.Context, id uint) error {
	if err := v.validateUserExists(ctx, id); err != nil {
		return err
	}

	return nil
}

func (v *userValidator) validateNewUserData(user *dto.User) error {
	if user == nil {
		return fmt.Errorf("user object cannot be nil")
	}

	if err := v.validateName(user.Name); err != nil {
		return err
	}

	if err := v.validateAge(user.Age); err != nil {
		return err
	}

	return nil
}

func (v *userValidator) validateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("name is required")
	}

	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return fmt.Errorf("name cannot exceed 100 characters")
	}

	return nil
}

func (v *userValidator) validateAge(age int) error {
	if age <= 0 {
		return fmt.Errorf("age must be positive")
	}

	if age > 150 {
		return fmt.Errorf("age seems unrealistic")
	}

	return nil
}

func (v *userValidator) validateUserExists(ctx context.Context, id uint) error {
	exists, err := v.userRepo.Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

package service

import (
	"context"

	"crud_app/dto"
	"crud_app/repository"
)

type User interface {
	List(ctx context.Context) ([]dto.User, error)
	Create(ctx context.Context, user *dto.User) (*dto.User, error)
	Update(ctx context.Context, user *dto.User, id uint) error
	Delete(ctx context.Context, id uint) error
}

type user struct {
	userValidator UserValidator
	userRepo      repository.UserRepo
}

func NewUser(userValidator UserValidator, userRepo repository.UserRepo) User {
	return &user{
		userValidator: userValidator,
		userRepo:      userRepo,
	}
}

func (s *user) List(ctx context.Context) ([]dto.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, err
}

func (s *user) Create(ctx context.Context, user *dto.User) (*dto.User, error) {
	err := s.userValidator.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *user) Update(ctx context.Context, user *dto.User, id uint) error {
	err := s.userValidator.Update(ctx, user, id)
	if err != nil {
		return err
	}

	err = s.userRepo.Update(ctx, user, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *user) Delete(ctx context.Context, id uint) error {
	err := s.userValidator.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

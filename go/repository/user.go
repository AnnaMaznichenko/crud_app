package repository

import (
	"context"

	"gorm.io/gorm"

	"crud_app/dto"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks_$GOPACKAGE/mock_$GOFILE

const tableName string = "users"

type UserRepo interface {
	List(ctx context.Context) ([]dto.User, error)
	Create(ctx context.Context, user *dto.User) (*dto.User, error)
	Update(ctx context.Context, user *dto.User, id uint) error
	Delete(ctx context.Context, id uint) error
	Exists(ctx context.Context, id uint) (bool, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) List(ctx context.Context) ([]dto.User, error) {
	var users []dto.User

	return users, r.db.WithContext(ctx).
		Table(tableName).
		Where("deleted_at is null").
		Order("id").
		Find(&users).
		Error
}

func (r *userRepo) Create(ctx context.Context, user *dto.User) (*dto.User, error) {
	err := r.db.WithContext(ctx).
		Table(tableName).
		Create(user).
		Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *dto.User, id uint) error {
	return r.db.WithContext(ctx).
		Table(tableName).
		Where("id = ?", id).
		Updates(user).
		Error
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Table(tableName).
		Delete(&dto.User{}, id).
		Error
}

func (r *userRepo) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.Table(tableName).
		Where("id = ?", id).
		Where("deleted_at is null").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

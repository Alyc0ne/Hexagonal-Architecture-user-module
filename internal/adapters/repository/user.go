package repository

import (
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/jinzhu/gorm"
)

type (
	UserRepositoryService interface {
		FindUserByEmail(email string) (*domain.User, error)
		FindUserByResetToken(resetToken string) (*domain.ForgetPassword, error)

		CreateUser(userModel *domain.User) error
		UpdateUser(userModel *domain.User) error
		CreateForgetPassword(forgetPasswordModel *domain.ForgetPassword) error
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepositoryService {
	return &userRepository{db}
}

func (r *userRepository) FindUserByEmail(email string) (*domain.User, error) {
	user := new(domain.User)

	tx := r.db.First(user, "email = ?", email)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (r *userRepository) FindUserByResetToken(resetToken string) (*domain.ForgetPassword, error) {
	forgetPassword := new(domain.ForgetPassword)
	tx := r.db.First(forgetPassword, "reset_token = ?", resetToken)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return nil, tx.Error
	}

	return forgetPassword, nil
}

func (r *userRepository) CreateUser(userModel *domain.User) error {
	tx := r.db.Create(userModel)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *userRepository) UpdateUser(userModel *domain.User) error {
	tx := r.db.Model(userModel).Where("id = ?", userModel.ID).Update(userModel)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *userRepository) CreateForgetPassword(forgetPasswordModel *domain.ForgetPassword) error {
	tx := r.db.Create(forgetPasswordModel)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

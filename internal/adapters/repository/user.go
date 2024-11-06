package repository

import (
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/jinzhu/gorm"
)

type (
	UserRepositoryService interface {
		ReadUsers() (*[]domain.User, error)
		CountUserByEmail(email string) (int, error)
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

func (r *userRepository) ReadUsers() (*[]domain.User, error) {
	query := "select id, email, password, role from users"

	users := new([]domain.User)
	tx := r.db.Raw(query).Scan(users)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func (r *userRepository) CountUserByEmail(email string) (int, error) {
	query := "select count(*) as count from users where email = ?"

	var result domain.CountResult
	tx := r.db.Raw(query, email).Scan(&result)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return result.Count, nil
}

func (r *userRepository) FindUserByEmail(email string) (*domain.User, error) {
	query := "select id, email, password, role from users where email = ?"

	user := new(domain.User)
	tx := r.db.Raw(query, email).Scan(user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (r *userRepository) FindUserByResetToken(resetToken string) (*domain.ForgetPassword, error) {
	query := "select id, email, reset_token from forget_passwords where reset_token = ?"

	forgetPassword := new(domain.ForgetPassword)
	tx := r.db.Raw(query, resetToken).Scan(forgetPassword)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return nil, tx.Error
	}

	return forgetPassword, nil
}

func (r *userRepository) CreateUser(userModel *domain.User) error {
	tx := r.db.Create(userModel)
	if tx.RowsAffected == 0 || tx.Error != nil {
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

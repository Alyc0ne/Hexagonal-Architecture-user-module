package services_test

import (
	"errors"
	"testing"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	countUserByEmail     func(email string) (int, error)
	findUserByEmail      func(email string) (*domain.User, error)
	findUserByResetToken func(resetToken string) (*domain.ForgetPassword, error)

	createUser           func(userModel *domain.User) error
	updateUser           func(userModel *domain.User) error
	createForgetPassword func(forgetPasswordModel *domain.ForgetPassword) error
}

func (m *mockUserRepository) CountUserByEmail(email string) (int, error) {
	return m.countUserByEmail(email)
}

func (m *mockUserRepository) FindUserByEmail(email string) (*domain.User, error) {
	return m.findUserByEmail(email)
}

func (m *mockUserRepository) FindUserByResetToken(resetToken string) (*domain.ForgetPassword, error) {
	return m.findUserByResetToken(resetToken)
}

func (m *mockUserRepository) CreateUser(userModel *domain.User) error {
	return m.createUser(userModel)
}

func (m *mockUserRepository) UpdateUser(userModel *domain.User) error {
	return m.updateUser(userModel)
}

func (m *mockUserRepository) CreateForgetPassword(forgetPasswordModel *domain.ForgetPassword) error {
	return m.createForgetPassword(forgetPasswordModel)
}

func TestLoginUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	t.Run("Login Success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			findUserByEmail: func(email string) (*domain.User, error) {
				return &domain.User{ID: uuid.New().String(), Email: email, Password: string(hashedPassword)}, nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		res, err := service.LoginUser(&domain.LoginUserReq{
			Email:    "quardruple@gmail.com",
			Password: "password",
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Login Fail, Password Not Matched", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			findUserByEmail: func(email string) (*domain.User, error) {
				return &domain.User{ID: uuid.New().String(), Email: email, Password: string(hashedPassword)}, nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		req := &domain.LoginUserReq{
			Email:    "quardruple@gmail.com",
			Password: "wrong-password",
		}

		res, err := service.LoginUser(req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "password not matched", err.Error())
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Create User Success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			countUserByEmail: func(email string) (int, error) {
				return 0, nil
			},
			createUser: func(userModel *domain.User) error {
				return nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		err := service.CreateUser(&domain.CreateUserReq{
			Email:    "test@example.com",
			Password: "password",
		})
		assert.NoError(t, err)
	})

	t.Run("Create User Fail, User Already Exists", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			countUserByEmail: func(email string) (int, error) {
				return 1, nil
			},
			createUser: func(userModel *domain.User) error {
				return nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		err := service.CreateUser(&domain.CreateUserReq{
			Email:    "quardruple@gmail.com",
			Password: "password",
		})
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})
}

func TestForgetPassword(t *testing.T) {
	t.Run("Forget Password Success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			findUserByEmail: func(email string) (*domain.User, error) {
				return &domain.User{ID: uuid.New().String(), Email: email}, nil
			},
			createForgetPassword: func(forgetPassword *domain.ForgetPassword) error {
				return nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		req := &domain.ForgetPasswordReq{
			Email: "quardruple@gmail.com",
		}

		res, err := service.ForgetPassword(req)
		assert.NoError(t, err)
		assert.Equal(t, "123456789", res.ResetToken)
	})

	t.Run("Forget Password Fail, User Not Found", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			findUserByEmail: func(email string) (*domain.User, error) {
				return nil, errors.New("user not found")
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		req := &domain.ForgetPasswordReq{
			Email: "quardruple@gmail.com",
		}

		res, err := service.ForgetPassword(req)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestResetPassword(t *testing.T) {
	t.Run("Reset Password Success", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			findUserByResetToken: func(token string) (*domain.ForgetPassword, error) {
				return &domain.ForgetPassword{Email: "quardruple@gmail.com", ResetToken: "reset-token"}, nil
			},
			findUserByEmail: func(email string) (*domain.User, error) {
				return &domain.User{ID: uuid.New().String(), Email: email}, nil
			},
			updateUser: func(user *domain.User) error {
				return nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		req := &domain.ResetPasswordReq{
			ResetToken: "reset-token",
			Password:   "new_password",
		}

		err := service.ResetPassword(req)
		assert.NoError(t, err)
	})

	t.Run("Reset Password Fail, Reset Token Not Match", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			findUserByResetToken: func(token string) (*domain.ForgetPassword, error) {
				return &domain.ForgetPassword{Email: "quardruple@gmail.com", ResetToken: "reset-token"}, nil
			},
			findUserByEmail: func(email string) (*domain.User, error) {
				return &domain.User{ID: uuid.New().String(), Email: email}, nil
			},
			updateUser: func(user *domain.User) error {
				return nil
			},
		}

		service := services.NewUserUsecase("jwtSecret", mockRepo)
		req := &domain.ResetPasswordReq{
			ResetToken: "wrong-reset-token",
			Password:   "new_password",
		}

		err := service.ResetPassword(req)
		assert.Error(t, err)
		assert.Equal(t, "reset token not matched", err.Error())
	})
}

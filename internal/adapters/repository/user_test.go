package repository_test

import (
	"errors"
	"log"
	"testing"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/repository"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	_ "github.com/go-sql-driver/mysql"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	gormDB.LogMode(false)

	return gormDB, mock
}

func TestReadUsers(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	t.Run("Read Users Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow(uuid.New().String(), "quardruple@gmail.com", "password", "member")
		mock.ExpectQuery("select id, email, password, role from users").WillReturnRows(rows)

		user, err := repo.ReadUsers()

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEqual(t, 0, len(*user))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Find User By Email Not Found", func(t *testing.T) {
		mock.ExpectQuery("select id, email, password, role from users").WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.ReadUsers()

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFindUserByEmail(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	t.Run("Find User By Email Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow(uuid.New().String(), "quardruple@gmail.com", "password", "member")
		mock.ExpectQuery("select id, email, password, role from users where email = \\?").WithArgs("quardruple@gmail.com").WillReturnRows(rows)

		user, err := repo.FindUserByEmail("quardruple@gmail.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "quardruple@gmail.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Find User By Email Not Found", func(t *testing.T) {
		mock.ExpectQuery("select id, email, password, role from users where email = \\?").WithArgs("quardruple@gmail.com").WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.FindUserByEmail("quardruple@gmail.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFindUserByResetToken(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	t.Run("Find User By Reset Token Found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "reset_token"}).AddRow(uuid.New().String(), "quardruple@gmail.com", "reset-token-key")
		mock.ExpectQuery("select id, email, reset_token from forget_passwords where reset_token = \\?").WithArgs("reset-token-key").WillReturnRows(rows)

		forgetPassword, err := repo.FindUserByResetToken("reset-token-key")

		assert.NoError(t, err)
		assert.NotNil(t, forgetPassword)
		assert.Equal(t, "reset-token-key", forgetPassword.ResetToken)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Find User By Reset Token Not Found", func(t *testing.T) {
		mock.ExpectQuery("select id, email, reset_token from forget_passwords where reset_token = \\?").WithArgs("reset-token-key").WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.FindUserByResetToken("reset-token-key")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	t.Run("Create User Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateUser(&domain.User{
			ID:       uuid.New().String(),
			Email:    "quardruple@gmail.com",
			Password: "password",
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create User Failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").WillReturnError(gorm.Errors{errors.New("duplicate entry")})
		mock.ExpectRollback()

		err := repo.CreateUser(&domain.User{
			ID:       uuid.New().String(),
			Email:    "quardruple@gmail.com",
			Password: "password",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate entry")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	t.Run("Update User Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateUser(&domain.User{
			ID:       uuid.New().String(),
			Email:    "quardruple@gmail.com",
			Password: "new_password",
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update User Failure	", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users`").WillReturnError(gorm.Errors{errors.New("database error")})
		mock.ExpectRollback()

		err := repo.UpdateUser(&domain.User{
			ID:       uuid.New().String(),
			Email:    "quardruple@gmail.com",
			Password: "new_password",
		})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateForgetPassword(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	t.Run("Create ForgetPassword Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `forget_passwords`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.CreateForgetPassword(&domain.ForgetPassword{
			ID:         uuid.New().String(),
			Email:      "quardruple@gmail.com",
			ResetToken: "reset-token",
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create ForgetPassword Failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `forget_passwords`").WillReturnError(gorm.Errors{errors.New("database error")})
		mock.ExpectRollback()

		err := repo.CreateForgetPassword(&domain.ForgetPassword{
			ID:         uuid.New().String(),
			Email:      "quardruple@gmail.com",
			ResetToken: "reset-token",
		})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

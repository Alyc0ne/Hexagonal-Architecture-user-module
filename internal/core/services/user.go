package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/adapters/repository"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userPb "github.com/LordMoMA/Hexagonal-Architecture/internal/core/proto/user"
)

type (
	UserUsecaseService interface {
		ReadUsers() (*[]domain.User, error)

		LoginGrpc(req *userPb.LoginGrpcReq) (*userPb.LoginGrpcRes, error)
		LoginUser(req *domain.LoginUserReq) (*domain.LoginUserRes, error)

		CreateUserGrpc(req *userPb.CreateUserGrpcReq) (*userPb.CreateUserGrpcRes, error)
		CreateUser(req *domain.CreateUserReq) error

		ForgetPasswordGrpc(req *userPb.ForgetPasswordReq) (*userPb.ForgetPasswordRes, error)
		ForgetPassword(req *domain.ForgetPasswordReq) (*domain.ForgetPasswordRes, error)

		ResetPasswordGrpc(req *userPb.ResetPasswordReq) (*userPb.ResetPasswordRes, error)
		ResetPassword(req *domain.ResetPasswordReq) error
	}

	userUsecase struct {
		jwtSecret      string
		userRepository repository.UserRepositoryService
	}
)

func NewUserUsecase(jwtSecret string, userRepository repository.UserRepositoryService) UserUsecaseService {
	return &userUsecase{jwtSecret, userRepository}
}

func (u *userUsecase) ReadUsers() (*[]domain.User, error) {
	users, err := u.userRepository.ReadUsers()
	if err != nil {
		return nil, errors.New("users not found")
	}

	if users == nil {
		return nil, errors.New("users not found")
	}

	return users, nil
}

func (u *userUsecase) verifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) checkHasUser(email string) (*domain.User, error) {
	user, err := u.userRepository.FindUserByEmail(email)
	if err != nil {
		fmt.Println("Error finding user: ", err)
		return nil, errors.New("user not found")
	}

	if user == nil {
		fmt.Println("user is nil")
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *userUsecase) generateAccessToken(userID, jwtSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "LordMoMA-access",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour).UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (u *userUsecase) generateRefreshToken(userID, jwtSecret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "LordMoMA-refresh",
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour).UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (u *userUsecase) generateAccessTokenByLogin(JWTSecret, userId string) (*domain.AccessToken, error) {
	accessToken, err := u.generateAccessToken(userId, JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.generateRefreshToken(userId, JWTSecret)

	if err != nil {
		return nil, err
	}

	return &domain.AccessToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *userUsecase) loginLogic(req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	user, err := u.checkHasUser(req.Email)
	if err != nil {
		return nil, err
	}

	err = u.verifyPassword(user.Password, req.Password)
	if err != nil {
		fmt.Println("Error verifying password: ", err)
		return nil, errors.New("password not matched")
	}

	token, err := u.generateAccessTokenByLogin(u.jwtSecret, user.ID)
	if err != nil {
		fmt.Println("Error generating access token: ", err)
		return nil, errors.New("error generating access token")
	}

	return &domain.LoginUserRes{
		ID:           user.ID,
		Email:        user.Email,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (u *userUsecase) LoginGrpc(req *userPb.LoginGrpcReq) (*userPb.LoginGrpcRes, error) {
	loginLogicReq := &domain.LoginUserReq{
		Email:    req.Email,
		Password: req.Password,
	}
	loginLogicRes, err := u.loginLogic(loginLogicReq)
	if err != nil {
		return nil, err
	}

	return &userPb.LoginGrpcRes{
		Id:           loginLogicRes.ID,
		Email:        loginLogicRes.Email,
		AccessToken:  loginLogicRes.AccessToken,
		RefreshToken: loginLogicRes.RefreshToken,
	}, nil
}

func (u *userUsecase) LoginUser(req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	return u.loginLogic(req)
}

func (u *userUsecase) createUserLogic(req *domain.CreateUserReq) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password not hashed: %v", err)
	}

	countUser, err := u.userRepository.CountUserByEmail(req.Email)
	fmt.Println("err: ", err)
	if err != nil {
		return errors.New("user not found")
	}

	if countUser > 0 {
		return errors.New("user already exists")
	}

	user := &domain.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "member",
	}

	err = u.userRepository.CreateUser(user)
	if err != nil {
		return fmt.Errorf("user not created: %v", err)
	}

	return nil
}

func (u *userUsecase) CreateUserGrpc(req *userPb.CreateUserGrpcReq) (*userPb.CreateUserGrpcRes, error) {
	createUserReq := &domain.CreateUserReq{
		Email:    req.Email,
		Password: req.Password,
	}
	err := u.createUserLogic(createUserReq)
	if err != nil {
		return nil, err
	}

	return &userPb.CreateUserGrpcRes{
		Message: "New user created successfully",
	}, nil
}

func (u *userUsecase) CreateUser(req *domain.CreateUserReq) error {
	return u.createUserLogic(req)
}

func (u *userUsecase) forgetPasswordLogic(req *domain.ForgetPasswordReq) (*domain.ForgetPasswordRes, error) {
	_, err := u.checkHasUser(req.Email)
	if err != nil {
		return nil, err
	}

	forgetPassword := &domain.ForgetPassword{
		ID:         uuid.New().String(),
		Email:      req.Email,
		ResetToken: "123456789",
	}

	err = u.userRepository.CreateForgetPassword(forgetPassword)
	if err != nil {
		return nil, fmt.Errorf("forget password not created: %v", err)
	}

	return &domain.ForgetPasswordRes{
		ResetToken: forgetPassword.ResetToken, //จำลอง Token
	}, nil
}

func (u *userUsecase) ForgetPasswordGrpc(req *userPb.ForgetPasswordReq) (*userPb.ForgetPasswordRes, error) {
	forgetPasswordReq := &domain.ForgetPasswordReq{
		Email: req.Email,
	}
	forgetPasswordRes, err := u.forgetPasswordLogic(forgetPasswordReq)
	if err != nil {
		return nil, err
	}

	return &userPb.ForgetPasswordRes{
		ResetToken: forgetPasswordRes.ResetToken,
	}, nil
}

func (u *userUsecase) ForgetPassword(req *domain.ForgetPasswordReq) (*domain.ForgetPasswordRes, error) {
	return u.forgetPasswordLogic(req)
}

func (u *userUsecase) resetPasswordLogic(req *domain.ResetPasswordReq) error {
	forgetPassword, err := u.userRepository.FindUserByResetToken(req.ResetToken)
	if err != nil {
		return fmt.Errorf("reset token not matched: %v", err)
	}

	if forgetPassword == nil {
		return errors.New("reset token not matched")
	}

	if forgetPassword.ResetToken != req.ResetToken {
		return errors.New("reset token not matched")
	}

	user, err := u.checkHasUser(forgetPassword.Email)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password not hashed: %v", err)
	}

	user.Password = string(hashedPassword)

	err = u.userRepository.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("user not updated: %v", err)
	}

	return nil
}

func (u *userUsecase) ResetPasswordGrpc(req *userPb.ResetPasswordReq) (*userPb.ResetPasswordRes, error) {
	resetPasswordReq := &domain.ResetPasswordReq{
		ResetToken: req.ResetToken,
		Password:   req.Password,
	}
	err := u.resetPasswordLogic(resetPasswordReq)
	if err != nil {
		return nil, err
	}

	return &userPb.ResetPasswordRes{
		Message: "Password updated successfully",
	}, nil
}

func (u *userUsecase) ResetPassword(req *domain.ResetPasswordReq) error {
	return u.resetPasswordLogic(req)
}

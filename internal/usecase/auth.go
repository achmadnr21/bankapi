package usecase

import (
	"time"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type AuthUsecase struct {
	UserRepository domain.UserRepository
}

func NewAuthUsecase(userRepository domain.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		UserRepository: userRepository,
	}
}

func (u *AuthUsecase) Login(username, password string) (string, string, error) {
	// return is token, refreshToken, error
	user, err := u.UserRepository.FindByUsername(username)
	if err != nil {
		err = &utils.NotFoundError{Message: "user not found"}
		return "", "", err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		err = &utils.UnauthorizedError{Message: "user or password not match"}
		return "", "", err
	}

	// generate token
	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		err = &utils.InternalServerError{Message: "failed to generate token"}
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		err = &utils.InternalServerError{Message: "failed to generate refresh token"}
		return "", "", err
	}
	return token, refreshToken, nil

}

func (u *AuthUsecase) RefreshToken(refreshToken string) (string, string, error) {
	// lakukan validasi token
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", &utils.UnauthorizedError{Message: "invalid refresh token"}
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", "", &utils.UnauthorizedError{Message: "refresh token expired"}
	}
	// generate token jika claims valid dan tidak expired
	token, err := utils.GenerateAccessToken(claims.UserId)
	if err != nil {
		return "", "", &utils.InternalServerError{Message: "failed to generate token"}
	}
	return token, refreshToken, nil

}

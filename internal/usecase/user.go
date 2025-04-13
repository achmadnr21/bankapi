package usecase

import (
	"fmt"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type UserUsecase struct {
	UserRepository     domain.UserRepository
	UserRoleRepository domain.UserRoleRepository
}

func NewUserUsecase(userRepository domain.UserRepository, userRoleRepository domain.UserRoleRepository) *UserUsecase {
	return &UserUsecase{
		UserRepository:     userRepository,
		UserRoleRepository: userRoleRepository,
	}
}

func (u *UserUsecase) Search(askId, nik, username, email string) ([]domain.User, error) {
	// get the asker info
	asker, err := u.UserRepository.FindByID(askId)

	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := u.UserRoleRepository.FindByID(asker.RoleID)

	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add user
	if !role.CanAddUser {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}

	// search user by nik, username, email
	users, err := u.UserRepository.Search(nik, username, email)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to search user"}
	}
	return users, nil
}

func (u *UserUsecase) AddUser(proposerId string, user domain.User) (domain.User, error) {
	// get the asker info
	asker, err := u.UserRepository.FindByID(proposerId)
	if err != nil {
		return domain.User{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := u.UserRoleRepository.FindByID(asker.RoleID)
	if err != nil {
		return domain.User{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add user
	if !role.CanAddUser {
		return domain.User{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// put the proposer id to the user's employee id
	user.EmployeeID = &proposerId

	// hash the user password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, &utils.InternalServerError{Message: "failed to hash password"}
	}
	// check if user already exists
	existingUser, err := u.UserRepository.FindByNIK(user.NIK)
	if err == nil {
		return domain.User{}, &utils.ConflictError{
			Message: fmt.Sprintf("user with NIK %s already exists (username: %s)",
				user.NIK, existingUser.Username),
		}
	}
	user.Password = hashedPassword
	user, err = u.UserRepository.Save(user)
	if err != nil {
		return domain.User{}, &utils.InternalServerError{Message: "failed to save user"}
	}
	user.Password = ""
	return user, nil
}

func (u *UserUsecase) GetByNIK(askerId, nik string) (domain.User, error) {
	// get the asker info
	asker, err := u.UserRepository.FindByID(askerId)
	if err != nil {
		return domain.User{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := u.UserRoleRepository.FindByID(asker.RoleID)
	if err != nil {
		return domain.User{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add user
	if !role.CanAddUser {
		return domain.User{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}

	user, err := u.UserRepository.FindByNIK(nik)
	if err != nil {
		return domain.User{}, &utils.NotFoundError{Message: "user not found"}
	}
	user.Password = ""
	return user, nil
}

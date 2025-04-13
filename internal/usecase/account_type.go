package usecase

import (
	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type AccountTypeUsecase struct {
	repo         domain.AccountTypeRepository
	UserRepo     domain.UserRepository
	UserRoleRepo domain.UserRoleRepository
}

func NewAccountTypeUsecase(repo domain.AccountTypeRepository, userRepo domain.UserRepository, userRoleRepo domain.UserRoleRepository) *AccountTypeUsecase {
	return &AccountTypeUsecase{
		repo:         repo,
		UserRepo:     userRepo,
		UserRoleRepo: userRoleRepo,
	}
}
func (uc *AccountTypeUsecase) GetAllAccountTypes() ([]domain.AccountType, error) {
	accountTypes, _ := uc.repo.FindAll()
	// if err != nil {
	// 	return nil, &utils.NotFoundError{Message: "Account types not found"}
	// }
	return accountTypes, nil
}

func (uc *AccountTypeUsecase) GetAccountTypeByID(id int) (domain.AccountType, error) {
	accountType, err := uc.repo.FindByID(id)
	if err != nil {
		return accountType, &utils.NotFoundError{Message: "Account type not found"}
	}
	return accountType, nil
}

func (uc *AccountTypeUsecase) AddAccountType(proposerId string, accountType domain.AccountType) (domain.AccountType, error) {
	// check the proposer data
	asker, err := uc.UserRepo.FindByID(proposerId)
	if err != nil {
		return domain.AccountType{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := uc.UserRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.AccountType{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add account type
	if !role.CanAddAccountTypes {
		return domain.AccountType{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// save account type
	accountType, err = uc.repo.Save(accountType)
	if err != nil {
		return domain.AccountType{}, &utils.InternalServerError{Message: "failed to save account type"}
	}
	return accountType, nil
}

func (uc *AccountTypeUsecase) UpdateAccountType(proposerId string, accountTypeId int, accountType domain.AccountType) (domain.AccountType, error) {
	// check the proposer data
	asker, err := uc.UserRepo.FindByID(proposerId)
	if err != nil {
		return domain.AccountType{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := uc.UserRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.AccountType{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add account type
	if !role.CanAddAccountTypes {
		return domain.AccountType{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	accountType.ID = accountTypeId
	accountType, err = uc.repo.Update(accountType)
	if err != nil {
		return domain.AccountType{}, &utils.InternalServerError{Message: "failed to update account type"}
	}
	return accountType, nil
}

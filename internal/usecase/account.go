package usecase

import (
	"fmt"
	"unicode"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type AccountUsecase struct {
	repo            domain.AccountRepository
	userRepo        domain.UserRepository
	userRoleRepo    domain.UserRoleRepository
	branchRepo      domain.BranchRepository
	accountTypeRepo domain.AccountTypeRepository
	currencyRepo    domain.CurrencyRepository
}

func NewAccountUsecase(repo domain.AccountRepository, user domain.UserRepository, userRole domain.UserRoleRepository, branch domain.BranchRepository, accountType domain.AccountTypeRepository, currency domain.CurrencyRepository) *AccountUsecase {
	return &AccountUsecase{
		repo:            repo,
		userRepo:        user,
		userRoleRepo:    userRole,
		branchRepo:      branch,
		accountTypeRepo: accountType,
		currencyRepo:    currency,
	}
}
func (uc *AccountUsecase) GetAllAccounts(proposerId string) ([]domain.Account, error) {
	// check the proposer data
	asker, err := uc.userRepo.FindByID(proposerId)
	if err != nil {
		return nil, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := uc.userRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add account type
	if !role.CanAddUser {
		return nil, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// get all accounts
	accounts, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (uc *AccountUsecase) GetAccountByAccountNumber(proposerId string, accountNumber string) (domain.Account, error) {
	// check the proposer data
	asker, err := uc.userRepo.FindByID(proposerId)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := uc.userRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Account{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add account type
	if !role.CanAddUser {
		return domain.Account{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// get account by account number
	account, err := uc.repo.FindByAccountNumber(accountNumber)
	if err != nil {
		return account, &utils.NotFoundError{Message: "account not found"}
	}
	return account, nil
}

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return s != "" // to avoid empty string being considered numeric
}

func (uc *AccountUsecase) AddAccount(proposerId string, branchId int, accountTypeId int, nik string, account domain.Account) (domain.Account, error) {
	// =============> PROPOSER CHECKING UNIT
	// check the proposer data
	asker, err := uc.userRepo.FindByID(proposerId)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := uc.userRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Account{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add account type
	if !role.CanAddUser {
		return domain.Account{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}

	// =============> USER WHO ADD ACCOUNT CHECKING UNIT
	// get the existing user by nik
	user, err := uc.userRepo.FindByNIK(nik)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "user not found"}
	}

	// =============> BRANCH CHECKING UNIT
	// check if the branch id is valid
	branch, err := uc.branchRepo.FindByID(branchId)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "branch not found"}
	}
	// =============> ACCOUNT TYPE CHECKING UNIT
	// check if the account type id is valid
	accountType, err := uc.accountTypeRepo.FindByID(accountTypeId)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "account type not found"}
	}
	// =============> CURRENCY CHECKING UNIT
	// Currency ID default is IDR
	if account.CurrencyID == "" {
		account.CurrencyID = "IDR"
	}
	// check if the currency id is valid
	currency, err := uc.currencyRepo.FindByID(account.CurrencyID)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "currency not found"}
	}
	// checking flow
	// 1. Put the user id to account
	account.UserID = user.ID
	// 2. put asker id to account
	account.EmployeeID = asker.ID
	// 3. put the branch id to account
	account.BranchID = branch.ID
	// 4. put the account type id to account
	account.AccountTypeID = accountType.ID
	// 5. put the currency id to account
	account.CurrencyID = currency.ID
	// 6. put the balance to account
	account.Balance = 0.00
	// =============> PIN HASH CHECKING UNIT
	// 1. Check if the pin hash is not empty
	if account.PinHash == "" {
		return domain.Account{}, &utils.BadRequestError{Message: "pin hash is required"}
	}
	// 2. Check if the pin hash must be 6 digits and numeric
	if len(account.PinHash) != 6 || !isNumeric(account.PinHash) {
		return domain.Account{}, &utils.BadRequestError{Message: "pin hash must be 6 digits and numeric"}
	}
	// 3. Perform hashing on the pin hash
	account.PinHash, err = utils.HashPassword(account.PinHash)
	if err != nil {
		return domain.Account{}, &utils.InternalServerError{Message: "failed to hash pin"}
	}

	fmt.Printf("2.1. Account Usecase:\n%+v\n", account)
	account, err = uc.repo.Save(account)
	fmt.Printf("2.2. Account Usecase:\n%+v\n", account)
	if err != nil {
		return domain.Account{}, &utils.InternalServerError{Message: "failed to save account"}
	}
	return account, nil
}

/*
	type Account struct {
		ID             string    `json:"id" db:"id"`
		UserID         string    `json:"user_id" db:"user_id"`
		EmployeeID     string    `json:"employee_id" db:"employee_id"`
		BranchID       int       `json:"branch_id" db:"branch_id"`
		AccountTypeID  int       `json:"account_type_id" db:"account_type_id"`
		SequenceNumber int       `json:"sequence_number" db:"sequence_number"`
		AccountNumber  string    `json:"account_number" db:"account_number"`
		PinHash        string    `json:"pin_hash" db:"pin_hash"`
		CurrencyID     string    `json:"currency_id" db:"currency_id"`
		Balance        float64   `json:"balance" db:"balance"`
		CreatedAt      time.Time `json:"created_at" db:"created_at
	}

yang dapat di update adalah:
1. PinHash
*/
func (uc *AccountUsecase) UpdateAccount(proposerId string, accountNumber string, account domain.Account) (domain.Account, error) {
	// check the proposer data
	asker, err := uc.userRepo.FindByID(proposerId)
	if err != nil {
		return domain.Account{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := uc.userRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Account{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// get account by account number
	existedAccount, err := uc.repo.FindByAccountNumber(accountNumber)
	if err != nil {
		return existedAccount, &utils.NotFoundError{Message: "account not found"}
	}
	// jika user tidak bisa adduser atau asker.ID tidak sama dengan proposerId, maka return unauthorized
	fmt.Printf("%t && %t\n", !role.CanAddUser, asker.ID != existedAccount.UserID)
	if !role.CanAddUser && asker.ID != existedAccount.UserID {

		return domain.Account{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}

	// START THE FLOW
	// jika CanAddUser, maka employee id diubah menjadi asker id
	if role.CanAddUser {
		existedAccount.EmployeeID = asker.ID
	}
	// check jika account.PinHash kosong, maka tidak perlu ubah dari existedAccount
	if account.PinHash != "" {
		// cek jika pin hash tidak 6 digit dan numeric
		if len(account.PinHash) != 6 || !isNumeric(account.PinHash) {
			return domain.Account{}, &utils.BadRequestError{Message: "pin hash must be 6 digits and numeric"}
		}
		// jika pin hash tidak kosong, maka ubah pin hash dengan hashed account.PinHash
		existedAccount.PinHash, err = utils.HashPassword(account.PinHash)
		if err != nil {
			return domain.Account{}, &utils.InternalServerError{Message: "failed to hash pin"}
		}
	}
	account, err = uc.repo.Update(existedAccount)
	if err != nil {
		return domain.Account{}, &utils.InternalServerError{Message: "failed to update account"}
	}
	// return the updated account
	return account, nil
}

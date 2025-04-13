package usecase

import (
	"fmt"
	"strings"

	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type TransactionFeeUsecase struct {
	repo         domain.TransactionFeeRepository
	UserRepo     domain.UserRepository
	UserRoleRepo domain.UserRoleRepository
}

func NewTransactionFeeUsecase(repo domain.TransactionFeeRepository, userRepo domain.UserRepository, userRoleRepo domain.UserRoleRepository) *TransactionFeeUsecase {
	return &TransactionFeeUsecase{
		repo:         repo,
		UserRepo:     userRepo,
		UserRoleRepo: userRoleRepo,
	}
}
func (uc *TransactionFeeUsecase) GetAllTransactionFees() ([]domain.TransactionFee, error) {
	transactionFees, err := uc.repo.FindAll()
	if err != nil {
		return nil, &utils.NotFoundError{Message: "Transaction fees not found"}
	}
	return transactionFees, nil
}
func (uc *TransactionFeeUsecase) GetTransactionFeeByID(id string) (domain.TransactionFee, error) {
	transactionFee, err := uc.repo.FindByID(id)
	if err != nil {
		return domain.TransactionFee{}, &utils.NotFoundError{Message: "Transaction fee not found"}
	}
	return transactionFee, nil
}

func (uc *TransactionFeeUsecase) AddTransactionFee(proposerId string, transactionFee domain.TransactionFee) (domain.TransactionFee, error) {
	// Check if the user is allowed to add transaction fees
	user, err := uc.UserRepo.FindByID(proposerId)
	if err != nil {
		return domain.TransactionFee{}, &utils.NotFoundError{Message: "User not found"}
	}
	userRole, err := uc.UserRoleRepo.FindByID(user.RoleID)
	if err != nil {
		return domain.TransactionFee{}, &utils.NotFoundError{Message: "User role not found"}
	}
	if !userRole.CanAddAccountTypes {
		return domain.TransactionFee{}, &utils.UnauthorizedError{Message: "User is not authorized to add transaction fees"}
	}
	// Start adding the transaction fee

	// check that the id is not empty and the length must be 3
	if transactionFee.ID == "" || len(transactionFee.ID) != 3 {
		return domain.TransactionFee{}, &utils.BadRequestError{Message: "Transaction fee ID must be 3 characters long"}
	}
	// check if the transaction fee already exists
	existingTransactionFee, err := uc.repo.FindByID(transactionFee.ID)
	if err == nil {
		return domain.TransactionFee{}, &utils.ConflictError{Message: fmt.Sprintf("Transaction fee with ID %s already exists", existingTransactionFee.ID)}
	}
	// check that the fee value is not negative or zero
	if transactionFee.Fee <= 0 {
		return domain.TransactionFee{}, &utils.BadRequestError{Message: "Transaction fee must be greater than zero"}
	}
	// check that the name is not empty
	if transactionFee.Name == "" && len(transactionFee.Name) < 3 {
		return domain.TransactionFee{}, &utils.BadRequestError{Message: "Transaction fee name cannot be empty"}
	}

	// before saving make sure to perform to upper for the ID
	transactionFee.ID = strings.ToUpper(transactionFee.ID)

	// perform the insert (save)
	transactionFee, err = uc.repo.Save(transactionFee)
	if err != nil {
		return domain.TransactionFee{}, &utils.InternalServerError{Message: "Failed to add transaction fee"}
	}
	return transactionFee, nil
}
func (uc *TransactionFeeUsecase) UpdateTransactionFee(proposerId string, transactionFee domain.TransactionFee) (domain.TransactionFee, error) {
	// Check if the user is allowed to update the transaction fee
	user, err := uc.UserRepo.FindByID(proposerId)
	if err != nil {
		return domain.TransactionFee{}, &utils.NotFoundError{Message: "User not found"}
	}
	userRole, err := uc.UserRoleRepo.FindByID(user.RoleID)
	if err != nil {
		return domain.TransactionFee{}, &utils.NotFoundError{Message: "User role not found"}
	}
	if !userRole.CanAddAccountTypes {
		return domain.TransactionFee{}, &utils.UnauthorizedError{Message: "User is not authorized to update transaction fees"}
	}
	// Start updating the transaction fee
	// get old transaction fee
	oldTransactionFee, err := uc.repo.FindByID(transactionFee.ID)
	if err != nil {
		return domain.TransactionFee{}, &utils.NotFoundError{Message: "Transaction fee not found"}
	}
	// update the old transaction fee by checking and inserting the data from new transaction fee
	// == first check : if the name is empty, then use the old name
	if transactionFee.Name != "" && len(transactionFee.Name) > 3 {
		oldTransactionFee.Name = transactionFee.Name
	}
	// == second check : if the fee is empty, then use the old fee
	if transactionFee.Fee != 0 {
		oldTransactionFee.Fee = transactionFee.Fee
	}
	// update the transaction fee
	updatedTransactionFee, err := uc.repo.Update(oldTransactionFee)
	if err != nil {
		return domain.TransactionFee{}, &utils.InternalServerError{Message: "Failed to update transaction fee"}
	}
	return updatedTransactionFee, nil
}
func (uc *TransactionFeeUsecase) DeleteTransactionFee(proposerId, id string) error {
	// Check if the user is allowed to delete the transaction fee
	user, err := uc.UserRepo.FindByID(proposerId)
	if err != nil {
		return &utils.NotFoundError{Message: "User not found"}
	}
	userRole, err := uc.UserRoleRepo.FindByID(user.RoleID)
	if err != nil {
		return &utils.NotFoundError{Message: "User role not found"}
	}
	if !userRole.CanAddAccountTypes {
		return &utils.UnauthorizedError{Message: "User is not authorized to delete transaction fees"}
	}
	// Start deleting the transaction fee
	err = uc.repo.Delete(id)
	if err != nil {
		return &utils.InternalServerError{Message: err.Error()}
	}
	return nil
}

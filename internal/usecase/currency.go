/*
package domain

import "time"

type Currency struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	RateInIDR float64   `json:"rate_in_idr" db:"rate_in_idr"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CurrencyRepository interface {
	FindAll() ([]Currency, error)
	FindByID(id string) (Currency, error)
	FindByName(name string) (Currency, error)
	Save(currency Currency) (Currency, error)
	Update(currency Currency) (Currency, error)
	Delete(id string) error
}
*/

package usecase

import (
	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type CurrencyUsecase struct {
	repo         domain.CurrencyRepository
	userRepo     domain.UserRepository
	userRoleRepo domain.UserRoleRepository
}

func NewCurrencyUsecase(repo domain.CurrencyRepository, userRepo domain.UserRepository, userRoleRepo domain.UserRoleRepository) *CurrencyUsecase {
	return &CurrencyUsecase{
		repo:         repo,
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (uc *CurrencyUsecase) GetAllCurrencies() ([]domain.Currency, error) {
	currencies, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

func (uc *CurrencyUsecase) GetCurrencyByID(id string) (domain.Currency, error) {
	// cek bahwa id harus ada. Panjang ID harus antara 1<= len(id) <= 3
	if id == "" || len(id) < 1 || len(id) > 3 {
		return domain.Currency{}, &utils.BadRequestError{
			Message: "currency id must be between 1 and 3 characters long"}
	}
	currency, err := uc.repo.FindByID(id)
	if err != nil {
		return domain.Currency{}, err
	}
	return currency, nil
}

func (uc *CurrencyUsecase) AddCurrency(proposerId string, currency domain.Currency) (domain.Currency, error) {
	// Get proposer data
	asker, err := uc.userRepo.FindByID(proposerId)
	if err != nil {
		return domain.Currency{}, &utils.UnauthorizedError{
			Message: "proposer not found"}
	}
	// Check if proposer is a currency admin
	role, err := uc.userRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Currency{}, &utils.UnauthorizedError{
			Message: "proposer not found"}
	}
	// check if CanAddCurrencies
	if !role.CanAddCurrencies {
		return domain.Currency{}, &utils.UnauthorizedError{
			Message: "proposer not authorized"}
	}
	// Check if currency already exists
	existingCurrency, err := uc.repo.FindByID(currency.ID)
	if err == nil {
		return domain.Currency{}, &utils.ConflictError{
			Message: "currency already exists with id " + existingCurrency.ID}
	}
	// Cek jika currency id kosong dan panjang string lebih dari 3 digit, maka akan return error
	if currency.ID == "" || len(currency.ID) < 3 {
		return domain.Currency{}, &utils.BadRequestError{
			Message: "currency id must be at least 3 characters long"}
	}
	// Cek jika currency name kosong dan panjang string lebih dari 3 digit, maka akan return error
	if currency.Name == "" || len(currency.Name) < 4 {
		return domain.Currency{}, &utils.BadRequestError{
			Message: "currency name must be at least 4 characters long"}
	}
	// Save currency
	currency, err = uc.repo.Save(currency)
	if err != nil {
		return domain.Currency{}, err
	}
	return currency, nil
}

func (uc *CurrencyUsecase) UpdateCurrency(proposerId string, id string, currency domain.Currency) (domain.Currency, error) {
	// Get proposer data
	asker, err := uc.userRepo.FindByID(proposerId)
	if err != nil {
		return domain.Currency{}, &utils.UnauthorizedError{
			Message: "proposer not found"}
	}
	// Check if proposer is a currency admin
	role, err := uc.userRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Currency{}, &utils.UnauthorizedError{
			Message: "proposer not found"}
	}
	// check if CanUpdateCurrencies
	if !role.CanAddCurrencies {
		return domain.Currency{}, &utils.UnauthorizedError{
			Message: "proposer not authorized"}
	}
	// Find currency by ID
	existingCurrency, err := uc.repo.FindByID(id)
	if err != nil {
		return domain.Currency{}, err
	}
	// Update currency
	currency.ID = existingCurrency.ID

	// Cek jika currency name kosong, maka masukkan dari existing currency
	if currency.Name == "" {
		currency.Name = existingCurrency.Name
	}

	// cek jika currency rate_in_idr kosong, maka masukkan dari existing currency
	if currency.RateInIDR == 0 {
		currency.RateInIDR = existingCurrency.RateInIDR
	}

	currency, err = uc.repo.Update(currency)
	if err != nil {
		return domain.Currency{}, err
	}
	return currency, nil
}

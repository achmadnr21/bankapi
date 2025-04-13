package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type AccountTypeRepository struct {
	db *sql.DB
}

func NewAccountTypeRepository(DB *sql.DB) *AccountTypeRepository {
	return &AccountTypeRepository{db: DB}
}
func (r *AccountTypeRepository) FindAll() ([]domain.AccountType, error) {
	rows, err := r.db.Query("SELECT * FROM account.account_types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accountTypes []domain.AccountType
	for rows.Next() {
		var accountType domain.AccountType
		err := rows.Scan(&accountType.ID, &accountType.Name, &accountType.CanEToll, &accountType.CanIndomart, &accountType.CanSteam, &accountType.CreatedAt)
		if err != nil {
			return nil, err
		}
		accountTypes = append(accountTypes, accountType)
	}

	return accountTypes, nil
}
func (r *AccountTypeRepository) FindByID(id int) (domain.AccountType, error) {
	var accountType domain.AccountType
	err := r.db.QueryRow("SELECT * FROM account.account_types WHERE id = $1", id).Scan(&accountType.ID, &accountType.Name, &accountType.CanEToll, &accountType.CanIndomart, &accountType.CanSteam, &accountType.CreatedAt)
	if err != nil {
		return domain.AccountType{}, err
	}
	return accountType, nil
}

func (r *AccountTypeRepository) FindByName(name string) (domain.AccountType, error) {
	var accountType domain.AccountType
	err := r.db.QueryRow("SELECT * FROM account.account_types WHERE lower(name) LIKE lower($1)", "%"+name+"%").Scan(
		&accountType.ID,
		&accountType.Name,
		&accountType.CanEToll,
		&accountType.CanIndomart,
		&accountType.CanSteam,
		&accountType.CreatedAt)
	if err != nil {
		return domain.AccountType{}, err
	}
	return accountType, nil
}

func (r *AccountTypeRepository) Save(accountType domain.AccountType) (domain.AccountType, error) {
	fmt.Println("accountType: ", accountType)
	err := r.db.QueryRow(
		"INSERT INTO account.account_types (name, can_etoll, can_indomart, can_steam) VALUES ($1, $2, $3, $4) RETURNING id",
		accountType.Name,
		accountType.CanEToll,
		accountType.CanIndomart,
		accountType.CanSteam).Scan(&accountType.ID)
	if err != nil {
		fmt.Println("error: ", err)
		return domain.AccountType{}, err
	}
	return accountType, nil
}
func (r *AccountTypeRepository) Update(accountType domain.AccountType) (domain.AccountType, error) {
	err := r.db.QueryRow(
		"UPDATE account.account_types SET name = $1, can_etoll = $2, can_indomart = $3, can_steam = $4 WHERE id = $5 RETURNING id",
		accountType.Name,
		accountType.CanEToll,
		accountType.CanIndomart,
		accountType.CanSteam,
		accountType.ID).Scan(&accountType.ID)
	if err != nil {
		return domain.AccountType{}, err
	}
	return accountType, nil
}
func (r *AccountTypeRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM account.account_types WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

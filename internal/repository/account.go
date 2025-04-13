package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(DB *sql.DB) *AccountRepository {
	return &AccountRepository{db: DB}
}
func (r *AccountRepository) FindAll() ([]domain.Account, error) {
	rows, err := r.db.Query("SELECT * FROM account.accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		err := rows.Scan(&account.ID, &account.UserID, &account.EmployeeID, &account.BranchID, &account.AccountTypeID, &account.SequenceNumber, &account.AccountNumber, &account.PinHash, &account.CurrencyID, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
func (r *AccountRepository) FindByID(id string) (domain.Account, error) {
	row := r.db.QueryRow("SELECT * FROM account.accounts WHERE id = $1", id)
	var account domain.Account
	err := row.Scan(&account.ID, &account.UserID, &account.EmployeeID, &account.BranchID, &account.AccountTypeID, &account.SequenceNumber, &account.AccountNumber, &account.PinHash, &account.CurrencyID, &account.Balance, &account.CreatedAt)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}
func (r *AccountRepository) FindByUserID(userID string) ([]domain.Account, error) {
	rows, err := r.db.Query("SELECT * FROM account.accounts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		err := rows.Scan(&account.ID, &account.UserID, &account.EmployeeID, &account.BranchID, &account.AccountTypeID, &account.SequenceNumber, &account.AccountNumber, &account.PinHash, &account.CurrencyID, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
func (r *AccountRepository) FindByAccountNumber(accountNumber string) (domain.Account, error) {
	row := r.db.QueryRow("SELECT * FROM account.accounts WHERE account_number = $1", accountNumber)
	var account domain.Account
	err := row.Scan(&account.ID, &account.UserID, &account.EmployeeID, &account.BranchID, &account.AccountTypeID, &account.SequenceNumber, &account.AccountNumber, &account.PinHash, &account.CurrencyID, &account.Balance, &account.CreatedAt)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

/*
-- TABEL SEQUENCE PER KOMB CABANG + TIPE
CREATE TABLE account.account_sequences (

	branch_id INTEGER NOT NULL,
	account_type_id INTEGER NOT NULL,
	last_sequence INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (branch_id, account_type_id),
	FOREIGN KEY (branch_id) REFERENCES account.branches(id),
	FOREIGN KEY (account_type_id) REFERENCES account.account_types(id)

);
we will use this  to get the next sequence number for the account
*/
func generateAccountNumber(branchID int, accountTypeID int, sequenceNumber int) string {
	// Generate account number based on sequence number and branch ID
	// This is a placeholder implementation. You can customize it as needed.
	var accountNumber string = fmt.Sprintf("%04d%04d%012d", branchID, accountTypeID, sequenceNumber)
	// Ensure the account number is unique
	return accountNumber
}
func (r *AccountRepository) getNextSequenceNumber(branchID int, accountTypeID int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var seq int
	err = tx.QueryRow("SELECT last_sequence FROM account.account_sequences WHERE branch_id = $1 AND account_type_id = $2 FOR UPDATE", branchID, accountTypeID).Scan(&seq)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no row exists, insert a new row with the initial sequence number
			_, err = tx.Exec("INSERT INTO account.account_sequences (branch_id, account_type_id, last_sequence) VALUES ($1, $2, 0)", branchID, accountTypeID)
			if err != nil {
				return 0, err
			}
			seq = 0
		} else {
			return 0, err
		}
	}
	// Increment the sequence number
	seq++
	_, err = tx.Exec("UPDATE account.account_sequences SET last_sequence = $1 WHERE branch_id = $2 AND account_type_id = $3", seq, branchID, accountTypeID)
	if err != nil {
		return 0, err
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	// Return the sequence number and account number
	return seq, nil
}

/*
the save function must be using begin and select for update to ensure that the sequence number is latest on its pair
branch and account type id.
*/
func (r *AccountRepository) Save(account domain.Account) (domain.Account, error) {
	// Get the next sequence number for the account
	seq, err := r.getNextSequenceNumber(account.BranchID, account.AccountTypeID)
	if err != nil {
		return domain.Account{}, err
	}
	// Generate the account number based on the sequence number
	account.AccountNumber = generateAccountNumber(account.BranchID, account.AccountTypeID, seq)
	account.SequenceNumber = seq

	query := "INSERT INTO account.accounts (user_id, employee_id, branch_id, account_type_id, sequence_number, account_number, pin_hash, currency_id, balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	err = r.db.QueryRow(query,
		account.UserID,
		account.EmployeeID,
		account.BranchID,
		account.AccountTypeID,
		account.SequenceNumber,
		account.AccountNumber,
		account.PinHash,
		account.CurrencyID,
		account.Balance).Scan(&account.ID)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

func (r *AccountRepository) Update(account domain.Account) (domain.Account, error) {
	query := "UPDATE account.accounts SET user_id = $1, employee_id = $2, branch_id = $3, account_type_id = $4, sequence_number = $5, account_number = $6, pin_hash = $7, currency_id = $8, balance = $9 WHERE id = $10"
	_, err := r.db.Exec(query,
		account.UserID,
		account.EmployeeID,
		account.BranchID,
		account.AccountTypeID,
		account.SequenceNumber,
		account.AccountNumber,
		account.PinHash,
		account.CurrencyID,
		account.Balance,
		account.ID)
	if err != nil {
		return domain.Account{}, err
	}
	return account, nil
}
func (r *AccountRepository) Delete(id string) error {
	query := "DELETE FROM account.accounts WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

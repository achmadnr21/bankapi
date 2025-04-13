/*
create table transfer.transaction_fees(
	id char(3) unique primary key,
	name varchar(255) not null,
	fee decimal(15,2) default 0,
	created_at TIMESTAMP DEFAULT NOW()
);
insert into transfer.transaction_fees(id, name, fee)
values
('INT', 'Same Bank TF', 0.0),
('BIF', 'BI-Fast', 2500.0);


package domain

type TransactionFee struct {
	ID        string  `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Fee       float64 `json:"fee" db:"fee"`
	CreatedAt string  `json:"created_at" db:"created_at"`
}

type TransactionFeeRepository interface {
	FindAll() ([]TransactionFee, error)
	FindByID(id string) (TransactionFee, error)
	FindByName(name string) (TransactionFee, error)
	Save(transactionFee TransactionFee) (TransactionFee, error)
	Update(transactionFee TransactionFee) (TransactionFee, error)
	Delete(id string) error
}
*/

package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type TransactionFeeRepository struct {
	db *sql.DB
}

func NewTransactionFeeRepository(db *sql.DB) *TransactionFeeRepository {
	return &TransactionFeeRepository{
		db: db,
	}
}

func (r *TransactionFeeRepository) FindAll() ([]domain.TransactionFee, error) {
	rows, err := r.db.Query("SELECT id, name, fee, created_at FROM transfer.transaction_fees")
	// rows, err := r.db.Query("SELECT * FROM transfer.transaction_fees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionFees []domain.TransactionFee
	for rows.Next() {
		var transactionFee domain.TransactionFee
		err := rows.Scan(&transactionFee.ID, &transactionFee.Name, &transactionFee.Fee, &transactionFee.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactionFees = append(transactionFees, transactionFee)
	}

	return transactionFees, nil
}

func (r *TransactionFeeRepository) FindByID(id string) (domain.TransactionFee, error) {
	var transactionFee domain.TransactionFee
	err := r.db.QueryRow("SELECT id, name, fee, created_at FROM transfer.transaction_fees WHERE id = $1", id).Scan(&transactionFee.ID, &transactionFee.Name, &transactionFee.Fee, &transactionFee.CreatedAt)
	if err != nil {
		return domain.TransactionFee{}, err
	}
	return transactionFee, nil
}
func (r *TransactionFeeRepository) FindByName(name string) (domain.TransactionFee, error) {
	var transactionFee domain.TransactionFee
	err := r.db.QueryRow("SELECT id, name, fee, created_at FROM transfer.transaction_fees WHERE lower(name) LIKE lower($1)", "%"+name+"%").Scan(
		&transactionFee.ID,
		&transactionFee.Name,
		&transactionFee.Fee,
		&transactionFee.CreatedAt)
	if err != nil {
		return domain.TransactionFee{}, err
	}
	return transactionFee, nil
}
func (r *TransactionFeeRepository) Save(transactionFee domain.TransactionFee) (domain.TransactionFee, error) {
	_, err := r.db.Exec("INSERT INTO transfer.transaction_fees (id, name, fee) VALUES ($1, $2, $3)",
		transactionFee.ID,
		transactionFee.Name,
		transactionFee.Fee)
	if err != nil {
		return domain.TransactionFee{}, err
	}
	return transactionFee, nil
}
func (r *TransactionFeeRepository) Update(transactionFee domain.TransactionFee) (domain.TransactionFee, error) {
	_, err := r.db.Exec("UPDATE transfer.transaction_fees SET name = $1, fee = $2 WHERE id = $3",
		transactionFee.Name,
		transactionFee.Fee,
		transactionFee.ID)
	if err != nil {
		return domain.TransactionFee{}, err
	}
	return transactionFee, nil
}
func (r *TransactionFeeRepository) Delete(id string) error {
	// debug print the id
	fmt.Println("Deleting transaction fee with id:", id)
	result, err := r.db.Exec("DELETE FROM transfer.transaction_fees WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no transaction fee found with id: %s", id)
	}
	return nil
}

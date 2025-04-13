/*


create table account.currencies(
	id 			char(3) unique not null primary key,
	name 		varchar(255) not null,
	rate_in_idr DECIMAL(15, 2) not null,
	created_at  TIMESTAMP DEFAULT NOW()
);

 package domain

 import "time"
 type Currency struct {
	 ID         string    `json:"id" db:"id"`
	 Name       string    `json:"name" db:"name"`
	 RateInIDR  float64   `json:"rate_in_idr" db:"rate_in_idr"`
	 CreatedAt  time.Time `json:"created_at" db:"created_at"`
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

package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) domain.CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) FindAll() ([]domain.Currency, error) {
	query := "SELECT id, name, rate_in_idr, created_at FROM account.currencies"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var currencies []domain.Currency
	for rows.Next() {
		var currency domain.Currency
		if err := rows.Scan(&currency.ID, &currency.Name, &currency.RateInIDR, &currency.CreatedAt); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return currencies, nil
}
func (r *CurrencyRepository) FindByID(id string) (domain.Currency, error) {
	query := "SELECT id, name, rate_in_idr, created_at FROM account.currencies WHERE id = $1"
	row := r.db.QueryRow(query, id)
	var currency domain.Currency
	if err := row.Scan(&currency.ID, &currency.Name, &currency.RateInIDR, &currency.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Currency{}, fmt.Errorf("currency with id %s not found", id)
		}
		return domain.Currency{}, err
	}
	return currency, nil
}
func (r *CurrencyRepository) FindByName(name string) (domain.Currency, error) {
	query := "SELECT id, name, rate_in_idr, created_at FROM account.currencies WHERE name = $1"
	row := r.db.QueryRow(query, name)
	var currency domain.Currency
	if err := row.Scan(&currency.ID, &currency.Name, &currency.RateInIDR, &currency.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Currency{}, fmt.Errorf("currency with name %s not found", name)
		}
		return domain.Currency{}, err
	}
	return currency, nil
}
func (r *CurrencyRepository) Save(currency domain.Currency) (domain.Currency, error) {
	query := "INSERT INTO account.currencies (id, name, rate_in_idr, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, currency.ID, currency.Name, currency.RateInIDR, time.Now()).Scan(&currency.ID)
	if err != nil {
		return domain.Currency{}, err
	}
	return currency, nil
}
func (r *CurrencyRepository) Update(currency domain.Currency) (domain.Currency, error) {
	query := "UPDATE account.currencies SET name = $1, rate_in_idr = $2 WHERE id = $3"
	_, err := r.db.Exec(query, currency.Name, currency.RateInIDR, currency.ID)
	if err != nil {
		return domain.Currency{}, err
	}
	return currency, nil
}
func (r *CurrencyRepository) Delete(id string) error {
	query := "DELETE FROM account.currencies WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

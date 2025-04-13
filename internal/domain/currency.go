/*

create table account.currencies(
	id 			char(3) unique not null primary key,
	name 		varchar(255) not null,
	rate_in_idr DECIMAL(15, 2) not null,
	created_at  TIMESTAMP DEFAULT NOW()
);

*/

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

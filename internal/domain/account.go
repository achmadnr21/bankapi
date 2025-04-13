/*


CREATE TABLE account.accounts (
    id           	UUID unique PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      	UUID NOT null,
    employee_id  	UUID null,
    branch_id    	integer not null,
    account_type_id integer not null,
    sequence_number integer not null,
    account_number 	VARCHAR(20) UNIQUE NOT NULL,
    pin_hash 		varchar(255) not null,
    currency_id     VARCHAR(3) NOT NULL DEFAULT 'IDR',
    balance      	DECIMAL(15,2) DEFAULT 0.00,
    created_at   	TIMESTAMP DEFAULT NOW(),
    foreign key (user_id) references profile.users(id) ON DELETE CASCADE,
    foreign key (employee_id) references profile.users(id),
    foreign key (branch_id) references account.branches(id),
    foreign key (account_type_id) references account.account_types(id),
    foreign key (currency_id) references account.currencies(id)
);

*/

package domain

import (
	"time"
)

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
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
type AccountRepository interface {
	FindAll() ([]Account, error)
	FindByID(id string) (Account, error)
	FindByUserID(userID string) ([]Account, error)
	FindByAccountNumber(accountNumber string) (Account, error)
	Save(account Account) (Account, error)
	Update(account Account) (Account, error)
	Delete(id string) error
}

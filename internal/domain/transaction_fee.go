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
*/

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

/*
create table account.account_types(

	id SERIAL unique not null primary key,
	name varchar(255) not null,
	can_etoll boolean default false,
	can_indomart boolean default false,
	can_steam boolean default false,
	created_at  TIMESTAMP DEFAULT NOW()

);
*/
package domain

type AccountType struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	CanEToll    bool   `json:"can_etoll" db:"can_etoll"`
	CanIndomart bool   `json:"can_indomart" db:"can_indomart"`
	CanSteam    bool   `json:"can_steam" db:"can_steam"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}
type AccountTypeRepository interface {
	FindAll() ([]AccountType, error)
	FindByID(id int) (AccountType, error)
	FindByName(name string) (AccountType, error)
	Save(accountType AccountType) (AccountType, error)
	Update(accountType AccountType) (AccountType, error)
	Delete(id int) error
}

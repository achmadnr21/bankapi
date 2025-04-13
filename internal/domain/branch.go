/*
create table account.branches(
	id SERIAL unique not null primary key,
	name varchar(255) not null,
	address text not null,
	created_at  TIMESTAMP DEFAULT NOW()
);
*/

package domain

import "time"

type Branch struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BranchRepository interface {
	FindAll() ([]Branch, error)
	FindByID(id int) (Branch, error)
	FindByName(name string) (Branch, error)
	Save(branch Branch) (Branch, error)
	Update(branch Branch) (Branch, error)
	Delete(id int) error
}

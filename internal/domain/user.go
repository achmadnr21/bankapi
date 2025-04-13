/*
create table profile.user_roles(
	id char(3) unique primary key,
	name varchar(255) not null,
	can_add_user boolean default false,
	can_add_branches boolean default false,
	can_add_account_types boolean default false,
	can_add_currencies boolean default false,
	can_add_accounts boolean default false,
	created_at  TIMESTAMP DEFAULT NOW()
);

create table profile.users(
	id uuid unique default gen_random_uuid() primary key,
	employee_id  UUID null,
	role_id char(3) not null default 'USR',
	full_name   VARCHAR(255) NOT NULL,
	username	VARCHAR(64) unique not null,
    email       VARCHAR(255) UNIQUE NOT NULL,
    phone       VARCHAR(20) UNIQUE NOT NULL,
    password_hash    TEXT NOT NULL,  -- Disimpan dalam bentuk hash
    created_at  TIMESTAMP DEFAULT NOW(),
    foreign key (role_id) references profile.user_roles(id),
    foreign key (employee_id) references profile.users(id)
);

*/

package domain

import "time"

type User struct {
	ID         string    `json:"id"`
	EmployeeID *string   `json:"employee_id"`
	NIK        string    `json:"nik"`
	RoleID     string    `json:"role_id"`
	FullName   string    `json:"full_name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserRepository interface {
	FindAll() ([]User, error)
	Search(nik, username, email string) ([]User, error)
	FindByID(id string) (User, error)
	FindByNIK(nik string) (User, error)
	FindByUsername(username string) (User, error)
	FindByEmail(email string) (User, error)
	FindByPhone(phone string) (User, error)
	Save(user User) (User, error)
	Update(user User) (User, error)
	Delete(id string) error
}

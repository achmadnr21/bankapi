package domain

import "time"

type UserRole struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	CanAddUser         bool      `json:"can_add_user"`
	CanAddBranches     bool      `json:"can_add_branches"`
	CanAddAccountTypes bool      `json:"can_add_account_types"`
	CanAddCurrencies   bool      `json:"can_add_currencies"`
	CanAddAccounts     bool      `json:"can_add_accounts"`
	CreatedAt          time.Time `json:"created_at"`
}

type UserRoleRepository interface {
	FindAll() ([]UserRole, error)
	FindByID(id string) (UserRole, error)
	Save(userRole UserRole) (UserRole, error)
	Update(userRole UserRole) (UserRole, error)
	Delete(id string) error
}

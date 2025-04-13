package repository

import (
	"database/sql"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type UserRoleRepository struct {
	DB *sql.DB
}

func NewUserRoleRepository(db *sql.DB) *UserRoleRepository {
	return &UserRoleRepository{DB: db}
}

func (r *UserRoleRepository) FindAll() ([]domain.UserRole, error) {
	rows, err := r.DB.Query("SELECT * FROM user_roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []domain.UserRole
	for rows.Next() {
		var userRole domain.UserRole
		err := rows.Scan(&userRole.ID, &userRole.Name, &userRole.CanAddUser, &userRole.CanAddBranches, &userRole.CanAddAccountTypes, &userRole.CanAddCurrencies, &userRole.CanAddAccounts, &userRole.CreatedAt)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

func (r *UserRoleRepository) FindByID(id string) (domain.UserRole, error) {
	var userRole domain.UserRole
	err := r.DB.QueryRow("SELECT * FROM profile.user_roles WHERE id = $1", id).Scan(&userRole.ID, &userRole.Name, &userRole.CanAddUser, &userRole.CanAddBranches, &userRole.CanAddAccountTypes, &userRole.CanAddCurrencies, &userRole.CanAddAccounts, &userRole.CreatedAt)
	if err != nil {
		return domain.UserRole{}, err
	}

	return userRole, nil
}

func (r *UserRoleRepository) Save(userRole domain.UserRole) (domain.UserRole, error) {
	_, err := r.DB.Exec("INSERT INTO profile.user_roles (id, name, can_add_user, can_add_branches, can_add_account_types, can_add_currencies, can_add_accounts, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", userRole.ID, userRole.Name, userRole.CanAddUser, userRole.CanAddBranches, userRole.CanAddAccountTypes, userRole.CanAddCurrencies, userRole.CanAddAccounts, userRole.CreatedAt)
	if err != nil {
		return domain.UserRole{}, err
	}

	return userRole, nil
}

func (r *UserRoleRepository) Update(userRole domain.UserRole) (domain.UserRole, error) {
	_, err := r.DB.Exec("UPDATE profile.user_roles SET name = $2, can_add_user = $3, can_add_branches = $4, can_add_account_types = $5, can_add_currencies = $6, can_add_accounts = $7, created_at = $8 WHERE id = $1", userRole.ID, userRole.Name, userRole.CanAddUser, userRole.CanAddBranches, userRole.CanAddAccountTypes, userRole.CanAddCurrencies, userRole.CanAddAccounts, userRole.CreatedAt)
	if err != nil {
		return domain.UserRole{}, err
	}

	return userRole, nil
}

func (r *UserRoleRepository) Delete(id string) error {
	_, err := r.DB.Exec("DELETE FROM profile.user_roles WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

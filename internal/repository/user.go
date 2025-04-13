package repository

import (
	"database/sql"
	"fmt"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{db: DB}
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT * FROM profile.users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Search(nik, username, email string) ([]domain.User, error) {
	var users []domain.User
	query := "SELECT * FROM profile.users WHERE nik = $1 OR username = $2 OR email = $3"
	rows, err := r.db.Query(query, nik, username, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		// clear the user's password
		user.Password = ""
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no user found with nik: %s, username: %s, email: %s", nik, username, email)
	}
	return users, nil
}

func (r *UserRepository) FindByID(id string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT * FROM profile.users WHERE id = $1", id).Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) FindByNIK(nik string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT * FROM profile.users WHERE nik = $1", nik).Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (domain.User, error) {
	var user domain.User
	// Use parameterized query to prevent SQL injection
	query := `SELECT 
	id, employee_id, nik, role_id, full_name, username, email, phone, password_hash, created_at 
	FROM profile.users 
	WHERE username = $1`

	row := r.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, fmt.Errorf("user not found with username: %s", username)
		}
		return domain.User{}, fmt.Errorf("error finding user by username: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT * FROM profile.users WHERE email = $1", email).Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) FindByPhone(phone string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT * FROM profile.users WHERE phone = $1", phone).Scan(&user.ID, &user.EmployeeID, &user.NIK, &user.RoleID, &user.FullName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Save(user domain.User) (domain.User, error) {

	if user.RoleID == "" {
		user.RoleID = "USR"
	}
	_, err := r.db.Exec("INSERT INTO profile.users (employee_id, nik,  role_id, full_name, username, email, phone, password_hash, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", *user.EmployeeID, user.NIK, user.RoleID, user.FullName, user.Username, user.Email, user.Phone, user.Password, user.CreatedAt)
	if err != nil {
		fmt.Println("Error updating user:", err)
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Update(user domain.User) (domain.User, error) {

	_, err := r.db.Exec("UPDATE profile.users SET employee_id = $2, role_id = $3, full_name = $4, username = $5, email = $6, phone = $7, password_hash = $8 WHERE id = $1", user.ID, user.EmployeeID, user.RoleID, user.FullName, user.Username, user.Email, user.Phone, user.Password)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM profile.users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

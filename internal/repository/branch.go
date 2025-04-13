package repository

import (
	"database/sql"

	"github.com/achmadnr21/bankapi/internal/domain"
)

type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository(DB *sql.DB) *BranchRepository {
	return &BranchRepository{db: DB}
}
func (r *BranchRepository) FindAll() ([]domain.Branch, error) {
	rows, err := r.db.Query("SELECT * FROM account.branches")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var branch domain.Branch
		err := rows.Scan(&branch.ID, &branch.Name, &branch.Address, &branch.CreatedAt)
		if err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}

	return branches, nil
}
func (r *BranchRepository) FindByID(id int) (domain.Branch, error) {
	row := r.db.QueryRow("SELECT * FROM account.branches WHERE id = $1", id)
	var branch domain.Branch
	err := row.Scan(&branch.ID, &branch.Name, &branch.Address, &branch.CreatedAt)
	if err != nil {
		return domain.Branch{}, err
	}
	return branch, nil
}

/*
Untuk Find by Name, akan digunakan %name% untuk mencari nama yang mengandung string name,
juga semua dianggap lowercase
*/
func (r *BranchRepository) FindByName(name string) (domain.Branch, error) {
	row := r.db.QueryRow("SELECT * FROM account.branches WHERE LOWER(name) LIKE LOWER($1)", "%"+name+"%")
	var branch domain.Branch
	err := row.Scan(&branch.ID, &branch.Name, &branch.Address, &branch.CreatedAt)
	if err != nil {
		return domain.Branch{}, err
	}
	return branch, nil
}
func (r *BranchRepository) Save(branch domain.Branch) (domain.Branch, error) {
	query := "INSERT INTO account.branches (name, address) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, branch.Name, branch.Address).Scan(&branch.ID)
	if err != nil {
		return domain.Branch{}, err
	}
	return branch, nil
}
func (r *BranchRepository) Update(branch domain.Branch) (domain.Branch, error) {
	query := "UPDATE account.branches SET name = $1, address = $2 WHERE id = $3"
	_, err := r.db.Exec(query, branch.Name, branch.Address, branch.ID)
	if err != nil {
		return domain.Branch{}, err
	}
	return branch, nil
}
func (r *BranchRepository) Delete(id int) error {
	query := "DELETE FROM account.branches WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

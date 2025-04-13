package usecase

import (
	"github.com/achmadnr21/bankapi/internal/domain"
	"github.com/achmadnr21/bankapi/internal/utils"
)

type BranchUsecase struct {
	BranchRepo   domain.BranchRepository
	UserRoleRepo domain.UserRoleRepository
	UserRepo     domain.UserRepository
}

func NewBranchUsecase(branchRepo domain.BranchRepository, userRoleRepo domain.UserRoleRepository, userRepo domain.UserRepository) *BranchUsecase {
	return &BranchUsecase{
		BranchRepo:   branchRepo,
		UserRoleRepo: userRoleRepo,
		UserRepo:     userRepo,
	}
}

func (u *BranchUsecase) GetAllBranches() ([]domain.Branch, error) {
	branches, err := u.BranchRepo.FindAll()
	if err != nil {
		return nil, &utils.InternalServerError{Message: "failed to get branches"}
	}
	return branches, nil
}

func (u *BranchUsecase) GetBranchByID(id int) (domain.Branch, error) {
	branch, err := u.BranchRepo.FindByID(id)
	if err != nil {
		return domain.Branch{}, &utils.NotFoundError{Message: "branch not found"}
	}
	return branch, nil
}

func (u *BranchUsecase) AddBranch(proposerId string, branch domain.Branch) (domain.Branch, error) {
	// get the asker info
	asker, err := u.UserRepo.FindByID(proposerId)
	if err != nil {
		return domain.Branch{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := u.UserRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Branch{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add user
	if !role.CanAddBranches {
		return domain.Branch{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// save branch
	branch, err = u.BranchRepo.Save(branch)
	if err != nil {
		return domain.Branch{}, &utils.InternalServerError{Message: "failed to save branch"}
	}
	return branch, nil
}

func (u *BranchUsecase) UpdateBranch(proposerId string, branchId int, branch domain.Branch) (domain.Branch, error) {
	// get the asker info
	asker, err := u.UserRepo.FindByID(proposerId)
	if err != nil {
		return domain.Branch{}, &utils.NotFoundError{Message: "user not found"}
	}
	// check if asker is admin by checking the role
	role, err := u.UserRoleRepo.FindByID(asker.RoleID)
	if err != nil {
		return domain.Branch{}, &utils.InternalServerError{Message: "failed to get user role"}
	}
	// check user is admin if can add user
	if !role.CanAddBranches {
		return domain.Branch{}, &utils.UnauthorizedError{Message: "user not authorized"}
	}
	// update branch
	branch.ID = branchId
	// check if branch exists
	existingBranch, err := u.BranchRepo.FindByID(branchId)
	if err != nil {
		return domain.Branch{}, &utils.NotFoundError{Message: "branch not found"}
	}
	if existingBranch.ID != branchId {
		return domain.Branch{}, &utils.NotFoundError{Message: "branch not found"}
	}
	branch, err = u.BranchRepo.Update(branch)
	if err != nil {
		return domain.Branch{}, &utils.InternalServerError{Message: "failed to update branch"}
	}
	return branch, nil
}

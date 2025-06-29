package services

import (
	"accounts/internal/api/v1/role_policies/domain/repositories"
)

type RolePoliciesService struct {
	role_policies_repository repositories.RolePoliciesRepository
}

func NewRolePoliciesService(role_policies_repository repositories.RolePoliciesRepository) *RolePoliciesService {
	return &RolePoliciesService{
		role_policies_repository: role_policies_repository,
	}
}

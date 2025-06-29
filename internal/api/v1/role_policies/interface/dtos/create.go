package dtos

import "accounts/internal/api/v1/role_policies/domain/commands"

type CreateRolePoliciesDTO struct {
	RoleID   string `json:"role_id"`
	PolicyID string `json:"policy_id"`
}

func (dto CreateRolePoliciesDTO) Validate() error {
	return nil
}

func (dto CreateRolePoliciesDTO) ToCommand() commands.CreateRolePoliciesCommand {
	return commands.CreateRolePoliciesCommand{
		RoleID:   dto.RoleID,
		PolicyID: dto.PolicyID,
	}
}

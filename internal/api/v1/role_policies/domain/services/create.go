package services

import (
	"accounts/internal/api/v1/role_policies/domain/commands"
	"accounts/internal/api/v1/role_policies/domain/entities"
	"accounts/internal/common/logger"
	"accounts/internal/utils"
	"context"
	"net/http"
)

func (s *RolePoliciesService) Create(ctx context.Context, command commands.CreateRolePoliciesCommand) utils.Responses[entities.RolePoliciesEntity] {
	entry := logger.FromContext(ctx)

	entity := entities.RolePoliciesEntity{
		RoleID:   command.RoleID,
		PolicyID: command.PolicyID,
	}

	res := s.role_policies_repository.SaveEntity(entity)

	if res.Err != nil {
		entry.Error("Error creating policy", "error", res.Err)
		return utils.Responses[entities.RolePoliciesEntity]{
			Err: res.Err,
		}
	}

	return utils.Responses[entities.RolePoliciesEntity]{
		Body:       res.Data,
		StatusCode: http.StatusCreated,
		Success:    true,
	}
}

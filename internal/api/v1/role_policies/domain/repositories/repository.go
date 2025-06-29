package repositories

import (
	"accounts/internal/api/v1/role_policies/domain/entities"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/utils"
)

// --------------------------------
// DOMAIN
// --------------------------------
// Role Policies Repository
// --------------------------------

type RolePoliciesRepository interface {
	SaveEntity(role_policies entities.RolePoliciesEntity) utils.Either[entities.RolePoliciesEntity]
	Search(uuid string) (entities.RolePoliciesEntity, error)
	SearchAll() ([]entities.RolePoliciesEntity, error)
	Delete(uuid string) error
	UpdateByFields(uuid string, fields map[string]interface{}) error
	Matching(criteria criteria.Criteria) ([]entities.RolePoliciesEntity, error)
	View(data []entities.RolePoliciesEntity)
}

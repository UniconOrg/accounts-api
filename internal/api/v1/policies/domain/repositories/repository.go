package repositories

import (
	"accounts/internal/api/v1/policies/domain/entities"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/utils"
)

// --------------------------------
// DOMAIN
// --------------------------------
// Policy Repository
// --------------------------------

type PolicyRepository interface {
	SaveEntity(policy entities.PolicyEntity) utils.Either[entities.PolicyEntity]
	Search(uuid string) (entities.PolicyEntity, error)
	SearchAll() ([]entities.PolicyEntity, error)
	Delete(uuid string) error
	UpdateByFields(uuid string, fields map[string]interface{}) error
	Matching(criteria criteria.Criteria) ([]entities.PolicyEntity, error)
	View(data []entities.PolicyEntity)
}

package policies_pg

import (
	"accounts/internal/api/v1/policies/domain/entities"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/db/postgres"
	"accounts/internal/utils"
	"time"

	"gorm.io/gorm"
)

// --------------------------------
// INFRASTRUCTURE
// --------------------------------
// Policies Postgres Repository
// --------------------------------

type PoliciesPostgresRepository struct {
	postgres.PostgresRepository[entities.PolicyEntity, PolicyModel]
}

func NewPoliciesPostgresRepository(connection *gorm.DB) *PoliciesPostgresRepository {
	return &PoliciesPostgresRepository{
		PostgresRepository: postgres.PostgresRepository[entities.PolicyEntity, PolicyModel]{
			Connection: connection,
		},
	}
}

func (r *PoliciesPostgresRepository) Matching(cr criteria.Criteria) ([]entities.PolicyEntity, error) {

	model := &PolicyModel{}

	return r.MatchingLow(cr, model)
}

func (r *PoliciesPostgresRepository) SaveEntity(policy entities.PolicyEntity) utils.Either[entities.PolicyEntity] {

	res := r.Save(policy)

	if res.Err != nil {
		return utils.Either[entities.PolicyEntity]{
			Err: res.Err,
		}
	}

	policy.ID = res.Data
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	return utils.Either[entities.PolicyEntity]{
		Data: policy,
	}

}

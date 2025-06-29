package rolepolicies_pg

import (
	"accounts/internal/api/v1/role_policies/domain/entities"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/db/postgres"
	"accounts/internal/utils"
	"time"

	"gorm.io/gorm"
)

// --------------------------------
// INFRASTRUCTURE
// --------------------------------
// Role Policies Postgres Repository
// --------------------------------

type RolePoliciesPostgresRepository struct {
	postgres.PostgresRepository[entities.RolePoliciesEntity, RolePoliciesModel]
}

func NewRolePoliciesPostgresRepository(connection *gorm.DB) *RolePoliciesPostgresRepository {
	return &RolePoliciesPostgresRepository{
		PostgresRepository: postgres.PostgresRepository[entities.RolePoliciesEntity, RolePoliciesModel]{
			Connection: connection,
		},
	}
}

func (r *RolePoliciesPostgresRepository) Matching(cr criteria.Criteria) ([]entities.RolePoliciesEntity, error) {

	model := &RolePoliciesModel{}

	return r.MatchingLow(cr, model)
}

func (r *RolePoliciesPostgresRepository) SaveEntity(role_policies entities.RolePoliciesEntity) utils.Either[entities.RolePoliciesEntity] {

	res := r.Save(role_policies)

	if res.Err != nil {
		return utils.Either[entities.RolePoliciesEntity]{
			Err: res.Err,
		}
	}

	role_policies.ID = res.Data
	role_policies.CreatedAt = time.Now()
	role_policies.UpdatedAt = time.Now()

	return utils.Either[entities.RolePoliciesEntity]{
		Data: role_policies,
	}

}

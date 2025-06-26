package postgres

import (
	"accounts/internal/api/v1/codes/domain/entities"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/db/postgres"

	"gorm.io/gorm"
)

// --------------------------------
// INFRASTRUCTURE
// --------------------------------
// PendingRegistrations Postgres Repository
// --------------------------------

type PendingRegistrationsPostgresRepository struct {
	postgres.PostgresRepository[entities.Code, PendingRegistrationModel]
}

func NewPendingRegistrationsPostgresRepository(connection *gorm.DB) *PendingRegistrationsPostgresRepository {
	return &PendingRegistrationsPostgresRepository{
		PostgresRepository: postgres.PostgresRepository[entities.Code, PendingRegistrationModel]{
			Connection: connection,
		},
	}
}

func (r *PendingRegistrationsPostgresRepository) Matching(cr criteria.Criteria) ([]entities.Code, error) {

	model := &PendingRegistrationModel{}

	return r.MatchingLow(cr, model)
}

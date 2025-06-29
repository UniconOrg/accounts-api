package repositories

import (
	"accounts/internal/api/v1/pending_registrations/domain/entities"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/utils"
)

// --------------------------------
// DOMAIN
// --------------------------------
// Pending Registrations Repository
// --------------------------------

type PendingRegistrationsRepository interface {
	Save(role entities.PendingRegistration) utils.Either[string]
	Search(uuid string) (entities.PendingRegistration, error)
	SearchAll() ([]entities.PendingRegistration, error)
	Delete(uuid string) error
	UpdateByFields(uuid string, fields map[string]interface{}) error
	Matching(criteria criteria.Criteria) ([]entities.PendingRegistration, error)
	View(data []entities.PendingRegistration)
}

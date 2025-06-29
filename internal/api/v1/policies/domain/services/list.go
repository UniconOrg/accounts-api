package services

import (
	"accounts/internal/api/v1/policies/domain/entities"
	"context"
)

func (s *PoliciesService) List(ctx context.Context) ([]entities.PolicyEntity, error) {
	return s.policies_repository.SearchAll()
}

package services

import (
	"accounts/internal/api/v1/emails/domain/entities"
	"accounts/internal/common/logger"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/utils"
	"context"
)

func (s EmailsService) ActivateV2(ctx context.Context, entity entities.Activate) utils.Responses[entities.ActivateV2Response] {

	entry := logger.FromContext(ctx)

	entry.Info("Activating user", "email", entity.Email)

	// buscamos el pending registration por el email

	criteria_email := criteria.Criteria{
		Filters: *criteria.NewFilters(
			[]criteria.Filter{
				{
					Field:    "email",
					Value:    entity.Email,
					Operator: criteria.OperatorEqual,
				},
			},
		),
	}

	pending_registrations, err := s.pending_registrations_repository.Matching(criteria_email)
	if err != nil {
		return utils.Responses[entities.ActivateV2Response]{
			StatusCode: 500,
			Errors:     []string{err.Error()},
			Success:    false,
		}
	}

	if len(pending_registrations) == 0 {
		return utils.Responses[entities.ActivateV2Response]{
			StatusCode: 404,
			Errors:     []string{"pending registration not found"},
			Success:    false,
		}
	}

	pending_registration := pending_registrations[0]

	// buscamos el code por el code_id

	criteria_code := criteria.Criteria{
		Filters: *criteria.NewFilters(
			[]criteria.Filter{
				{
					Field:    "id",
					Value:    pending_registration.CodeID,
					Operator: criteria.OperatorEqual,
				},
			},
		),
	}

	codes, err := s.codes_repository.Matching(criteria_code)
	if err != nil {
		return utils.Responses[entities.ActivateV2Response]{
			StatusCode: 500,
			Errors:     []string{err.Error()},
			Success:    false,
		}
	}

	if len(codes) == 0 {
		return utils.Responses[entities.ActivateV2Response]{
			StatusCode: 404,
			Errors:     []string{"code not found"},
			Success:    false,
		}
	}

	code := codes[0]

	// verificamos que sean iguales
	if code.Code != entity.Code {
		return utils.Responses[entities.ActivateV2Response]{
			StatusCode: 400,
			Errors:     []string{"code not valid"},
			Success:    false,
		}
	}

	// crear token con expiracion de 15 minutos
	token, err := s.jwt_controller.GenerateToken(
		ctx,
		map[string]interface{}{
			"user_id": code.UserID,
			"id":      pending_registration.ID,
		},
		15*60,
	)

	if err != nil {
		return utils.Responses[entities.ActivateV2Response]{
			StatusCode: 500,
			Errors:     []string{err.Error()},
			Success:    false,
		}
	}

	return utils.Responses[entities.ActivateV2Response]{
		StatusCode: 200,
		Body:       entities.ActivateV2Response{JWT: token},
		Success:    true,
	}

}

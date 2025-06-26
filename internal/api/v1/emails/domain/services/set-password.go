package services

import (
	"accounts/internal/api/v1/emails/domain/commands"
	"accounts/internal/api/v1/emails/domain/entities"
	login_ents "accounts/internal/api/v1/login_methods/domain/entities"
	"accounts/internal/common/logger"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/utils"
	"context"
)

func (s EmailsService) SetPassword(ctx context.Context, entity commands.SetPassword, jwt string) utils.Responses[entities.ActivateResponse] {

	entry := logger.FromContext(ctx)

	entry.Info("Setting password", "email", entity.Email)

	// validar el JWT
	claims, err := s.jwt_controller.ValidateToken(ctx, jwt)
	if err != nil {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 401,
			Errors:     []string{err.Error()},
			Success:    false,
		}
	}

	// obtenemos el user_id del jwt
	user_id := claims["user_id"].(string)

	// obtenemos el id del pending registration
	pending_registration_id := claims["id"].(string)

	// buscamos el pending registration por el id

	cri := criteria.Criteria{
		Filters: *criteria.NewFilters(
			[]criteria.Filter{
				{
					Field:    "id",
					Value:    pending_registration_id,
					Operator: criteria.OperatorEqual,
				},
				{
					Field:    "email",
					Value:    entity.Email,
					Operator: criteria.OperatorEqual,
				},
			},
		),
	}

	pending_registrations, err := s.pending_registrations_repository.Matching(cri)
	if err != nil {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 500,
			Errors:     []string{err.Error()},
			Success:    false,
		}
	}

	if len(pending_registrations) == 0 {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 404,
			Errors:     []string{"pending registration not found"},
			Success:    false,
		}
	}

	pending_registration := pending_registrations[0]

	// cxreamos eñ email
	email := entities.Email{
		Email:    entity.Email,
		UserID:   user_id,
		Password: entity.Password,
	}

	emailResult := s.repository.Save(email)
	if emailResult.Err != nil {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 500,
			Errors:     []string{emailResult.Err.Error()},
			Success:    false,
		}
	}

	// creamos un login method vinculando el email y le user
	login_method := login_ents.LoginMethod{
		EntityID:   emailResult.Data,
		EntityType: "email",
		UserID:     user_id,
		IsActive:   true,
		IsVerify:   true,
	}

	login_method_result := s.login_methods_repository.Save(login_method)
	if login_method_result.Err != nil {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 500,
			Errors:     []string{login_method_result.Err.Error()},
			Success:    false,
		}
	}

	login_method.ID = login_method_result.Data

	// creamos el refresh token
	refresh_token := s.createRefreshToken(ctx, login_method)
	if refresh_token.Err != nil {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 500,
			Errors:     []string{refresh_token.Err.Error()},
			Success:    false,
		}
	}

	s.codes_repository.UpdateByFields(pending_registration.CodeID, map[string]interface{}{
		"is_removed": true,
	})

	result := s.generateTokens(ctx, login_method, refresh_token.Data)

	if result.Err != nil {
		return utils.Responses[entities.ActivateResponse]{
			StatusCode: 500,
			Errors:     []string{result.Err.Error()},
		}
	}

	s.publishActivationUserEvent(entity.Email, entity.Email)

	return utils.Responses[entities.ActivateResponse]{
		StatusCode: 200,
		Body: entities.ActivateResponse{
			JWT:          result.Data.jwt,
			RefreshToken: result.Data.refresh_token,
		},
	}
}

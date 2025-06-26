package services

import (
	codes_entities "accounts/internal/api/v1/codes/domain/entities"
	email_events "accounts/internal/api/v1/emails/domain/events"
	"accounts/internal/api/v1/pending_registrations/domain/commands"
	pending_registrations_entities "accounts/internal/api/v1/pending_registrations/domain/entities"
	users_entities "accounts/internal/api/v1/users/domain/entities"
	"accounts/internal/common/logger"
	"accounts/internal/core/domain"
	"accounts/internal/core/domain/criteria"
	"accounts/internal/core/domain/event"
	"accounts/internal/utils"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func (s *PendingRegistrationsService) Create(ctx context.Context, command commands.CreatePendingRegistrationCommand) utils.Responses[string] {

	entry := logger.FromContext(ctx)

	entry.Info("Creating pending registration")

	// Validar la exixtencia del role
	cri := criteria.Criteria{
		Filters: *criteria.NewFilters(
			[]criteria.Filter{
				{
					Field:    "name",
					Operator: criteria.OperatorEqual,
					Value:    command.Role,
				},
			},
		),
	}

	roles, err := s.rolesRepository.Matching(cri)
	if err != nil {
		entry.Error("Error getting roles", "error", err)
		return utils.Responses[string]{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("error getting roles %s", err.Error()),
		}
	}

	if len(roles) == 0 {
		entry.Error("Role not found")
		return utils.Responses[string]{
			StatusCode: http.StatusNotFound,
			Err:        fmt.Errorf("role not found %s", command.Role),
		}
	}

	// Crear el user

	user := users_entities.User{
		UserName: command.UserName,
		Role:     command.Role,
		RoleID:   roles[0].ID,
	}

	result := s.usersRepository.Save(user)
	if result.Err != nil {
		entry.Error("Error saving user", "error", result.Err)
		return utils.Responses[string]{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("error saving user %s", result.Err.Error()),
		}
	}

	// Crear el code
	code := codes_entities.Code{
		Entity: domain.Entity{},
		Code:   generateCode(6),
		Type:   "registration",
		UserID: result.Data,
	}

	result = s.codeRepository.Save(code)
	if result.Err != nil {
		entry.Error("Error saving code", "error", result.Err)
		return utils.Responses[string]{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("error saving code %s", result.Err.Error()),
		}
	}

	code.ID = result.Data

	// Crear el pending registration
	pendingRegistration := pending_registrations_entities.PendingRegistration{
		Email:    command.Email,
		UserName: command.UserName,
		Role:     command.Role,
		CodeID:   code.ID,
	}

	result = s.repository.Save(pendingRegistration)
	if result.Err != nil {
		entry.Error("Error saving pending registration", "error", result.Err)
		return utils.Responses[string]{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("error saving pending registration %s", result.Err.Error()),
		}
	}

	pendingRegistration.ID = result.Data

	entry.Info("Pending registration created successfully", "pendingRegistration", pendingRegistration)

	// Enviar el email
	s.publishValidateRegistrationEvent(command.Email, command.UserName, code.Code)

	return utils.Responses[string]{
		Body:       "Check your email to validate your registration",
		StatusCode: http.StatusOK,
	}
}

func generateCode(longitud int) string {
	const numeros = "0123456789"
	resultado := make([]byte, longitud)

	// Se crea un generador local de números aleatorios con semilla.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < longitud; i++ {
		resultado[i] = numeros[r.Intn(len(numeros))]
	}
	return string(resultado)
}

func (s PendingRegistrationsService) publishValidateRegistrationEvent(email string, user_name string, code string) {

	user_event := email_events.UserRegistered{
		Email:            email,
		CodeVerification: code,
		UserName:         user_name,
	}

	// Agregar el mensaje a la cola "new-users"
	if err := s.event_bus.Publish([]event.DomainEvent{
		user_event,
	}); err != nil {
		log.Println("Error al publicar el evento new-users")
		log.Println(err)
	}
}

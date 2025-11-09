package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diogomcd/fake-mill-api/pkg/logger"
	"github.com/go-playground/validator/v10"
)

var (
	validate   *validator.Validate
	cpfRegex   = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$|^\d{11}$`)
	cnpjRegex  = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}$|^\d{14}$`)
	rgRegex    = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}-[\dX]$|^[\d]{8}[\dX]$`)
	cepRegex   = regexp.MustCompile(`^\d{5}-\d{3}$|^\d{8}$`)
	phoneRegex = regexp.MustCompile(`^\+\d{2}\s\(\d{2}\)\s\d{4,5}-\d{4}$|\(\d{2}\)\s\d{4,5}-\d{4}$|^\d{10,11}$`)
)

func init() {
	validate = validator.New()

	// Register custom validators
	if err := validate.RegisterValidation("cpf", validateCPF); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to register CPF validator")
	}

	if err := validate.RegisterValidation("cnpj", validateCNPJ); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to register CNPJ validator")
	}

	if err := validate.RegisterValidation("rg", validateRG); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to register RG validator")
	}

	if err := validate.RegisterValidation("cep", validateCEP); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to register CEP validator")
	}

	if err := validate.RegisterValidation("br_phone", validateBRPhone); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to register BR phone validator")
	}

	if err := validate.RegisterValidation("br_state", validateBRState); err != nil {
		logger.Get().Fatal().Err(err).Msg("Failed to register BR state validator")
	}

	logger.Get().Debug().Msg("Custom validators registered")
}

// Validate validates a struct using validation tags
func Validate(data interface{}) error {
	if err := validate.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// formatValidationErrors formats validation errors into a readable message
func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		messages = append(messages, fmt.Sprintf(
			"field '%s' failed validation '%s'",
			err.Field(),
			err.Tag(),
		))
	}
	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}

// validateCPF validates CPF format (with or without mask)
func validateCPF(fl validator.FieldLevel) bool {
	cpf := fl.Field().String()
	return cpfRegex.MatchString(cpf)
}

// validateCNPJ validates CNPJ format (with or without mask)
func validateCNPJ(fl validator.FieldLevel) bool {
	cnpj := fl.Field().String()
	return cnpjRegex.MatchString(cnpj)
}

// validateRG validates RG format (with or without mask)
func validateRG(fl validator.FieldLevel) bool {
	rg := fl.Field().String()
	return rgRegex.MatchString(rg)
}

// validateCEP validates CEP format (with or without mask)
func validateCEP(fl validator.FieldLevel) bool {
	cep := fl.Field().String()
	return cepRegex.MatchString(cep)
}

// validateBRPhone validates Brazilian phone format
func validateBRPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

// validateBRState validates Brazilian state codes
func validateBRState(fl validator.FieldLevel) bool {
	state := fl.Field().String()
	validStates := []string{
		"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA",
		"MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN",
		"RS", "RO", "RR", "SC", "SP", "SE", "TO",
	}

	for _, s := range validStates {
		if state == s {
			return true
		}
	}
	return false
}

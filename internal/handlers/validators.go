package handlers

import (
	"strconv"
	"strings"

	"github.com/diogomcd/fake-mill-api/internal/generators"
	"github.com/diogomcd/fake-mill-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/ttacon/libphonenumber"
)

// numberTypeToString converte o enum PhoneNumberType da libphonenumber para uma string.
func numberTypeToString(t libphonenumber.PhoneNumberType) string {
	switch t {
	case libphonenumber.FIXED_LINE:
		return "FIXED_LINE"
	case libphonenumber.MOBILE:
		return "MOBILE"
	case libphonenumber.FIXED_LINE_OR_MOBILE:
		return "FIXED_LINE_OR_MOBILE"
	case libphonenumber.TOLL_FREE:
		return "TOLL_FREE"
	case libphonenumber.PREMIUM_RATE:
		return "PREMIUM_RATE"
	case libphonenumber.SHARED_COST:
		return "SHARED_COST"
	case libphonenumber.VOIP:
		return "VOIP"
	case libphonenumber.PERSONAL_NUMBER:
		return "PERSONAL_NUMBER"
	case libphonenumber.PAGER:
		return "PAGER"
	case libphonenumber.UAN:
		return "UAN"
	case libphonenumber.VOICEMAIL:
		return "VOICEMAIL"
	default:
		return "UNKNOWN"
	}
}

// ValidatePhone validates a phone number based on query parameters.
// @Summary Valida Telefone
// @Description Valida um número de telefone e retorna informações detalhadas sobre ele.
// @Tags Validação
// @Accept json
// @Produce json
// @Param phone_number query string true "Número de telefone a validar (com ou sem DDI)"
// @Param country_code query string false "Código do país ISO 3166-1 alpha-2 (ex: BR)"
// @Param ddi query string false "DDI - código de discagem internacional (ex: 55)"
// @Success 200 {object} models.PhoneValidationResponse
// @Failure 400 {object} map[string]string
// @Router /validate/phone [get]
func ValidatePhone(c *fiber.Ctx) error {
	phoneNumber := c.Query("phone_number")
	countryCode := strings.ToUpper(c.Query("country_code"))
	ddi := c.Query("ddi")

	if phoneNumber == "" {
		log.Warn().
			Str("handler", "ValidatePhone").
			Str("error_type", "missing_required_parameter").
			Msg("phone_number parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "phone_number is required",
			"code":  "missing_required_parameter",
		})
	}

	// Define the default region as BR and override it if country_code or ddi are provided.
	region := "BR"
	if countryCode != "" {
		region = countryCode
		log.Debug().
			Str("handler", "ValidatePhone").
			Str("region_source", "country_code").
			Str("region", region).
			Msg("Region determined from country_code")
	} else if ddi != "" {
		if ddiInt, err := strconv.Atoi(ddi); err == nil {
			region = libphonenumber.GetRegionCodeForCountryCode(ddiInt)
			log.Debug().
				Str("handler", "ValidatePhone").
				Str("region_source", "ddi").
				Int("ddi", ddiInt).
				Str("region", region).
				Msg("Region determined from ddi")
		} else {
			log.Warn().
				Err(err).
				Str("handler", "ValidatePhone").
				Str("ddi_input", ddi).
				Str("error_type", "invalid_ddi_format").
				Msg("Invalid ddi format, using default region")
		}
	}

	// Attempts to parse the number. If it fails, return as invalid.
	num, err := libphonenumber.Parse(phoneNumber, region)
	if err != nil {
		log.Warn().
			Err(err).
			Str("handler", "ValidatePhone").
			Str("phone_number", phoneNumber).
			Str("region", region).
			Str("error_type", "parse_failed").
			Msg("Failed to parse phone number")
		return c.Status(fiber.StatusOK).JSON(models.PhoneValidationResponse{
			Valid:       false,
			PhoneNumber: phoneNumber,
		})
	}

	// Checks if the number is valid.
	isValid := libphonenumber.IsValidNumber(num)

	log.Debug().
		Str("handler", "ValidatePhone").
		Str("phone_number", phoneNumber).
		Bool("is_valid", isValid).
		Str("country_code", libphonenumber.GetRegionCodeForNumber(num)).
		Int("ddi", int(num.GetCountryCode())).
		Msg("Phone number validation processed")

	// Prepares the base response.
	response := models.PhoneValidationResponse{
		Valid:       isValid,
		PhoneNumber: phoneNumber,
		CountryCode: libphonenumber.GetRegionCodeForNumber(num),
		DDI:         num.GetCountryCode(),
		NumberType:  "UNKNOWN",
	}

	// If the number is valid, fills in the formatting fields.
	if isValid {
		response.NationalFormat = libphonenumber.Format(num, libphonenumber.NATIONAL)
		response.InternationalFormat = libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
		response.E164Format = libphonenumber.Format(num, libphonenumber.E164)
		response.NumberType = numberTypeToString(libphonenumber.GetNumberType(num))
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ValidateCPFHandler validates a CPF
// @Summary Valida CPF
// @Description Verifica se um número de CPF é válido de acordo com o algoritmo oficial.
// @Tags Validação
// @Accept json
// @Produce json
// @Param cpf path string true "Número de CPF a validar (com ou sem formatação)"
// @Success 200 {object} models.CPFValidationResponse
// @Failure 400 {object} map[string]string
// @Router /validate/cpf/{cpf} [get]
func ValidateCPFHandler(c *fiber.Ctx) error {
	cpf := c.Params("cpf")

	if cpf == "" {
		log.Warn().
			Str("handler", "ValidateCPFHandler").
			Str("error_type", "missing_required_parameter").
			Msg("cpf parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cpf parameter is required",
			"code":  "missing_required_parameter",
		})
	}

	isValid := generators.ValidateCPF(cpf)

	log.Debug().
		Str("handler", "ValidateCPFHandler").
		Str("cpf", cpf).
		Bool("is_valid", isValid).
		Msg("CPF validation processed")

	return c.JSON(models.CPFValidationResponse{
		CPF:   cpf,
		Valid: isValid,
	})
}

// ValidateCNPJHandler validates a CNPJ
// @Summary Valida CNPJ
// @Description Verifica se um número de CNPJ é válido de acordo com o algoritmo oficial.
// @Tags Validação
// @Accept json
// @Produce json
// @Param cnpj path string true "Número de CNPJ a validar (com ou sem formatação)"
// @Success 200 {object} models.CNPJValidationResponse
// @Failure 400 {object} map[string]string
// @Router /validate/cnpj/{cnpj} [get]
func ValidateCNPJHandler(c *fiber.Ctx) error {
	cnpj := c.Params("cnpj")

	if cnpj == "" {
		log.Warn().
			Str("handler", "ValidateCNPJHandler").
			Str("error_type", "missing_required_parameter").
			Msg("cnpj parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cnpj parameter is required",
			"code":  "missing_required_parameter",
		})
	}

	isValid := generators.ValidateCNPJ(cnpj)

	log.Debug().
		Str("handler", "ValidateCNPJHandler").
		Str("cnpj", cnpj).
		Bool("is_valid", isValid).
		Msg("CNPJ validation processed")

	return c.JSON(models.CNPJValidationResponse{
		CNPJ:  cnpj,
		Valid: isValid,
	})
}

// ValidateRGHandler validates a RG
// @Summary Valida RG
// @Description Verifica se um número de RG tem o formato básico correto.
// @Tags Validação
// @Accept json
// @Produce json
// @Param rg path string true "Número de RG a validar (com ou sem formatação)"
// @Success 200 {object} models.RGValidationResponse
// @Failure 400 {object} map[string]string
// @Router /validate/rg/{rg} [get]
func ValidateRGHandler(c *fiber.Ctx) error {
	rg := c.Params("rg")

	if rg == "" {
		log.Warn().
			Str("handler", "ValidateRGHandler").
			Str("error_type", "missing_required_parameter").
			Msg("rg parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rg parameter is required",
			"code":  "missing_required_parameter",
		})
	}

	isValid := generators.ValidateRG(rg)

	log.Debug().
		Str("handler", "ValidateRGHandler").
		Str("rg", rg).
		Bool("is_valid", isValid).
		Msg("RG validation processed")

	return c.JSON(models.RGValidationResponse{
		RG:    rg,
		Valid: isValid,
	})
}

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/diogomcd/fake-mill-api/internal/config"
	"github.com/diogomcd/fake-mill-api/internal/generators"
	"github.com/diogomcd/fake-mill-api/internal/handlers"
	"github.com/diogomcd/fake-mill-api/internal/middleware"
	"github.com/diogomcd/fake-mill-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
)

// @title Fake Mill API
// @version 1.0.0
// @description API REST para geração de dados fake brasileiros para testes, com foco em performance e simplicidade.
// @termsOfService http://swagger.io/terms/

// @license.name GPL v3
// @license.url https://www.gnu.org/licenses/gpl-3.0.html

// @host fakemill.com
// @BasePath /api/v1

// @schemes http https

// @tag.name Pessoa
// @tag.description Endpoints para geração de dados pessoais

// @tag.name Documentos
// @tag.description Endpoints para geração de documentos (CPF, CNPJ, RG)

// @tag.name Contato
// @tag.description Endpoints para geração de emails e telefones

// @tag.name Financeiro
// @tag.description Endpoints para geração de dados bancários e cartões de crédito

// @tag.name Endereço
// @tag.description Endpoints para geração de endereços e CEPs

// @tag.name Empresa
// @tag.description Endpoints para geração de dados de empresas

// @tag.name Validação
// @tag.description Endpoints para validação de documentos e telefones

const (
	apiVersion = "1.0.0"
)

func main() {
	// Load configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	logger.Init()

	log.Info().Msg("Starting Fake Data API")

	// Initialize DataStore (load data)
	dataStore, err := generators.NewDataStore()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize DataStore")
	}

	// Create Generator with DataStore injected
	generator := generators.NewGenerator(dataStore)

	// Create Fiber application
	app := fiber.New(fiber.Config{
		Prefork:       cfg.Server.Prefork,
		CaseSensitive: false,
		ReadTimeout:   cfg.Server.ReadTimeout,
		WriteTimeout:  cfg.Server.WriteTimeout,
		IdleTimeout:   cfg.Server.IdleTimeout,
		AppName:       cfg.Server.AppName,
	})

	// Middleware
	app.Use(middleware.CORSWithConfig(cfg.CORS))
	app.Use(middleware.Logger())

	// Injeta Generator no contexto
	app.Use(middleware.InjectGenerator(generator))

	// Rate Limiting
	if cfg.RateLimit.Enabled {
		rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.Limit, cfg.RateLimit.Window)
		app.Use(rateLimiter.Middleware())
	}

	// Endpoints
	v1 := app.Group("/api/v1")

	// Endpoint: GET /api/v1/person
	v1.Get("/person", handlers.PersonHandler)

	// Document endpoints
	v1.Get("/cpf", handlers.CPFHandler)
	v1.Get("/cnpj", handlers.CNPJHandler)
	v1.Get("/rg", handlers.RGHandler)

	// Contact endpoints
	v1.Get("/email", handlers.EmailHandler)
	v1.Get("/phone", handlers.PhoneHandler)

	// Financial endpoints
	v1.Get("/bank-account", handlers.BankAccountHandler)
	v1.Get("/credit-card", handlers.CreditCardHandler)

	// Address endpoints
	v1.Get("/address", handlers.AddressHandler)
	v1.Get("/zipcode", handlers.ZipcodeHandler)

	// Company endpoint
	v1.Get("/company", handlers.CompanyHandler)

	// Validation endpoints
	v1.Get("/validate/cpf/:cpf", handlers.ValidateCPFHandler)
	v1.Get("/validate/cnpj/:cnpj", handlers.ValidateCNPJHandler)
	v1.Get("/validate/rg/:rg", handlers.ValidateRGHandler)
	v1.Get("/validate/phone", handlers.ValidatePhone)

	// Health check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "ok",
			"version":     apiVersion,
			"name":        cfg.Server.AppName,
			"timestamp":   time.Now().Unix(),
			"environment": cfg.Logging.Env,
		})
	})

	// Swagger JSON file
	app.Get("/api/docs/swagger.json", func(c *fiber.Ctx) error {
		swaggerPath := filepath.Join("docs", "swagger.json")
		content, err := os.ReadFile(swaggerPath)
		if err != nil {
			log.Error().Err(err).Str("path", swaggerPath).Msg("Failed to read swagger.json")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load swagger.json",
			})
		}

		c.Set("Content-Type", "application/json")
		return c.Send(content)
	})

	// Swagger UI
	app.Get("/api/docs/*", swagger.New(swagger.Config{
		URL:         "/api/docs/swagger.json",
		DeepLinking: true,
	}))

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "not found",
			"path":  c.Path(),
		})
	})

	// Start server
	addr := cfg.Address()
	log.Info().
		Str("host", cfg.Server.Host).
		Str("port", cfg.Server.Port).
		Str("addr", addr).
		Msg("🚀 Server started")

	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

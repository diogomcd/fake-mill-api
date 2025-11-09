package generators

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/diogomcd/fake-mill-api/internal/models"
)

// GenerateAddress generates a Brazilian address using the real database
func (g *Generator) GenerateAddress(stateCode, city string) *models.Address {
	ds := g.dataStore

	maxRetries := 10
	var realAddr *RealAddress
	var street, neighborhood string

	for i := 0; i < maxRetries; i++ {
		realAddr = ds.GetRandomRealAddress(stateCode)
		street = strings.TrimSpace(realAddr.Name)
		neighborhood = strings.TrimSpace(realAddr.District)

		if len(street) >= 3 && len(neighborhood) >= 3 {
			break
		}
	}

	if len(street) < 3 {
		if stateCode != "" {
			if cityName := ds.GetRandomCity(stateCode); len(cityName) > 0 {
				street = fmt.Sprintf("Rua %s", cityName)
			} else {
				street = "Rua Principal"
			}
		} else {
			street = "Rua Principal"
		}
		if len(street) < 3 {
			street = "Rua Principal"
		}
	}

	if len(neighborhood) < 3 {
		neighborhood = "Centro"
	}

	number := strconv.Itoa(1 + rand.Intn(9999))

	// Complement, hard code for temporary purposes
	complement := ""
	if rand.Float32() > 0.5 {
		complementOptions := []string{
			"Apto 101", "Apto 202", "Apto 305", "Apto 401",
			"Sala 01", "Sala 102", "Loja 01", "Fundos",
			"Bloco A", "Bloco B", "Bloco C",
		}
		complement = complementOptions[rand.Intn(len(complementOptions))]
	}

	zipcode := formatCEP(realAddr.CEP)
	lat, lng := parseCoordinates(realAddr.Latitude, realAddr.Longitude)

	stateCodeFinal := stateCode
	if stateCodeFinal == "" {
		stateCodeFinal = getStateCodeFromName(realAddr.State)
	}

	cityName := strings.TrimSpace(realAddr.City)
	if len(cityName) < 2 {
		if stateCodeFinal != "" {
			if randomCity := ds.GetRandomCity(stateCodeFinal); len(randomCity) > 0 {
				cityName = randomCity
			} else {
				cityName = "São Paulo"
			}
		} else {
			cityName = "São Paulo"
		}
	}

	return &models.Address{
		Street:       street,
		Number:       number,
		Complement:   complement,
		Neighborhood: neighborhood,
		City:         cityName,
		State:        stateCodeFinal,
		Zipcode:      zipcode,
		Coordinates: &models.Coordinates{
			Lat: lat,
			Lng: lng,
		},
	}
}

// formatCEP formats the CEP in the XXXXX-XXX format
func formatCEP(cep string) string {
	cep = strings.ReplaceAll(cep, "-", "")
	cep = strings.ReplaceAll(cep, ".", "")
	cep = strings.TrimSpace(cep)

	if len(cep) == 8 {
		return fmt.Sprintf("%s-%s", cep[:5], cep[5:])
	}

	return cep
}

// parseCoordinates converts strings of latitude and longitude to float64
func parseCoordinates(lat, lng string) (float64, float64) {
	var latFloat, lngFloat float64

	fmt.Sscanf(lat, "%f", &latFloat)
	fmt.Sscanf(lng, "%f", &lngFloat)

	return latFloat, lngFloat
}

// getStateCodeFromName returns the state code by name
func getStateCodeFromName(stateName string) string {
	stateMap := map[string]string{
		"Acre":                "AC",
		"Alagoas":             "AL",
		"Amapá":               "AP",
		"Amazonas":            "AM",
		"Bahia":               "BA",
		"Ceará":               "CE",
		"Distrito Federal":    "DF",
		"Espírito Santo":      "ES",
		"Goiás":               "GO",
		"Maranhão":            "MA",
		"Mato Grosso":         "MT",
		"Mato Grosso do Sul":  "MS",
		"Minas Gerais":        "MG",
		"Pará":                "PA",
		"Paraíba":             "PB",
		"Paraná":              "PR",
		"Pernambuco":          "PE",
		"Piauí":               "PI",
		"Rio de Janeiro":      "RJ",
		"Rio Grande do Norte": "RN",
		"Rio Grande do Sul":   "RS",
		"Rondônia":            "RO",
		"Roraima":             "RR",
		"Santa Catarina":      "SC",
		"São Paulo":           "SP",
		"Sergipe":             "SE",
		"Tocantins":           "TO",
	}

	if code, exists := stateMap[stateName]; exists {
		return code
	}

	// If the state is not found by name, return the first available
	return "AC"
}

// GenerateZipcode generates a valid CEP
func GenerateZipcode() string {
	firstPart := rand.Intn(100000)
	secondPart := rand.Intn(1000)
	return fmt.Sprintf("%05d-%03d", firstPart, secondPart)
}

// GenerateZipcodeDetails generates all the details of a CEP using real addresses
func (g *Generator) GenerateZipcodeDetails(stateCode string) (formatted, unformatted, state, city string) {
	realAddr := g.dataStore.GetRandomRealAddress(stateCode)

	formatted = formatCEP(realAddr.CEP)
	unformatted = strings.ReplaceAll(formatted, "-", "")
	city = realAddr.City

	if stateCode != "" {
		state = stateCode
	} else {
		state = getStateCodeFromName(realAddr.State)
	}

	return
}

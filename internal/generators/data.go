package generators

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// DataStore stores all the necessary data for generation
type DataStore struct {
	maleFirstNames   []string
	femaleFirstNames []string
	lastNames        []string
	colors           []string
	bloodTypes       []string
	zodiacSigns      []ZodiacSign
	states           []StateData
	stateMap         map[string]*StateData
	dddData          []DDDData
	dddMap           map[string][]int
	professions      []ProfessionData
	educationLevels  []string
	maritalStatuses  []string
	realAddresses    map[string][]RealAddress
	emailShortNames  []string
	emailExtensions  []string
	mu               sync.RWMutex
}

// ZodiacSign represents a zodiac sign
type ZodiacSign struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// ProfessionData represents data of a profession
type ProfessionData struct {
	Title string `json:"title"`
	Area  string `json:"area"`
}

// StateData contains information about a state
type StateData struct {
	Code   string   `json:"code"`
	Name   string   `json:"name"`
	DDD    string   `json:"ddd"`
	Cities []string `json:"cities"`
}

// DDDData contains information about DDD by state
type DDDData struct {
	State     string `json:"state"`
	Code      string `json:"code"`
	DDD_codes []int  `json:"ddd_codes"`
}

// RealAddress represents a real address from the database
type RealAddress struct {
	Name        string `json:"name"`
	AddressType string `json:"addressType"`
	CEP         string `json:"cep"`
	State       string `json:"state"`
	City        string `json:"city"`
	District    string `json:"district"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

// StateAddressData represents the structure of the JSON file of a state
type StateAddressData struct {
	Code   string     `json:"code"`
	Name   string     `json:"name"`
	DDD    string     `json:"ddd"`
	Cities []CityData `json:"cities"`
}

// CityData represents a city in the JSON file
type CityData struct {
	Name      string         `json:"name"`
	IBGECode  string         `json:"ibgeCode"`
	Districts []DistrictData `json:"districts"`
}

// DistrictData represents a neighborhood in the JSON file
type DistrictData struct {
	Name    string        `json:"name"`
	Streets []RealAddress `json:"streets"`
}

// personData struct to deserialize person data
type personData struct {
	FirstNames struct {
		Male   []string `json:"male"`
		Female []string `json:"female"`
	} `json:"firstNames"`
	LastNames       []string         `json:"lastNames"`
	Colors          []string         `json:"colors"`
	BloodTypes      []string         `json:"bloodTypes"`
	ZodiacSigns     []ZodiacSign     `json:"zodiacSigns"`
	Professions     []ProfessionData `json:"professions"`
	EducationLevels []string         `json:"educationLevels"`
	MaritalStatuses []string         `json:"maritalStatuses"`
}

// emailData struct to deserialize email data
type emailData struct {
	DomainExtensions []string `json:"domainExtensions"`
	ShortNames       []string `json:"shortNames"`
}

// Generator encapsulates the logic of generating fake data
// Receives DataStore via dependency injection
type Generator struct {
	dataStore *DataStore
}

// NewGenerator creates a new instance of Generator with DataStore injected
func NewGenerator(ds *DataStore) *Generator {
	return &Generator{
		dataStore: ds,
	}
}

// DataBasePath is the base path for the data
var DataBasePath = "data"

// findProjectRoot finds the root directory of the project by searching for go.mod
func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		// Check if go.mod exists in the current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Go up one level
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root of the file system
			break
		}
		dir = parent
	}

	// If go.mod is not found, return the current directory
	return wd, nil
}

// getDataPath returns the complete path for a data file
func getDataPath(filename string) string {
	projectRoot, err := findProjectRoot()
	if err != nil {
		// Fallback to relative path if the root cannot be found
		return fmt.Sprintf("%s/%s", DataBasePath, filename)
	}
	return filepath.Join(projectRoot, DataBasePath, filename)
}

// NewDataStore creates and initializes a new DataStore with all the data
// Validate that all necessary files exist and contain valid data
func NewDataStore() (*DataStore, error) {
	ds := &DataStore{
		stateMap:      make(map[string]*StateData),
		dddMap:        make(map[string][]int),
		realAddresses: make(map[string][]RealAddress),
	}

	// Load person data
	if err := ds.loadPersonData(); err != nil {
		return nil, fmt.Errorf("error loading person data: %w", err)
	}

	// Load DDD data
	if err := ds.loadDDDData(); err != nil {
		return nil, fmt.Errorf("error loading area code data: %w", err)
	}

	// Load real addresses of states (which also loads cities and states)
	if err := ds.loadRealAddresses(); err != nil {
		return nil, fmt.Errorf("error loading real addresses: %w", err)
	}

	// Load email data
	if err := ds.loadEmailData(); err != nil {
		return nil, fmt.Errorf("error loading email data: %w", err)
	}

	// Validate that all necessary data has been loaded
	if err := ds.validateRequiredData(); err != nil {
		return nil, fmt.Errorf("data validation failed: %w", err)
	}

	return ds, nil
}

// dataValidation represents a data validation
type dataValidation struct {
	name      string
	validator func() bool
	count     int
}

// validateRequiredData validates that all required data has been loaded
func (ds *DataStore) validateRequiredData() error {
	validations := []dataValidation{
		{"male first names", func() bool { return len(ds.maleFirstNames) > 0 }, len(ds.maleFirstNames)},
		{"female first names", func() bool { return len(ds.femaleFirstNames) > 0 }, len(ds.femaleFirstNames)},
		{"last names", func() bool { return len(ds.lastNames) > 0 }, len(ds.lastNames)},
		{"colors", func() bool { return len(ds.colors) > 0 }, len(ds.colors)},
		{"blood types", func() bool { return len(ds.bloodTypes) > 0 }, len(ds.bloodTypes)},
		{"zodiac signs", func() bool { return len(ds.zodiacSigns) > 0 }, len(ds.zodiacSigns)},
		{"professions", func() bool { return len(ds.professions) > 0 }, len(ds.professions)},
		{"education levels", func() bool { return len(ds.educationLevels) > 0 }, len(ds.educationLevels)},
		{"marital statuses", func() bool { return len(ds.maritalStatuses) > 0 }, len(ds.maritalStatuses)},
		{"states", func() bool { return len(ds.states) > 0 }, len(ds.states)},
		{"DDD data", func() bool { return len(ds.dddData) > 0 }, len(ds.dddData)},
		{"real addresses", func() bool { return len(ds.realAddresses) > 0 }, len(ds.realAddresses)},
		{"email short names", func() bool { return len(ds.emailShortNames) > 0 }, len(ds.emailShortNames)},
		{"email extensions", func() bool { return len(ds.emailExtensions) > 0 }, len(ds.emailExtensions)},
	}

	for _, validation := range validations {
		if !validation.validator() {
			return fmt.Errorf("no %s loaded", validation.name)
		}
	}

	log.Info().
		Int("male_names", validations[0].count).
		Int("female_names", validations[1].count).
		Int("last_names", validations[2].count).
		Int("professions", validations[6].count).
		Int("states", validations[9].count).
		Int("ddd_entries", validations[10].count).
		Int("address_states", validations[11].count).
		Int("email_extensions", validations[13].count).
		Int("email_short_names", validations[12].count).
		Msg("All required data validated successfully")

	return nil
}

// GetDataStore returns the internal DataStore (for temporary compatibility)
func (g *Generator) GetDataStore() *DataStore {
	return g.dataStore
}

// loadPersonData loads person data from the JSON file
func (ds *DataStore) loadPersonData() error {
	filePath := getDataPath("person.json")
	log.Debug().Str("file", filePath).Msg("Loading person data")

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to read person data file")
		return err
	}

	var data personData
	if err := json.Unmarshal(content, &data); err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to parse person JSON")
		return err
	}

	ds.maleFirstNames = data.FirstNames.Male
	ds.femaleFirstNames = data.FirstNames.Female
	ds.lastNames = data.LastNames
	ds.colors = data.Colors
	ds.bloodTypes = data.BloodTypes
	ds.zodiacSigns = data.ZodiacSigns
	ds.professions = data.Professions
	ds.educationLevels = data.EducationLevels
	ds.maritalStatuses = data.MaritalStatuses

	log.Info().
		Int("male_names", len(ds.maleFirstNames)).
		Int("female_names", len(ds.femaleFirstNames)).
		Int("last_names", len(ds.lastNames)).
		Int("professions", len(ds.professions)).
		Msg("Person data loaded")

	return nil
}

// GetRandomMaleFirstName returns a random male first name
func (ds *DataStore) GetRandomMaleFirstName() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.maleFirstNames[rand.Intn(len(ds.maleFirstNames))]
}

// GetRandomFemaleFirstName returns a random female first name
func (ds *DataStore) GetRandomFemaleFirstName() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.femaleFirstNames[rand.Intn(len(ds.femaleFirstNames))]
}

// GetRandomLastName returns a random last name
func (ds *DataStore) GetRandomLastName() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.lastNames[rand.Intn(len(ds.lastNames))]
}

// GetRandomColor returns a random favorite color
func (ds *DataStore) GetRandomColor() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.colors[rand.Intn(len(ds.colors))]
}

// GetRandomBloodType returns a random blood type
func (ds *DataStore) GetRandomBloodType() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.bloodTypes[rand.Intn(len(ds.bloodTypes))]
}

// GetZodiacSign returns the zodiac sign based on the birthdate
func (ds *DataStore) GetZodiacSign(birthdate string) string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Parse the birthdate (format: YYYY-MM-DD)
	parts := strings.Split(birthdate, "-")
	if len(parts) != 3 {
		// If the format is invalid, return the first available sign
		return ds.zodiacSigns[0].Name
	}

	monthDay := fmt.Sprintf("%s-%s", parts[1], parts[2])

	for _, sign := range ds.zodiacSigns {
		if (monthDay >= sign.Start && monthDay <= sign.End) ||
			(sign.Name == "Capricórnio" && ((monthDay >= "12-22" && monthDay <= "12-31") || (monthDay >= "01-01" && monthDay <= "01-19"))) {
			return sign.Name
		}
	}

	// If no sign is found, return the first available sign
	return ds.zodiacSigns[0].Name
}

// GetRandomState returns a random state
func (ds *DataStore) GetRandomState() *StateData {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return &ds.states[rand.Intn(len(ds.states))]
}

// GetStateByCode returns a state by code (ex: "SP")
func (ds *DataStore) GetStateByCode(code string) *StateData {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.stateMap[strings.ToUpper(code)]
}

// ValidateAndSanitizeState validates and sanitizes a state code
// Returns the sanitized code in uppercase if valid, or an empty string if invalid
func (ds *DataStore) ValidateAndSanitizeState(stateCode string) string {
	if stateCode == "" {
		return ""
	}

	sanitized := strings.ToUpper(strings.TrimSpace(stateCode))

	if len(sanitized) != 2 {
		return ""
	}

	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if _, exists := ds.stateMap[sanitized]; !exists {
		return ""
	}

	return sanitized
}

// GetRandomCity returns a random city from a specific state
func (ds *DataStore) GetRandomCity(stateCode string) string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	state := ds.stateMap[strings.ToUpper(stateCode)]
	return state.Cities[rand.Intn(len(state.Cities))]
}

// GetRandomCityFromState returns a random city from a specific state
func (ds *DataStore) GetRandomCityFromState(state *StateData) string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return state.Cities[rand.Intn(len(state.Cities))]
}

// loadDDDData loads DDD data from the JSON file
func (ds *DataStore) loadDDDData() error {
	filePath := getDataPath("ddd.json")
	log.Debug().Str("file", filePath).Msg("Loading DDD data")

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to read DDD file")
		return err
	}

	if err := json.Unmarshal(content, &ds.dddData); err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to parse DDD JSON")
		return err
	}

	// Fill the DDD map by state code
	for _, dddItem := range ds.dddData {
		ds.dddMap[dddItem.Code] = dddItem.DDD_codes
	}

	log.Info().Int("ddd_entries", len(ds.dddData)).Msg("DDD data loaded")

	return nil
}

// GetDDDForState returns a random DDD for a specific state
func (ds *DataStore) GetDDDForState(stateCode string) string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	dddCodes := ds.dddMap[strings.ToUpper(stateCode)]
	ddd := dddCodes[rand.Intn(len(dddCodes))]
	return fmt.Sprintf("%02d", ddd)
}

// GetRandomProfession returns a random profession
func (ds *DataStore) GetRandomProfession() ProfessionData {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.professions[rand.Intn(len(ds.professions))]
}

// GetRandomEducationLevel returns a random education level
func (ds *DataStore) GetRandomEducationLevel() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.educationLevels[rand.Intn(len(ds.educationLevels))]
}

// GetRandomMaritalStatus returns a random marital status
func (ds *DataStore) GetRandomMaritalStatus() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.maritalStatuses[rand.Intn(len(ds.maritalStatuses))]
}

// loadRealAddresses loads all real addresses from the JSON files of states
func (ds *DataStore) loadRealAddresses() error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	statesPath := filepath.Join(projectRoot, DataBasePath, "states")
	log.Info().Str("path", statesPath).Msg("Starting real addresses loading")

	// List the files in the states directory
	entries, err := os.ReadDir(statesPath)
	if err != nil {
		return fmt.Errorf("failed to read states directory: %w", err)
	}

	totalAddresses := 0
	completeAddresses := 0

	// Map to track processed states (case-insensitive)
	processedStates := make(map[string]bool)

	for _, entry := range entries {
		// Ignore directories and files that are not JSON
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}

		// Extract the state code from the file name (remove .json and convert to lowercase)
		fileName := entry.Name()
		stateCode := strings.ToLower(strings.TrimSuffix(fileName, ".json"))

		// Ignore if we have already processed this state (case-insensitive)
		if processedStates[stateCode] {
			continue
		}
		processedStates[stateCode] = true

		filePath := filepath.Join(statesPath, fileName)

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Warn().Err(err).Str("file", filePath).Msg("Failed to read status file, continuing...")
			continue
		}

		var stateData StateAddressData
		if err := json.Unmarshal(content, &stateData); err != nil {
			log.Warn().Err(err).Str("file", filePath).Msg("Failed to parse state JSON, continuing...")
			continue
		}

		stateCodeUpper := strings.ToUpper(stateCode)
		addresses := []RealAddress{}
		cityNames := []string{}

		// Extract city names and addresses
		for _, city := range stateData.Cities {
			cityNames = append(cityNames, city.Name)
			for _, district := range city.Districts {
				for _, street := range district.Streets {
					totalAddresses++

					// Filter only complete addresses (with latitude and longitude)
					if street.Latitude != "" && street.Longitude != "" {
						addresses = append(addresses, street)
						completeAddresses++
					}
				}
			}
		}

		// Fill state data
		stateInfo := StateData{
			Code:   stateCodeUpper,
			Name:   stateData.Name,
			DDD:    stateData.DDD,
			Cities: cityNames,
		}

		ds.states = append(ds.states, stateInfo)
		ds.stateMap[stateCodeUpper] = &ds.states[len(ds.states)-1]

		if len(addresses) > 0 {
			ds.realAddresses[stateCodeUpper] = addresses
			log.Info().
				Str("state", stateCodeUpper).
				Int("cities", len(cityNames)).
				Int("total_addresses", len(addresses)).
				Msg("Addresses loaded for state")
		} else if len(cityNames) > 0 {
			log.Warn().
				Str("state", stateCodeUpper).
				Int("cities", len(cityNames)).
				Int("total_addresses", 0).
				Msg("State loaded but no complete addresses found (missing coordinates)")
		}
	}

	log.Info().
		Int("total_addresses", totalAddresses).
		Int("complete_addresses", completeAddresses).
		Int("states_loaded", len(ds.realAddresses)).
		Int("states_indexed", len(ds.states)).
		Msg("Real addresses loading completed")

	if completeAddresses == 0 {
		return fmt.Errorf("no complete addresses found")
	}

	return nil
}

// GetRandomRealAddress returns a random real address
// If stateCode is empty, choose a random state
func (ds *DataStore) GetRandomRealAddress(stateCode string) *RealAddress {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// If no state is specified, choose a random state
	if stateCode == "" {
		// Cria lista de estados disponíveis
		availableStates := make([]string, 0, len(ds.realAddresses))
		for state := range ds.realAddresses {
			availableStates = append(availableStates, state)
		}

		stateCode = availableStates[rand.Intn(len(availableStates))]
	}

	addresses := ds.realAddresses[strings.ToUpper(stateCode)]
	idx := rand.Intn(len(addresses))
	return &addresses[idx]
}

// loadEmailData loads email data from the JSON file
func (ds *DataStore) loadEmailData() error {
	filePath := getDataPath("email.json")
	log.Debug().Str("file", filePath).Msg("Loading email data")

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to read email data file")
		return err
	}

	var data emailData
	if err := json.Unmarshal(content, &data); err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Failed to parse email JSON")
		return err
	}

	ds.emailShortNames = data.ShortNames
	ds.emailExtensions = data.DomainExtensions

	log.Info().
		Int("email_short_names", len(ds.emailShortNames)).
		Int("email_extensions", len(ds.emailExtensions)).
		Msg("Email data loaded")

	return nil
}

// GetRandomEmailShortName returns a random short name for email
func (ds *DataStore) GetRandomEmailShortName() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.emailShortNames[rand.Intn(len(ds.emailShortNames))]
}

// GetRandomEmailExtension returns a random domain extension
func (ds *DataStore) GetRandomEmailExtension() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.emailExtensions[rand.Intn(len(ds.emailExtensions))]
}

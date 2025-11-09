package generators

import (
	"fmt"
	"math/rand"
	"strings"
)

// GenerateEmail generates a fake email with a realistic username and dynamic domain
func (g *Generator) GenerateEmail(customDomain string) (email, username, domain string) {
	// If a custom domain is provided, use it directly
	if customDomain != "" {
		email, username, domain = g.generateEmailWithDomain(customDomain)
		return
	}

	// Generate a random domain using short names and extensions
	domain = g.generateRandomDomain()
	email, username, _ = g.generateEmailWithDomain(domain)

	return
}

// generateRandomDomain generates a random domain combining short names and extensions
func (g *Generator) generateRandomDomain() string {
	ds := g.dataStore
	shortName := ds.GetRandomEmailShortName()
	extension := ds.GetRandomEmailExtension()
	return fmt.Sprintf("%s.%s", shortName, extension)
}

// generateEmailWithDomain generates an email with a realistic username for a specific domain
func (g *Generator) generateEmailWithDomain(domain string) (email, username string, returnDomain string) {
	ds := g.dataStore

	// Generate a random name for the email
	firstName := ds.GetRandomMaleFirstName()
	lastName := ds.GetRandomLastName()

	// Clean the names to use as base for the username
	firstParts := strings.Fields(firstName)
	lastParts := strings.Fields(lastName)

	cleanFirst := strings.ToLower(firstParts[0])
	cleanLast := ""
	if len(lastParts) > 0 {
		cleanLast = strings.ToLower(lastParts[0])
	}

	// Choose one of the most realistic and short formats (logic from person.go)
	format := rand.Intn(6)
	switch format {
	case 0:
		// name.lastname
		username = fmt.Sprintf("%s.%s", cleanFirst, cleanLast)
	case 1:
		// name_lastname
		username = fmt.Sprintf("%s_%s", cleanFirst, cleanLast)
	case 2:
		// name + initial of lastname (very common)
		if len(cleanLast) > 0 {
			username = fmt.Sprintf("%s%s", cleanFirst, cleanLast[:1])
		} else {
			username = cleanFirst
		}
	case 3:
		// only first name
		username = cleanFirst
	case 4:
		// name + birth year (realistic)
		year := 85 + rand.Intn(35) // 1985-2019
		username = fmt.Sprintf("%s%d", cleanFirst, year)
	case 5:
		// abbreviated name + abbreviated lastname
		if len(cleanFirst) > 3 && len(cleanLast) > 2 {
			username = fmt.Sprintf("%s_%s", cleanFirst[:3], cleanLast[:3])
		} else {
			username = fmt.Sprintf("%s.%s", cleanFirst, cleanLast)
		}
	}

	// Add numbers (high chance to be more realistic)
	if rand.Intn(100) < 70 {
		if format == 4 {
			// if already has year, add only some extra numbers
			if rand.Intn(100) < 50 {
				number := rand.Intn(99) + 1
				username = fmt.Sprintf("%s%d", username, number)
			}
		} else {
			// add random numbers
			number := rand.Intn(999) + 1
			if rand.Intn(100) < 60 {
				// only 1-2 digits to keep short
				number = rand.Intn(99) + 1
			}
			username = fmt.Sprintf("%s%d", username, number)
		}
	}

	// Remove accents and sanitize the username
	username = sanitizeUsername(removeAccents(username))

	// Format the email
	email = fmt.Sprintf("%s@%s", username, domain)

	return email, username, domain
}

// removeAccents remove accents from the text
func removeAccents(s string) string {
	replacements := map[string]string{
		"á": "a", "à": "a", "ã": "a", "â": "a", "ä": "a",
		"Á": "A", "À": "A", "Ã": "A", "Â": "A", "Ä": "A",
		"é": "e", "è": "e", "ê": "e", "ë": "e",
		"É": "E", "È": "E", "Ê": "E", "Ë": "E",
		"í": "i", "ì": "i", "î": "i", "ï": "i",
		"Í": "I", "Ì": "I", "Î": "I", "Ï": "I",
		"ó": "o", "ò": "o", "õ": "o", "ô": "o", "ö": "o",
		"Ó": "O", "Ò": "O", "Õ": "O", "Ô": "O", "Ö": "O",
		"ú": "u", "ù": "u", "û": "u", "ü": "u",
		"Ú": "U", "Ù": "U", "Û": "U", "Ü": "U",
		"ç": "c", "Ç": "C",
		"ñ": "n", "Ñ": "N",
	}

	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// sanitizeUsername remove invalid characters for email and ensure it is valid
func sanitizeUsername(username string) string {
	// Remove spaces and convert to lowercase
	username = strings.ToLower(strings.TrimSpace(username))

	// Remove invalid special characters, keeping only letters, numbers, dots, hyphens and underscores
	var result strings.Builder
	for _, r := range username {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	sanitized := result.String()

	// Ensure it does not start or end with dot, hyphen or underscore
	sanitized = strings.Trim(sanitized, ".-_")

	// Ensure it has at least 3 characters
	if len(sanitized) < 3 {
		sanitized = fmt.Sprintf("user%d", rand.Intn(9999)+1000)
	}

	// Ensure it does not have consecutive dots
	for strings.Contains(sanitized, "..") {
		sanitized = strings.ReplaceAll(sanitized, "..", ".")
	}

	return sanitized
}

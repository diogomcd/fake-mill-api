package generators

import (
	"fmt"
	"math/rand"
	"strings"
)

// GeneratePhone generates a Brazilian phone number
func (g *Generator) GeneratePhone(stateCode, requestedType string) (phone, ddd, state, phoneType string) {
	ds := g.dataStore

	var selectedState *StateData
	if stateCode != "" {
		selectedState = ds.GetStateByCode(stateCode)
	}

	if selectedState == nil {
		selectedState = ds.GetRandomState()
	}

	state = selectedState.Code
	ddd = ds.GetDDDForState(state)

	if requestedType == "mobile" || requestedType == "landline" {
		phoneType = requestedType
	} else {
		types := []string{"mobile", "landline"}
		phoneType = types[rand.Intn(len(types))]
	}

	var phoneNumber string
	if phoneType == "mobile" {
		// Mobile: 9XXXX-XXXX
		phoneNumber = fmt.Sprintf("9%d%d%d%d-%d%d%d%d",
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		)
	} else {
		// Landline: 3XXXX-XXXX or 4XXXX-XXXX, etc
		firstPart := 3 + rand.Intn(5) // 3 to 7
		phoneNumber = fmt.Sprintf("%d%d%d%d-%d%d%d%d",
			firstPart, rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		)
	}

	// Format with DDD: (XX) XXXXX-XXXX
	phone = fmt.Sprintf("(%s) %s", ddd, phoneNumber)

	return
}

// FormatPhone formats a phone number
func FormatPhone(phone string) string {
	clean := strings.ReplaceAll(phone, "(", "")
	clean = strings.ReplaceAll(clean, ")", "")
	clean = strings.ReplaceAll(clean, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")

	if len(clean) < 10 {
		return phone
	}

	// Extract DDD (first 2 digits)
	ddd := clean[:2]

	// Rest of the number
	rest := clean[2:]

	// Format: (XX) XXXXX-XXXX
	if len(rest) == 8 {
		return fmt.Sprintf("(%s) %s-%s", ddd, rest[:4], rest[4:])
	}

	if len(rest) == 9 {
		return fmt.Sprintf("(%s) %s-%s", ddd, rest[:5], rest[5:])
	}

	return phone
}

// UnformatPhone removes formatting from the phone number
func UnformatPhone(phone string) string {
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return phone
}

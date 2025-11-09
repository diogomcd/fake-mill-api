package generators

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// GenerateCreditCard generates fake credit card data
func (g *Generator) GenerateCreditCard(brand string) (number, cardBrand, cvv, expirationDate, holderName string) {
	brands := []string{"visa", "mastercard", "elo", "amex"}

	if brand != "" {
		cardBrand = capitalizeFirst(strings.ToLower(brand))
	} else {
		cardBrand = capitalizeFirst(brands[rand.Intn(len(brands))])
	}

	// Generate card number (simplified, only the format)
	number = fmt.Sprintf("%d%d%d%d %d%d%d%d %d%d%d%d %d%d%d%d",
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
	)

	cvv = fmt.Sprintf("%03d", rand.Intn(1000))

	// Generate expiration date (2 to 5 years from now) in the MM/YY format (5 characters)
	now := time.Now()
	expirationYear := now.Year() + 2 + rand.Intn(4)
	expirationMonth := 1 + rand.Intn(12)
	expirationDate = fmt.Sprintf("%02d/%02d", expirationMonth, expirationYear%100)

	ds := g.dataStore
	firstName := ds.GetRandomMaleFirstName() // Simplified
	lastName := ds.GetRandomLastName()
	holderName = fmt.Sprintf("%s %s", firstName, lastName)

	return
}

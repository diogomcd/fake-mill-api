package generators

import (
	"fmt"
	"math/rand"
)

// Bank represents a Brazilian bank
type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var banks = []Bank{
	{Code: "001", Name: "Banco do Brasil"},
	{Code: "237", Name: "Bradesco"},
	{Code: "341", Name: "Itaú Unibanco"},
	{Code: "104", Name: "Caixa Econômica Federal"},
	{Code: "033", Name: "Santander"},
	{Code: "260", Name: "Nu Pagamentos S.A. (Nubank)"},
	{Code: "077", Name: "Banco Inter"},
	{Code: "336", Name: "Banco C6 S.A."},
}

// GenerateBankAccount generates fake bank account data
func (g *Generator) GenerateBankAccount(bankCode string) (bank Bank, agency, account, accountType string) {
	if bankCode != "" {
		for _, b := range banks {
			if b.Code == bankCode {
				bank = b
				break
			}
		}
	}

	// If no bank is found or specified, choose one randomly
	if bank.Code == "" {
		bank = banks[rand.Intn(len(banks))]
	}

	agency = fmt.Sprintf("%04d-%d", rand.Intn(9999)+1, rand.Intn(10))
	account = fmt.Sprintf("%08d-%d", rand.Intn(99999999)+1, rand.Intn(10))
	types := []string{"checking", "savings"}
	accountType = types[rand.Intn(len(types))]

	return
}

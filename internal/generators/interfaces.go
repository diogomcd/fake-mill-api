package generators

import "github.com/diogomcd/fake-mill-api/internal/models"

// DocumentGenerator define interface for generating Brazilian documents
type DocumentGenerator interface {
	GenerateCPF(formatted bool, valid bool) string
	GenerateCNPJ(formatted bool, valid bool) string
	GenerateRG(stateCode string, formatted bool, valid bool) (rg, state, issuer, issueDate, expirationDate string)
}

// ContactGenerator define interface for generating contact information
type ContactGenerator interface {
	GenerateEmail(customDomain string) (email, username, domain string)
	GeneratePhone(stateCode, requestedType string) (phone, ddd, state, phoneType string)
}

// PersonGenerator define interface for generating person profiles
type PersonGenerator interface {
	GeneratePerson(gender, stateCode string) *models.Person
}

// AddressGenerator define interface for generating addresses
type AddressGenerator interface {
	GenerateAddress(stateCode, city string) *models.Address
	GenerateZipcodeDetails(stateCode string) (formatted, unformatted, state, city string)
}

// FinancialGenerator define interface for generating financial data
type FinancialGenerator interface {
	GenerateBankAccount(bankCode string) (bank Bank, agency, account, accountType string)
	GenerateCreditCard(brand string) (number, cardBrand, cvv, expirationDate, holderName string)
}

// CompanyGenerator define interface for generating companies
type CompanyGenerator interface {
	GenerateCompany() *models.CompanyResponse
}

// DataStoreProvider define interface for accessing the DataStore
type DataStoreProvider interface {
	GetDataStore() *DataStore
}

// IGenerator combines all specific interfaces
// Kept for compatibility with existing code
// The Generator (struct) implements this interface implicitly
type IGenerator interface {
	DocumentGenerator
	ContactGenerator
	PersonGenerator
	AddressGenerator
	FinancialGenerator
	CompanyGenerator
	DataStoreProvider
}

// CPFGenerator define interface for specific CPF generators
type CPFGenerator interface {
	Generate() interface{}
	Validate(value string) bool
}

// CNPJGenerator define interface for specific CNPJ generators
type CNPJGenerator interface {
	Generate() interface{}
	Validate(value string) bool
}

// RGGenerator define interface for specific RG generators
type RGGenerator interface {
	Generate() interface{}
	Validate(value string) bool
}

// EmailGenerator define interface for specific email generators
type EmailGenerator interface {
	Generate() interface{}
	Validate(value string) bool
}

// PhoneGenerator define interface for specific phone generators
type PhoneGenerator interface {
	Generate() interface{}
	Validate(value string) bool
}

package generators

import "github.com/diogomcd/fake-mill-api/internal/models"

// MockGenerator is an example of mock Generator for tests
// Implements the IGenerator interface
type MockGenerator struct {
	MockGenerateCPF            func(formatted bool, valid bool) string
	MockGenerateCNPJ           func(formatted bool, valid bool) string
	MockGenerateRG             func(stateCode string, formatted bool, valid bool) (rg, state, issuer, issueDate, expirationDate string)
	MockGenerateEmail          func(customDomain string) (email, username, domain string)
	MockGeneratePhone          func(stateCode, requestedType string) (phone, ddd, state, phoneType string)
	MockGeneratePerson         func(gender, stateCode string) *models.Person
	MockGenerateAddress        func(stateCode, city string) *models.Address
	MockGenerateZipcodeDetails func(stateCode string) (formatted, unformatted, state, city string)
	MockGenerateBankAccount    func(bankCode string) (bank Bank, agency, account, accountType string)
	MockGenerateCreditCard     func(brand string) (number, cardBrand, cvv, expirationDate, holderName string)
	MockGenerateCompany        func() *models.CompanyResponse
	MockGetDataStore           func() *DataStore
}

func (m *MockGenerator) GenerateCPF(formatted bool, valid bool) string {
	if m.MockGenerateCPF != nil {
		return m.MockGenerateCPF(formatted, valid)
	}
	return "00000000000"
}

func (m *MockGenerator) GenerateCNPJ(formatted bool, valid bool) string {
	if m.MockGenerateCNPJ != nil {
		return m.MockGenerateCNPJ(formatted, valid)
	}
	return "00000000000000"
}

func (m *MockGenerator) GenerateRG(stateCode string, formatted bool, valid bool) (rg, state, issuer, issueDate, expirationDate string) {
	if m.MockGenerateRG != nil {
		return m.MockGenerateRG(stateCode, formatted, valid)
	}
	return "000000000", "SP", "SSP", "2020-01-01", "2030-01-01"
}

func (m *MockGenerator) GenerateEmail(customDomain string) (email, username, domain string) {
	if m.MockGenerateEmail != nil {
		return m.MockGenerateEmail(customDomain)
	}
	return "test@example.com", "test", "example.com"
}

func (m *MockGenerator) GeneratePhone(stateCode, requestedType string) (phone, ddd, state, phoneType string) {
	if m.MockGeneratePhone != nil {
		return m.MockGeneratePhone(stateCode, requestedType)
	}
	return "(11) 99999-9999", "11", "SP", "mobile"
}

func (m *MockGenerator) GeneratePerson(gender, stateCode string) *models.Person {
	if m.MockGeneratePerson != nil {
		return m.MockGeneratePerson(gender, stateCode)
	}
	return &models.Person{}
}

func (m *MockGenerator) GenerateAddress(stateCode, city string) *models.Address {
	if m.MockGenerateAddress != nil {
		return m.MockGenerateAddress(stateCode, city)
	}
	return &models.Address{}
}

func (m *MockGenerator) GenerateZipcodeDetails(stateCode string) (formatted, unformatted, state, city string) {
	if m.MockGenerateZipcodeDetails != nil {
		return m.MockGenerateZipcodeDetails(stateCode)
	}
	return "00000-000", "00000000", "SP", "São Paulo"
}

func (m *MockGenerator) GenerateBankAccount(bankCode string) (bank Bank, agency, account, accountType string) {
	if m.MockGenerateBankAccount != nil {
		return m.MockGenerateBankAccount(bankCode)
	}
	return Bank{Code: "001", Name: "Banco do Brasil"}, "0001-0", "00000000-0", "checking"
}

func (m *MockGenerator) GenerateCreditCard(brand string) (number, cardBrand, cvv, expirationDate, holderName string) {
	if m.MockGenerateCreditCard != nil {
		return m.MockGenerateCreditCard(brand)
	}
	return "0000 0000 0000 0000", "visa", "000", "12/2025", "Test User"
}

func (m *MockGenerator) GenerateCompany() *models.CompanyResponse {
	if m.MockGenerateCompany != nil {
		return m.MockGenerateCompany()
	}
	return &models.CompanyResponse{}
}

func (m *MockGenerator) GetDataStore() *DataStore {
	if m.MockGetDataStore != nil {
		return m.MockGetDataStore()
	}
	return nil
}

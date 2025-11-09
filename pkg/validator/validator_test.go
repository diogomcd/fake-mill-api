package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	CPF   string `validate:"required,cpf"`
	CNPJ  string `validate:"required,cnpj"`
	Email string `validate:"required,email"`
	State string `validate:"required,br_state"`
}

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name  string
		cpf   string
		valid bool
	}{
		{"Valid masked CPF", "123.456.789-10", true},
		{"Valid unmasked CPF", "12345678910", true},
		{"Invalid CPF", "123.456.789", false},
		{"Empty CPF", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				CPF string `validate:"required,cpf"`
			}{CPF: tt.cpf}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateCNPJ(t *testing.T) {
	tests := []struct {
		name  string
		cnpj  string
		valid bool
	}{
		{"Valid masked CNPJ", "12.345.678/0001-90", true},
		{"Valid unmasked CNPJ", "12345678000190", true},
		{"Invalid CNPJ", "12.345.678/0001", false},
		{"Empty CNPJ", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				CNPJ string `validate:"required,cnpj"`
			}{CNPJ: tt.cnpj}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateRG(t *testing.T) {
	tests := []struct {
		name  string
		rg    string
		valid bool
	}{
		{"Valid masked RG", "12.345.678-9", true},
		{"Valid unmasked RG", "123456789", true},
		{"Valid masked RG with X", "12.345.678-X", true},
		{"Valid unmasked RG with X", "12345678X", true},
		{"Invalid RG", "12.345.678", false},
		{"Empty RG", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				RG string `validate:"required,rg"`
			}{RG: tt.rg}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name  string
		cep   string
		valid bool
	}{
		{"Valid masked CEP", "12345-678", true},
		{"Valid unmasked CEP", "12345678", true},
		{"Invalid CEP", "12345-67", false},
		{"Empty CEP", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				CEP string `validate:"required,cep"`
			}{CEP: tt.cep}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateBRState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		valid bool
	}{
		{"Valid state SP", "SP", true},
		{"Valid state RJ", "RJ", true},
		{"Valid state MG", "MG", true},
		{"Invalid state XX", "XX", false},
		{"Empty state", "", false},
		{"Lowercase state", "sp", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				State string `validate:"required,br_state"`
			}{State: tt.state}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateBRPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{"Valid international format", "+55 (11) 98765-4321", true},
		{"Valid unmasked mobile", "11987654321", true},
		{"Valid unmasked landline", "1134567890", true},
		{"Invalid phone", "123", false},
		{"Empty phone", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				Phone string `validate:"required,br_phone"`
			}{Phone: tt.phone}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Invalid email", "invalid-email", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				Email string `validate:"required,email"`
			}{Email: tt.email}

			err := Validate(data)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	t.Run("Valid struct", func(t *testing.T) {
		data := TestStruct{
			CPF:   "123.456.789-10",
			CNPJ:  "12.345.678/0001-90",
			Email: "test@example.com",
			State: "SP",
		}

		err := Validate(data)
		assert.NoError(t, err)
	})

	t.Run("Invalid struct - missing CPF", func(t *testing.T) {
		data := TestStruct{
			CPF:   "",
			CNPJ:  "12.345.678/0001-90",
			Email: "test@example.com",
			State: "SP",
		}

		err := Validate(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF")
	})

	t.Run("Invalid struct - invalid email", func(t *testing.T) {
		data := TestStruct{
			CPF:   "123.456.789-10",
			CNPJ:  "12.345.678/0001-90",
			Email: "invalid-email",
			State: "SP",
		}

		err := Validate(data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})
}

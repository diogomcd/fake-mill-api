package models

// Address represents a Brazilian address
type Address struct {
	Street       string       `json:"street" validate:"required,min=3,max=200"`
	Number       string       `json:"number" validate:"required,min=1,max=10"`
	Complement   string       `json:"complement" validate:"max=50"`
	Neighborhood string       `json:"neighborhood" validate:"required,min=3,max=150"`
	City         string       `json:"city" validate:"required,min=2,max=150"`
	State        string       `json:"state" validate:"required,br_state"`
	Zipcode      string       `json:"zipcode" validate:"required,cep"`
	Coordinates  *Coordinates `json:"coordinates,omitempty" validate:"omitempty"`
}

// Coordinates represents geographic coordinates
type Coordinates struct {
	Lat float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng float64 `json:"lng" validate:"required,min=-180,max=180"`
}

// Person represents complete fake person data
type Person struct {
	Name          PersonName       `json:"name" validate:"required"`
	CPF           PersonCPF        `json:"cpf" validate:"required"`
	RG            PersonRG         `json:"rg" validate:"required"`
	Birthdate     string           `json:"birthdate" validate:"required,datetime=2006-01-02"`
	Age           int              `json:"age" validate:"required,min=0,max=120"`
	Gender        string           `json:"gender" validate:"required,oneof=male female"`
	Height        Height           `json:"height" validate:"required"`
	Weight        Weight           `json:"weight" validate:"required"`
	BMI           float64          `json:"bmi" validate:"required,min=10,max=60"`
	ZodiacSign    string           `json:"zodiacSign" validate:"required,min=3,max=20"`
	FavoriteColor string           `json:"favoriteColor" validate:"required,min=3,max=20"`
	BloodType     string           `json:"bloodType" validate:"required,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Filiation     Filiation        `json:"filiation" validate:"required"`
	Email         PersonEmail      `json:"email" validate:"required"`
	Phone         PersonPhone      `json:"phone" validate:"required"`
	Address       Address          `json:"address" validate:"required"`
	Profession    PersonProfession `json:"profession" validate:"required"`
	Company       PersonCompany    `json:"company" validate:"required"`
	Education     string           `json:"education" validate:"required,min=3,max=50"`
	MaritalStatus string           `json:"maritalStatus" validate:"required,oneof=single married divorced widowed"`
	BirthCity     string           `json:"birthCity" validate:"required,min=2,max=50"`
}

// PersonName represents the full name of the person divided into parts
type PersonName struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
	FullName  string `json:"fullName" validate:"required,min=5,max=100"`
}

// PersonCPF represents the CPF with and without mask
type PersonCPF struct {
	Masked   string `json:"masked" validate:"required,cpf"`
	Unmasked string `json:"unmasked" validate:"required,cpf"`
}

// PersonRG represents the RG with additional information
type PersonRG struct {
	Masked         string `json:"masked" validate:"required,rg"`
	Unmasked       string `json:"unmasked" validate:"required,rg"`
	State          string `json:"state" validate:"required,br_state"`
	Issuer         string `json:"issuer" validate:"required,min=3,max=10"`
	IssueDate      string `json:"issueDate" validate:"required,datetime=2006-01-02"`
	ExpirationDate string `json:"expirationDate" validate:"required,datetime=2006-01-02"`
}

// Height represents the height in different metrics
type Height struct {
	Centimeters float64 `json:"centimeters" validate:"required,min=50,max=250"`
	Meters      float64 `json:"meters" validate:"required,min=0.5,max=2.5"`
	Inches      float64 `json:"inches" validate:"required,min=19.7,max=98.4"`
	Feet        float64 `json:"feet" validate:"required,min=1.64,max=8.2"`
}

// Weight represents the weight in different metrics
type Weight struct {
	Kilograms float64 `json:"kilograms" validate:"required,min=30,max=300"`
	Pounds    float64 `json:"pounds" validate:"required,min=66,max=661"`
	Grams     float64 `json:"grams" validate:"required,min=30000,max=300000"`
}

// Filiation represents the person's filiation
type Filiation struct {
	Father string `json:"father" validate:"required,min=5,max=100"`
	Mother string `json:"mother" validate:"required,min=5,max=100"`
}

// PersonEmail represents the email with password
type PersonEmail struct {
	Address  string `json:"address" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

// PersonPhone represents the phone with different formats
type PersonPhone struct {
	InternationalFormat string `json:"internationalFormat" validate:"required,br_phone"`
	NationalFormat      string `json:"nationalFormat" validate:"required,br_phone"`
	CountryCode         string `json:"countryCode" validate:"required,len=2"`
	DDI                 int    `json:"ddi" validate:"required,eq=55"`
	E164Format          string `json:"e164Format" validate:"required,e164"`
}

// PersonProfession represents the person's profession and job
type PersonProfession struct {
	Title  string `json:"title" validate:"required,min=3,max=50"`
	Area   string `json:"area" validate:"required,min=3,max=50"`
	Salary string `json:"salary,omitempty" validate:"omitempty,min=5,max=20"`
}

// PersonCompany represents the company where the person works
type PersonCompany struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
	CNPJ string `json:"cnpj" validate:"required,cnpj"`
}

// StateInfo contains information about Brazilian states
type StateInfo struct {
	Code string
	Name string
	DDD  string
}

// GeneratorRequest contains parameters of a request
type GeneratorRequest struct {
	Quantity int
	Gender   string
	State    string
}

// CPFResponse represents the response of the CPF generation
type CPFResponse struct {
	CPF   string `json:"cpf" validate:"required,cpf"`
	Valid bool   `json:"valid" validate:"required"`
}

// CNPJResponse represents the response of the CNPJ generation
type CNPJResponse struct {
	CNPJ  string `json:"cnpj" validate:"required,cnpj"`
	Valid bool   `json:"valid" validate:"required"`
}

// RGResponse represents the response of the RG generation
type RGResponse struct {
	RG             string `json:"rg" validate:"required,rg"`
	State          string `json:"state" validate:"required,br_state"`
	Issuer         string `json:"issuer" validate:"required,min=3,max=10"`
	IssueDate      string `json:"issueDate" validate:"required,datetime=2006-01-02"`
	ExpirationDate string `json:"expirationDate" validate:"required,datetime=2006-01-02"`
}

// EmailResponse represents the response of the email generation
type EmailResponse struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Domain   string `json:"domain" validate:"required,min=3,max=50"`
}

// PhoneResponse represents the response of the phone generation
type PhoneResponse struct {
	Phone       string `json:"phone" validate:"required,br_phone"`
	Formatted   string `json:"formatted" validate:"required,br_phone"`
	Unformatted string `json:"unformatted" validate:"required,br_phone"`
	Type        string `json:"type" validate:"required,oneof=mobile landline"`
	DDD         string `json:"ddd" validate:"required,len=2"`
	State       string `json:"state" validate:"required,br_state"`
}

// BankAccountResponse represents the response of the bank account generation
type BankAccountResponse struct {
	Bank        Bank   `json:"bank" validate:"required"`
	Agency      string `json:"agency" validate:"required,min=4,max=6"`
	Account     string `json:"account" validate:"required,min=5,max=15"`
	AccountType string `json:"accountType" validate:"required,oneof=checking savings"`
}

// Bank represents a Brazilian bank
type Bank struct {
	Code string `json:"code" validate:"required,len=3"`
	Name string `json:"name" validate:"required,min=3,max=50"`
}

// CreditCardResponse represents the response of the credit card generation
type CreditCardResponse struct {
	Number         string `json:"number" validate:"required,len=19"`
	Brand          string `json:"brand" validate:"required,oneof=Visa Mastercard Elo Amex"`
	CVV            string `json:"cvv" validate:"required,len=3"`
	ExpirationDate string `json:"expirationDate" validate:"required,len=5"`
	HolderName     string `json:"holderName" validate:"required,min=5,max=100"`
}

// ZipcodeResponse represents the response of the zipcode generation
type ZipcodeResponse struct {
	Zipcode     string `json:"zipcode" validate:"required,cep"`
	Formatted   string `json:"formatted" validate:"required,cep"`
	Unformatted string `json:"unformatted" validate:"required,cep"`
	State       string `json:"state" validate:"required,br_state"`
	City        string `json:"city" validate:"required,min=2,max=50"`
}

// CompanyResponse represents the response of the company generation
type CompanyResponse struct {
	Name              string  `json:"name" validate:"required,min=3,max=100"`
	TradeName         string  `json:"tradeName" validate:"required,min=3,max=100"`
	CNPJ              string  `json:"cnpj" validate:"required,cnpj"`
	StateRegistration string  `json:"stateRegistration" validate:"required,min=9,max=14"`
	Email             string  `json:"email" validate:"required,email"`
	Phone             string  `json:"phone" validate:"required,br_phone"`
	Area              string  `json:"area" validate:"required,min=3,max=50"`
	Size              string  `json:"size" validate:"required,min=3,max=50"`
	OpeningDate       string  `json:"openingDate" validate:"required,datetime=2006-01-02"`
	ShareCapital      string  `json:"shareCapital" validate:"required,min=3,max=50"`
	FoundedAt         string  `json:"foundedAt" validate:"required,datetime=2006-01-02"`
	Address           Address `json:"address" validate:"required"`
}

// CPFValidationResponse represents the response of the CPF validation
type CPFValidationResponse struct {
	CPF   string `json:"cpf" validate:"required"`
	Valid bool   `json:"valid" validate:"required"`
}

// CNPJValidationResponse represents the response of the CNPJ validation
type CNPJValidationResponse struct {
	CNPJ  string `json:"cnpj" validate:"required"`
	Valid bool   `json:"valid" validate:"required"`
}

// RGValidationResponse represents the response of the RG validation
type RGValidationResponse struct {
	RG    string `json:"rg" validate:"required"`
	Valid bool   `json:"valid" validate:"required"`
}

// PhoneValidationResponse represents the response of the phone validation
type PhoneValidationResponse struct {
	Valid               bool   `json:"valid" validate:"required"`
	PhoneNumber         string `json:"phone_number" validate:"required"`
	CountryCode         string `json:"country_code" validate:"required,len=2"`
	DDI                 int32  `json:"ddi" validate:"required,eq=55"`
	NationalFormat      string `json:"national_format" validate:"required"`
	InternationalFormat string `json:"international_format" validate:"required"`
	E164Format          string `json:"e164_format" validate:"required,e164"`
	NumberType          string `json:"number_type" validate:"required,oneof=MOBILE FIXED_LINE FIXED_LINE_OR_MOBILE"`
}

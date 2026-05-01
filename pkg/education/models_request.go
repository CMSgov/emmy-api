package education

type Request struct {
	Address     *Address `json:"address,omitempty"`
	FirstName   string   `json:"firstName"`
	MiddleName  string   `json:"middleName,omitempty"`
	LastName    string   `json:"lastName"`
	DateOfBirth string   `json:"dateOfBirth"`
	SSN         string   `json:"ssn,omitempty"`
}

type Address struct {
	Street1    string `json:"street1,omitempty"`
	Street2    string `json:"street2,omitempty"`
	Street3    string `json:"street3,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

type nscRequest struct {
	ContactEmail     string                     `json:"contactEmail,omitempty"`
	CaseReferenceID  string                     `json:"caseReferenceId,omitempty"`
	AccountID        string                     `json:"accountId"`
	DateOfBirth      string                     `json:"dateOfBirth"`
	LastName         string                     `json:"lastName"`
	FirstName        string                     `json:"firstName"`
	MiddleName       string                     `json:"middleName,omitempty"`
	SSN              string                     `json:"ssn,omitempty"`
	OrganizationName string                     `json:"organizationName,omitempty"`
	EndClient        string                     `json:"endClient"`
	Terms            string                     `json:"terms"`
	PreviousNames    []nscPreviousName          `json:"previousNames,omitempty"`
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	City     string `json:"city,omitempty"`
	State    string `json:"state,omitempty"`
	ZipCode  string `json:"zipCode,omitempty"`
}

type nscPreviousName struct {
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
}

type BatchRequest struct {
	BatchID     string         `json:"batchId"`
	SubmittedBy string         `json:"submittedBy"`
	CallbackURL string         `json:"callbackUrl"`
	Students    []BatchStudent `json:"students"`
}

type BatchStudent struct {
	RecordID    string `json:"recordId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
	SSN         string `json:"ssn"`
}

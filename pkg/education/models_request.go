package education

type Request struct {
	FirstName   string   `json:"firstName"`
	MiddleName  string   `json:"middleName,omitempty"`
	LastName    string   `json:"lastName"`
	DateOfBirth string   `json:"dateOfBirth"`
	SSN         string   `json:"ssn,omitempty"`
	Address     *Address `json:"address,omitempty"`
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
	AccountID        string                     `json:"accountId"`
	OrganizationName string                     `json:"organizationName,omitempty"`
	CaseReferenceID  string                     `json:"caseReferenceId,omitempty"`
	ContactEmail     string                     `json:"contactEmail,omitempty"`
	DateOfBirth      string                     `json:"dateOfBirth"`
	LastName         string                     `json:"lastName"`
	FirstName        string                     `json:"firstName"`
	MiddleName       string                     `json:"middleName,omitempty"`
	SSN              string                     `json:"ssn,omitempty"`
	IdentityDetails  *nscRequestIdentityDetails `json:"identityDetails,omitempty"`
	EndClient        string                     `json:"endClient"`
	PreviousNames    []nscPreviousName          `json:"previousNames,omitempty"`
	Terms            string                     `json:"terms"`
}

type nscRequestIdentityDetails struct {
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

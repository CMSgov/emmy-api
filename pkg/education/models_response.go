package education

type Response struct {
	EnrollmentStatus string `json:"enrollmentStatus"`
}

type NSCResponse struct {
	ClientData          ClientDataResponse          `json:"clientData"`
	IdentityDetails     []IdentityDetailsResponse   `json:"identityDetails"`
	Status              StatusResponse              `json:"status"`
	StudentInfoProvided StudentInfoProvidedResponse `json:"studentInfoProvided"`
	TransactionDetails  TransactionDetailsResponse  `json:"transactionDetails"`
}

type ClientDataResponse struct {
	AccountID        string `json:"zaccountID"`
	CaseReferenceID  string `json:"caseReferenceId"`
	ContactEmail     string `json:"contactEmail"`
	OrganizationName string `json:"organizationName"`
}

type IdentityDetailsResponse struct {
	MatchElementName string `json:"matchElementName"`
	MatchLevel       string `json:"matchLevel"`
}

type StatusResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type StudentInfoProvidedResponse struct {
	DateOfBirth string `json:"dateOfBirth"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

type TransactionDetailsResponse struct {
	NotifiedDate      string `json:"notifiedDate"`
	NSCHit            string `json:"nscHit"`
	OrderID           string `json:"orderId"`
	RequestedBy       string `json:"requestedBy"`
	RequestedDate     string `json:"requestedDate"`
	SalesTax          string `json:"salesTax"`
	TransactionFee    string `json:"transactionFee"`
	TransactionID     string `json:"transactionId"`
	TransactionStatus string `json:"transactionStatus"`
	TransactionTotal  string `json:"transactionTotal"`
}

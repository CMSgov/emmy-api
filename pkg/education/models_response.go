package education

type EnrollmentStatus string

const (
	EnrollmentStatusFullTime         EnrollmentStatus = "FULL_TIME"
	EnrollmentStatusPartTime         EnrollmentStatus = "PART_TIME"
	EnrollmentStatusLessThanPartTime EnrollmentStatus = "LESS_THAN_PART_TIME"
	EnrollmentStatusEnrolled         EnrollmentStatus = "ENROLLED"
)

type Response struct {
	EnrollmentStatus EnrollmentStatus `json:"enrollmentStatus"`
	RawData          any              `json:"rawData"`
	DataSource       string           `json:"dataSource"`
}

type nscResponse struct {
	ClientData          nscClientData          `json:"clientData"`
	IdentityDetails     []nscMatchDetail       `json:"identityDetails"`
	Status              nscStatus              `json:"status"`
	StudentInfoProvided nscStudentInfoProvided `json:"studentInfoProvided"`
	TransactionDetails  nscTransactionDetails  `json:"transactionDetails"`
	EnrollmentDetails   []nscEnrollmentDetails `json:"enrollmentDetails"`
}

type nscClientData struct {
	AccountID        string `json:"zaccountID"`
	CaseReferenceID  string `json:"caseReferenceId"`
	ContactEmail     string `json:"contactEmail"`
	OrganizationName string `json:"organizationName"`
}

type nscMatchDetail struct {
	MatchElementName string `json:"matchElementName"`
	MatchLevel       string `json:"matchLevel"`
}

type nscStatus struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type nscStudentInfoProvided struct {
	DateOfBirth string `json:"dateOfBirth"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	MiddleName  string `json:"middleName"`
}

type nscTransactionDetails struct {
	NotifiedDate              string `json:"notifiedDate"`
	NSCHit                    string `json:"nscHit"`
	OrderID                   string `json:"orderId"`
	RequestedBy               string `json:"requestedBy"`
	RequestedDate             string `json:"requestedDate"`
	SalesTax                  string `json:"salesTax"`
	TransactionFee            string `json:"transactionFee"`
	TransactionID             string `json:"transactionId"`
	TransactionStatus         string `json:"transactionStatus"`
	TransactionTotal          string `json:"transactionTotal"`
	AppliedLikeSchoolMatching string `json:"appliedLikeSchoolMatching"`
	StudentComments           string `json:"studentComments"`
}

type nscEnrollmentDetails struct {
	OfficialSchoolName      string              `json:"officialSchoolName"`
	EnrollmentSinceDate     string              `json:"enrollmentSinceDate"`
	CurrentEnrollmentStatus string              `json:"currentEnrollmentStatus"`
	EnrollmentData          []nscEnrollmentData `json:"enrollmentData"`
	StudentAddress          *nscStudentAddress  `json:"studentAddress"`
}

type nscEnrollmentData struct {
	EnrollmentStatus          string `json:"enrollmentStatus"`
	TermBeginDate             string `json:"termBeginDate"`
	TermEndDate               string `json:"termEndDate"`
	SchoolCertifiedOnDate     string `json:"schoolCertifiedOnDate"`
	AnticipatedGraduationDate string `json:"anticipatedGraduationDate"`
	StudentClassCredential    string `json:"studentClassCredential"`
}

type nscStudentAddress struct {
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zipCode"`
	Province   string `json:"province"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

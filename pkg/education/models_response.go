package education

import (
	"time"

	"github.com/cmsgov/emmy-api/pkg/core"
)

type EnrollmentStatus string

const (
	EnrollmentStatusFullTime          EnrollmentStatus = "FULL_TIME"
	EnrollmentStatusThreeQuartersTime EnrollmentStatus = "THREE_QUARTERS_TIME"
	EnrollmentStatusHalfTime          EnrollmentStatus = "HALF_TIME"
	EnrollmentStatusLessThanHalfTime  EnrollmentStatus = "LESS_THAN_HALF_TIME"
	EnrollmentStatusUnknown           EnrollmentStatus = "ENROLLMENT_STATUS_UNKNOWN_CREDIT_TIMING"
)

type Response struct {
	EnrollmentStatus  EnrollmentStatus   `json:"enrollmentStatus"`
	EnrollmentDetails []EnrollmentDetail `json:"enrollmentDetails"`
	RawData           any                `json:"rawData"`
	DataSource        core.DataSource    `json:"dataSource"`
	Metadata          Metadata           `json:"metadata"`
}

type EnrollmentDetail struct {
	SchoolName       string           `json:"schoolName"`
	TermBeginDate    string           `json:"termBeginDate"`
	TermEndDate      string           `json:"termEndDate"`
	EnrollmentStatus EnrollmentStatus `json:"enrollmentStatus"`
}

type Metadata struct {
	APIVersion               string `json:"apiVersion"`
	Environment              string `json:"environment"`
	RequestTimestamp         string `json:"requestTimestamp"`
	ResponseTimestamp        string `json:"responseTimestamp"`
	TransactionID            string `json:"transaction-id"` //nolint:tagliatelle // kebab-case is required for transaction-id for compatibility
	DatasourceDurationMillis int64  `json:"datasourceDurationMillis"`
}

type BatchJobStatusResponse struct {
	SubmittedAt             *time.Time `json:"submittedAt"`
	UpdatedAt               *time.Time `json:"updatedAt"`
	EstimatedCompletionTime *time.Time `json:"estimatedCompletionTime,omitempty"`
	BatchJobID              string     `json:"batchJobID"` //nolint:tagliatelle // kebab-case is required for batch-job-id for compatibility
	Status                  string     `json:"status"`
	TotalRecords            int        `json:"totalRecords"`
	ProcessedRecords        int        `json:"processedRecords"`
	SuccessCount            int        `json:"successCount"`
	FailureCount            int        `json:"failureCount"`
}

type BatchJobDetailsResponse struct {
	BatchJobID string               `json:"batchJobId"`
	Results    []BatchStudentResult `json:"results"`
}

type BatchStudentResult struct {
	Results         *StudentResults `json:"results,omitempty"`
	RecordID        string          `json:"recordId"`
	Status          string          `json:"status"`
	FoundEnrollment bool            `json:"foundEnrollment"`
}

type StudentResults struct {
	EnrollmentStatus  EnrollmentStatus   `json:"enrollmentStatus"`
	EnrollmentDetails []EnrollmentDetail `json:"enrollmentDetails"`
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
	AccountID        string `json:"zaccountID"` //nolint:tagliatelle // NSC payload uses this exact casing.
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
	StudentAddress          *nscStudentAddress  `json:"studentAddress"`
	OfficialSchoolName      string              `json:"officialSchoolName"`
	EnrollmentSinceDate     string              `json:"enrollmentSinceDate"`
	CurrentEnrollmentStatus string              `json:"currentEnrollmentStatus"`
	EnrollmentData          []nscEnrollmentData `json:"enrollmentData"`
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

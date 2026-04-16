package education

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestLookupEnrollmentStatus_SuccessFallsBackToEnrolledOnPositiveHit(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"status":{"code":"0","message":"Successful","severity":"Info"},
				"transactionDetails":{"nscHit":"Y","transactionStatus":"CNF"},
				"enrollmentDetails":[{"currentEnrollmentStatus":"CC"}]
			}`)),
		},
	}

	svc := New(&core.NSCConfig{
		AccountID: "10053523",
		SubmitURL: "https://example.test/submit",
	}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		MiddleName:  "Marie",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
		SSN:         "123-45-6789",
		Address: &Address{
			Street1:    "17020 Tortoise St",
			City:       "Round Rock",
			State:      "TX",
			PostalCode: "78664",
			Country:    "USA",
		},
	})
	require.NoError(t, err)

	require.True(t, ft.called)
	require.NotNil(t, ft.req)
	require.Equal(t, http.MethodPost, ft.req.Method)
	require.Equal(t, "application/json", ft.req.Header.Get("Content-Type"))
	require.Equal(t, "application/json", ft.req.Header.Get("Accept"))

	_, err = io.ReadAll(ft.req.Body)
	require.NoError(t, err)
	require.Equal(t, EnrollmentStatusEnrolled, out.EnrollmentStatus)
	require.Equal(t, core.DataSourceNSC, out.DataSource)

	rawData, ok := out.RawData.(map[string]any)
	require.True(t, ok, "RawData should be a map[string]any")
	transactionDetails := rawData["transactionDetails"].(map[string]any)
	require.Equal(t, "Y", transactionDetails["nscHit"])
	enrollmentDetails := rawData["enrollmentDetails"].([]any)
	firstDetail := enrollmentDetails[0].(map[string]any)
	require.Equal(t, "CC", firstDetail["currentEnrollmentStatus"])
}

func TestLookupEnrollmentStatus_MapsSpecificEnrollmentStatus(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y"},
				"enrollmentDetails":[{"currentEnrollmentStatus":"CC","enrollmentData":[{"enrollmentStatus":"H"}]}]
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.NoError(t, err)
	require.Equal(t, EnrollmentStatusPartTime, out.EnrollmentStatus)
	require.Len(t, out.EnrollmentDetails, 1)
	require.Equal(t, EnrollmentStatusPartTime, out.EnrollmentDetails[0].EnrollmentStatus)
}

func TestLookupEnrollmentStatus_MapsLessThanHalfTimeToLessThanPartTime(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y"},
				"enrollmentDetails":[{"currentEnrollmentStatus":"CC","enrollmentData":[{"enrollmentStatus":"L"}]}]
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.NoError(t, err)
	require.Equal(t, EnrollmentStatusLessThanPartTime, out.EnrollmentStatus)
}

func TestLookupEnrollmentStatus_NoHitReturnsNotFound(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"status":{"code":"0","message":"Successful","severity":"Info"},
				"transactionDetails":{"nscHit":"N","transactionStatus":"CNF"},
				"enrollmentDetails":[{"currentEnrollmentStatus":"CN"}]
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	_, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestLookupEnrollmentStatus_Non2xxReturnsError(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusBadGateway,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":"upstream failure"}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	_, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.Error(t, err)
	require.False(t, errors.Is(err, ErrNotFound))
}

func TestLookupEnrollmentStatus_CurrentlyNotEnrolledReturnsNotFound(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y","transactionStatus":"CNF"},
				"enrollmentDetails":[{"currentEnrollmentStatus":"CN"}]
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	_, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestLookupEnrollmentStatus_NullEnrollmentDetails(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y","transactionStatus":"CNF"},
				"enrollmentDetails":null
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.NoError(t, err)
	require.Equal(t, EnrollmentStatusEnrolled, out.EnrollmentStatus)
}

func TestLookupEnrollmentStatus_MapsMultipleEnrollmentDetails(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y"},
				"enrollmentDetails":[
					{
						"officialSchoolName":"University A",
						"enrollmentData":[
							{"enrollmentStatus":"F","termBeginDate":"2023-01-01","termEndDate":"2023-05-01"}
						]
					},
					{
						"officialSchoolName":"University B",
						"enrollmentData":[
							{"enrollmentStatus":"H","termBeginDate":"2022-09-01","termEndDate":"2022-12-31"}
						]
					}
				]
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupEnrollmentStatus(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
	})
	require.NoError(t, err)
	require.Equal(t, EnrollmentStatusFullTime, out.EnrollmentStatus)
	require.Len(t, out.EnrollmentDetails, 2)

	require.Equal(t, "University A", out.EnrollmentDetails[0].SchoolName)
	require.Equal(t, "2023-01-01", out.EnrollmentDetails[0].TermBeginDate)
	require.Equal(t, "2023-05-01", out.EnrollmentDetails[0].TermEndDate)
	require.Equal(t, EnrollmentStatusFullTime, out.EnrollmentDetails[0].EnrollmentStatus)

	require.Equal(t, "University B", out.EnrollmentDetails[1].SchoolName)
	require.Equal(t, "2022-09-01", out.EnrollmentDetails[1].TermBeginDate)
	require.Equal(t, "2022-12-31", out.EnrollmentDetails[1].TermEndDate)
	require.Equal(t, EnrollmentStatusPartTime, out.EnrollmentDetails[1].EnrollmentStatus)
}

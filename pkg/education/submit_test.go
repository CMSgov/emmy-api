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
				"enrollmentDetails":{"currentEnrollmentStatus":"CC"}
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
	require.Equal(t, "text/plain", ft.req.Header.Get("Content-Type"))
	require.Equal(t, "application/json", ft.req.Header.Get("Accept"))

	body, err := io.ReadAll(ft.req.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"accountId":"10053523",
		"dateOfBirth":"1988-10-24",
		"lastName":"Oyola",
		"firstName":"Lynette",
		"middleName":"Marie",
		"ssn":"123-45-6789",
		"identityDetails":{
			"address1":"17020 Tortoise St",
			"city":"Round Rock",
			"state":"TX",
			"zipCode":"78664"
		},
		"endClient":"CMS",
		"terms":"y"
	}`, string(body))

	require.Equal(t, EnrollmentStatusEnrolled, out.EnrollmentStatus)
}

func TestLookupEnrollmentStatus_MapsSpecificEnrollmentStatus(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y"},
				"enrollmentDetails":{"currentEnrollmentStatus":"CC","enrollmentData":[{"enrollmentStatus":"H"}]}
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
}

func TestLookupEnrollmentStatus_MapsLessThanHalfTimeToLessThanPartTime(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"transactionDetails":{"nscHit":"Y"},
				"enrollmentDetails":{"currentEnrollmentStatus":"CC","enrollmentData":[{"enrollmentStatus":"L"}]}
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
				"enrollmentDetails":{"currentEnrollmentStatus":"CN"}
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
				"enrollmentDetails":{"currentEnrollmentStatus":"CN"}
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

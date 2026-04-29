package education

import (
	"bytes"
	"context"
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
	require.Equal(t, EnrollmentStatusUnknown, out.EnrollmentStatus)
	require.Equal(t, core.DataSourceNSC, out.DataSource)

	rawData, ok := out.RawData.(map[string]any)
	require.True(t, ok, "RawData should be a map[string]any")
	transactionDetails, ok := rawData["transactionDetails"].(map[string]any)
	require.True(t, ok, "transactionDetails should be a map[string]any")
	require.Equal(t, "Y", transactionDetails["nscHit"])
	enrollmentDetails, ok := rawData["enrollmentDetails"].([]any)
	require.True(t, ok, "enrollmentDetails should be a []any")
	require.NotEmpty(t, enrollmentDetails)
	firstDetail, ok := enrollmentDetails[0].(map[string]any)
	require.True(t, ok, "first enrollment detail should be a map[string]any")
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
	require.Equal(t, EnrollmentStatusHalfTime, out.EnrollmentStatus)
	require.Len(t, out.EnrollmentDetails, 1)
	require.Equal(t, EnrollmentStatusHalfTime, out.EnrollmentDetails[0].EnrollmentStatus)
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
	require.Equal(t, EnrollmentStatusLessThanHalfTime, out.EnrollmentStatus)
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
	require.NotErrorIs(t, err, ErrNotFound)
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
	require.Equal(t, EnrollmentStatusUnknown, out.EnrollmentStatus)
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
	require.Equal(t, EnrollmentStatusHalfTime, out.EnrollmentDetails[1].EnrollmentStatus)
}

func TestResolveEnrollmentStatus_PrioritizesStatus(t *testing.T) {
	tests := []struct {
		name     string
		expected EnrollmentStatus
		resp     nscResponse
	}{
		{
			name: "FullTime overrides ThreeQuartersTime",
			resp: nscResponse{
				EnrollmentDetails: []nscEnrollmentDetails{
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "Q"},
						},
					},
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "F"},
						},
					},
				},
			},
			expected: EnrollmentStatusFullTime,
		},
		{
			name: "ThreeQuartersTime overrides HalfTime",
			resp: nscResponse{
				EnrollmentDetails: []nscEnrollmentDetails{
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "H"},
						},
					},
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "Q"},
						},
					},
				},
			},
			expected: EnrollmentStatusThreeQuartersTime,
		},
		{
			name: "HalfTime overrides LessThanPartTime",
			resp: nscResponse{
				EnrollmentDetails: []nscEnrollmentDetails{
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "L"},
						},
					},
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "H"},
						},
					},
				},
			},
			expected: EnrollmentStatusHalfTime,
		},
		{
			name: "LessThanPartTime overrides Unknown",
			resp: nscResponse{
				EnrollmentDetails: []nscEnrollmentDetails{
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "Y"},
						},
					},
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "L"},
						},
					},
				},
			},
			expected: EnrollmentStatusLessThanHalfTime,
		},
		{
			name: "Unknown when only Unknown is present",
			resp: nscResponse{
				EnrollmentDetails: []nscEnrollmentDetails{
					{
						EnrollmentData: []nscEnrollmentData{
							{EnrollmentStatus: "Y"},
						},
					},
				},
			},
			expected: EnrollmentStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, ok := resolveEnrollmentStatus(tt.resp)
			require.True(t, ok)
			require.Equal(t, tt.expected, status)
		})
	}
}

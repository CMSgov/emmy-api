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

func TestSubmit_Success(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"clientData":{"zaccountID":"","caseReferenceId":"","contactEmail":"","organizationName":"Lynette"},
				"identityDetails":null,
				"status":{"code":"0","message":"Successful","severity":"Info"},
				"studentInfoProvided":{"dateOfBirth":"1988-10-24","firstName":"Lynette","lastName":"Oyola"},
				"transactionDetails":{"notifiedDate":"","nscHit":"Y","orderId":"","requestedBy":"","requestedDate":"","salesTax":"0.00","transactionFee":"0.00","transactionId":"","transactionStatus":"CNF","transactionTotal":"0.00"}
			}`)),
		},
	}

	svc := New(&core.NSCConfig{SubmitURL: "https://example.test/submit"}, Options{
		HTTPClient: ft,
	})

	out, err := svc.Submit(context.Background(), Request{
		AccountID:        "10053523",
		OrganizationName: "Lynette",
		DateOfBirth:      "1988-10-24",
		LastName:         "Oyola",
		FirstName:        "Lynette",
		Terms:            "y",
		EndClient:        "CMS",
	})
	require.NoError(t, err)

	require.True(t, ft.called)
	require.NotNil(t, ft.req)
	require.Equal(t, http.MethodPost, ft.req.Method)
	require.Equal(t, "application/json", ft.req.Header.Get("Content-Type"))
	require.Equal(t, "application/json", ft.req.Header.Get("Accept"))

	require.Equal(t, "0", out.Status.Code)
	require.Equal(t, "Lynette", out.StudentInfoProvided.FirstName)
	require.Equal(t, "Oyola", out.StudentInfoProvided.LastName)
}

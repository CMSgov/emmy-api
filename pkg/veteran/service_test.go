package veteran

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/cmsgov/emmy-api/pkg/core"
	"github.com/stretchr/testify/require"
)

type fakeTransport struct {
	called bool
	req    *http.Request
	resp   *http.Response
	err    error
}

func (f *fakeTransport) Do(req *http.Request) (*http.Response, error) {
	f.called = true
	f.req = req
	return f.resp, f.err
}

func TestNew_UsesInjectedHTTPClient(t *testing.T) {
	cfg := &core.VAConfig{
		BaseURL:        "https://example.test",
		TokenURL:       "https://example.test/token",
		ClientID:       "id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: "/tmp/private.pem",
	}

	ft := &fakeTransport{}

	svc := New(cfg, Options{
		HTTPClient: ft,
	})

	impl, ok := svc.(*service)
	require.True(t, ok, "New should return *service implementation")
	require.Same(t, cfg, impl.cfg)
	require.Same(t, ft, impl.client)
}

func TestLookupDisabilityRating_Success(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(bytes.NewBufferString(`{
				"data": {
					"attributes": {
						"combined_disability_rating": 70
					}
				}
			}`)),
		},
	}

	svc := New(&core.VAConfig{
		BaseURL:        "https://example.test",
		TokenURL:       "https://example.test/token",
		ClientID:       "id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: "/tmp/private.pem",
	}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupDisabilityRating(context.Background(), Request{
		FirstName:   "Lynette",
		MiddleName:  "Marie",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
		SSN:         "123-45-6789",
		Address: &Address{
			Street1:    "17020 Tortoise St",
			Street2:    "Apt 4B",
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
	require.Equal(t, "https://example.test/restricted/disability_rating", ft.req.URL.String())
	require.Equal(t, "application/json", ft.req.Header.Get("Content-Type"))
	require.Equal(t, "application/json", ft.req.Header.Get("Accept"))

	reqBody, err := io.ReadAll(ft.req.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"ssn":"123-45-6789",
		"first_name":"Lynette",
		"middle_name":"Marie",
		"last_name":"Oyola",
		"birth_date":"1988-10-24",
		"street_address_line1":"17020 Tortoise St",
		"street_address_line2":"Apt 4B",
		"city":"Round Rock",
		"state":"TX",
		"zipcode":"78664",
		"country":"USA"
	}`, string(reqBody))

	require.Equal(t, 70, out.CombinedDisabilityRating)
	require.Equal(t, "Veteran's Affairs", out.DataSource)
	rawData, ok := out.RawData.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Veteran's Affairs", rawData["dataSource"])
	reqMap, ok := rawData["reqBody"].(Request)
	require.True(t, ok)
	require.Equal(t, "Lynette", reqMap.FirstName)
	require.Equal(t, "123-45-6789", reqMap.SSN)
}

func TestLookupDisabilityRating_AddressOnlySuccess(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(`{"data":{"attributes":{"combined_disability_rating":80}}}`)),
		},
	}

	svc := New(&core.VAConfig{
		BaseURL:        "https://example.test",
		TokenURL:       "https://example.test/token",
		ClientID:       "id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: "/tmp/private.pem",
	}, Options{
		HTTPClient: ft,
	})

	out, err := svc.LookupDisabilityRating(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
		Address: &Address{
			Street1:    "17020 Tortoise St",
			City:       "Round Rock",
			State:      "TX",
			PostalCode: "78664",
			Country:    "USA",
		},
	})
	require.NoError(t, err)
	require.Equal(t, 80, out.CombinedDisabilityRating)
	require.Equal(t, "Veteran's Affairs", out.DataSource)
	rawData, ok := out.RawData.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Veteran's Affairs", rawData["dataSource"])
	reqMap, ok := rawData["reqBody"].(Request)
	require.True(t, ok)
	require.Equal(t, "Lynette", reqMap.FirstName)
	require.Equal(t, "17020 Tortoise St", reqMap.Address.Street1)

	reqBody, err := io.ReadAll(ft.req.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"first_name":"Lynette",
		"last_name":"Oyola",
		"birth_date":"1988-10-24",
		"street_address_line1":"17020 Tortoise St",
		"city":"Round Rock",
		"state":"TX",
		"zipcode":"78664",
		"country":"USA"
	}`, string(reqBody))
}

func TestLookupDisabilityRating_NotFound(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		},
	}

	svc := New(&core.VAConfig{
		BaseURL:        "https://example.test",
		TokenURL:       "https://example.test/token",
		ClientID:       "id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: "/tmp/private.pem",
	}, Options{
		HTTPClient: ft,
	})

	_, err := svc.LookupDisabilityRating(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
		SSN:         "123-45-6789",
	})

	require.ErrorIs(t, err, ErrNotFound)
}

func TestLookupDisabilityRating_Non2xx(t *testing.T) {
	ft := &fakeTransport{
		resp: &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       io.NopCloser(bytes.NewBufferString(`{"detail":"upstream error"}`)),
		},
	}

	svc := New(&core.VAConfig{
		BaseURL:        "https://example.test",
		TokenURL:       "https://example.test/token",
		ClientID:       "id",
		TokenAudience:  "https://example.okta.com/oauth2/default/v1/token",
		PrivateKeyPath: "/tmp/private.pem",
	}, Options{
		HTTPClient: ft,
	})

	_, err := svc.LookupDisabilityRating(context.Background(), Request{
		FirstName:   "Lynette",
		LastName:    "Oyola",
		DateOfBirth: "1988-10-24",
		SSN:         "123-45-6789",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "status=502")
	require.NotErrorIs(t, err, ErrNotFound)
}

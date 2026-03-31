package education

import (
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
	cfg := &core.NSCConfig{
		ClientID:     "id",
		ClientSecret: "secret",
		TokenURL:     "https://example.com/token",
		SubmitURL:    "https://example.com/submit",
	}

	fd := &fakeTransport{}

	svc := New(cfg, Options{
		HTTPClient: fd,
	})

	impl, ok := svc.(*service)
	require.True(t, ok, "New should return *service implementation")
	require.Same(t, cfg, impl.cfg, "should preserve cfg pointer")
	require.Same(t, fd, impl.client, "should use injected HTTP client")
}

package education

import (
	"context"
	"net/http"

	"github.com/cmsgov/emmy-api/pkg/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func nscHTTPClient(ctx context.Context, cfg *core.NSCConfig) *http.Client {
	base := &http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				r.Header = via[0].Header.Clone()
			}
			return nil
		},
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, base)

	cc := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       []string{"vs.api.insights"},
	}

	return oauth2.NewClient(ctx, cc.TokenSource(ctx))
}
